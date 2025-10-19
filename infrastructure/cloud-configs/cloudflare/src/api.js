/**
 * API request handlers for Cloudflare Workers
 */

/**
 * Main API request handler
 */
export async function handleApiRequest(request, env, ctx) {
  const url = new URL(request.url);
  const path = url.pathname;

  // Route to appropriate handler
  if (path === '/api/v1/status') {
    return handleStatus(request, env);
  } else if (path === '/api/v1/health') {
    return handleHealth(request, env);
  } else if (path.startsWith('/api/v1/security/threats')) {
    return handleSecurityThreats(request, env);
  } else if (path.startsWith('/api/v1/network/')) {
    return handleNetwork(request, env);
  } else if (path.startsWith('/api/v1/devices')) {
    return handleDevices(request, env);
  } else if (path === '/api/v1/ml/detect') {
    return handleMLDetect(request, env);
  }

  return jsonResponse({ error: 'Not found' }, 404);
}

/**
 * System status endpoint
 */
async function handleStatus(request, env) {
  const stats = await getSystemStats(env);
  
  return jsonResponse({
    status: 'operational',
    version: '1.0.0',
    uptime: stats.uptime,
    services: {
      database: 'healthy',
      cache: 'healthy',
      websocket: 'healthy',
    },
    timestamp: new Date().toISOString(),
  });
}

/**
 * Health check endpoint
 */
async function handleHealth(request, env) {
  try {
    // Test D1 database connection
    const dbResult = await env.DB.prepare('SELECT 1 as health').first();
    
    return jsonResponse({
      healthy: true,
      checks: {
        database: dbResult?.health === 1,
        kv: true,
      },
    });
  } catch (error) {
    return jsonResponse(
      {
        healthy: false,
        error: error.message,
      },
      503
    );
  }
}

/**
 * Security threats endpoint
 */
async function handleSecurityThreats(request, env) {
  const url = new URL(request.url);
  const threatId = url.pathname.split('/').pop();

  if (request.method === 'GET') {
    // Query threats from D1 database
    const { results } = await env.DB.prepare(
      `SELECT * FROM threats 
       WHERE discovered_at > datetime('now', '-7 days') 
       ORDER BY discovered_at DESC 
       LIMIT 100`
    ).all();

    return jsonResponse({
      threats: results || [],
      total: results?.length || 0,
    });
  } else if (request.method === 'POST' && request.url.includes('/block')) {
    // Block threat
    await env.DB.prepare(
      'UPDATE threats SET status = ? WHERE id = ?'
    ).bind('blocked', threatId).run();

    return jsonResponse({ success: true, id: threatId });
  }

  return jsonResponse({ error: 'Method not allowed' }, 405);
}

/**
 * Network stats endpoint
 */
async function handleNetwork(request, env) {
  const url = new URL(request.url);
  
  if (url.pathname === '/api/v1/network/stats') {
    // Get cached stats or fetch from database
    const cacheKey = 'network:stats';
    let stats = await env.CACHE.get(cacheKey, 'json');

    if (!stats) {
      stats = await getNetworkStats(env);
      await env.CACHE.put(cacheKey, JSON.stringify(stats), {
        expirationTtl: 60, // Cache for 1 minute
      });
    }

    return jsonResponse(stats);
  } else if (url.pathname.startsWith('/api/v1/network/blocked-ips')) {
    const ip = url.pathname.split('/').pop();
    
    if (request.method === 'DELETE') {
      await env.DB.prepare('DELETE FROM blocked_ips WHERE ip = ?')
        .bind(ip)
        .run();
      
      return jsonResponse({ success: true, ip });
    }
  }

  return jsonResponse({ error: 'Not found' }, 404);
}

/**
 * Devices endpoint
 */
async function handleDevices(request, env) {
  const { results } = await env.DB.prepare(
    'SELECT * FROM devices ORDER BY last_seen DESC'
  ).all();

  return jsonResponse({
    devices: results || [],
    total: results?.length || 0,
  });
}

/**
 * ML threat detection endpoint
 */
async function handleMLDetect(request, env) {
  if (request.method !== 'POST') {
    return jsonResponse({ error: 'Method not allowed' }, 405);
  }

  try {
    const data = await request.json();
    
    // In a real implementation, this would call the ML service
    // For now, return a mock response
    const result = {
      threat_detected: false,
      confidence: 0.95,
      threat_type: 'normal',
      analysis_time_ms: 8,
      timestamp: new Date().toISOString(),
    };

    // Log to D1
    await env.DB.prepare(
      `INSERT INTO ml_detections (source_ip, threat_type, confidence, created_at) 
       VALUES (?, ?, ?, datetime('now'))`
    ).bind(
      data.source_ip || 'unknown',
      result.threat_type,
      result.confidence
    ).run();

    return jsonResponse(result);
  } catch (error) {
    return jsonResponse({ error: error.message }, 400);
  }
}

/**
 * Helper functions
 */

async function getSystemStats(env) {
  return {
    uptime: 0, // Workers don't have traditional uptime
    requests_today: 0, // Would need to track this
  };
}

async function getNetworkStats(env) {
  const { results } = await env.DB.prepare(`
    SELECT 
      COUNT(*) as total_connections,
      SUM(bytes_in) as total_bytes_in,
      SUM(bytes_out) as total_bytes_out
    FROM network_stats
    WHERE timestamp > datetime('now', '-1 hour')
  `).first();

  return {
    total_connections: results?.total_connections || 0,
    bandwidth: {
      in: results?.total_bytes_in || 0,
      out: results?.total_bytes_out || 0,
    },
    latency_ms: 12,
    packet_loss: 0.01,
  };
}

function jsonResponse(data, status = 200) {
  return new Response(JSON.stringify(data), {
    status,
    headers: {
      'Content-Type': 'application/json',
    },
  });
}

