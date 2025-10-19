/**
 * Security Platform - Cloudflare Workers
 * 
 * 部署說明:
 * 1. 前往 https://dash.cloudflare.com/
 * 2. 點擊 "Workers & Pages" → "Create application" → "Create Worker"
 * 3. 將此檔案的完整內容複製貼上到編輯器
 * 4. 點擊 "Save and Deploy"
 * 
 * API 端點:
 * - GET /api/v1/health - 健康檢查
 * - GET /api/v1/status - 系統狀態
 * - GET /api/v1/security/threats - 安全威脅
 * - GET /api/v1/network/stats - 網路統計
 * - GET /api/v1/devices - 設備列表
 */

// Simple rate limiter using in-memory storage
const rateLimitMap = new Map();

function checkRateLimit(request) {
  const ip = request.headers.get('CF-Connecting-IP') || 'unknown';
  const now = Date.now();
  const windowMs = 60 * 1000; // 1 minute
  const maxRequests = 150;

  if (!rateLimitMap.has(ip)) {
    rateLimitMap.set(ip, { count: 1, resetTime: now + windowMs });
    return { allowed: true };
  }

  const record = rateLimitMap.get(ip);
  
  if (now > record.resetTime) {
    record.count = 1;
    record.resetTime = now + windowMs;
    return { allowed: true };
  }

  if (record.count >= maxRequests) {
    return { 
      allowed: false, 
      retryAfter: Math.ceil((record.resetTime - now) / 1000) 
    };
  }

  record.count++;
  return { allowed: true };
}

// Simple router implementation
function route(request) {
  const url = new URL(request.url);
  const path = url.pathname;

  // Health check endpoint
  if (path === '/api/v1/health' && request.method === 'GET') {
    return new Response(JSON.stringify({
      healthy: true,
      timestamp: new Date().toISOString(),
      version: '1.0.0',
      uptime: 'running'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }

  // Status endpoint
  if (path === '/api/v1/status' && request.method === 'GET') {
    return new Response(JSON.stringify({
      status: 'operational',
      timestamp: new Date().toISOString(),
      version: '1.0.0',
      services: {
        worker: 'healthy',
        database: 'not_configured',
        cache: 'not_configured'
      }
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }

  // Security threats endpoint
  if (path === '/api/v1/security/threats' && request.method === 'GET') {
    return new Response(JSON.stringify({
      threats: [],
      total: 0,
      message: 'Database not configured yet. This is a demo endpoint.'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }

  // Network stats endpoint
  if (path === '/api/v1/network/stats' && request.method === 'GET') {
    return new Response(JSON.stringify({
      totalRequests: 0,
      blockedRequests: 0,
      activeConnections: 0,
      message: 'Database not configured yet. This is a demo endpoint.'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }

  // Devices endpoint
  if (path === '/api/v1/devices' && request.method === 'GET') {
    return new Response(JSON.stringify({
      devices: [],
      total: 0,
      message: 'Database not configured yet. This is a demo endpoint.'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }

  // 404 for all other routes
  return new Response(JSON.stringify({
    error: 'Not Found',
    message: `Endpoint ${path} not found`,
    available_endpoints: [
      '/api/v1/health',
      '/api/v1/status',
      '/api/v1/security/threats',
      '/api/v1/network/stats',
      '/api/v1/devices'
    ]
  }), { 
    status: 404,
    headers: { 'Content-Type': 'application/json' }
  });
}

/**
 * Main request handler
 */
async function handleRequest(request) {
  const url = new URL(request.url);

  // CORS headers
  const corsHeaders = {
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type, Authorization',
  };

  // Handle CORS preflight
  if (request.method === 'OPTIONS') {
    return new Response(null, { headers: corsHeaders });
  }

  // Check rate limit
  const rateLimitResult = checkRateLimit(request);
  if (!rateLimitResult.allowed) {
    return new Response(
      JSON.stringify({
        error: 'Rate limit exceeded',
        message: 'Too many requests. Please try again later.',
        retryAfter: rateLimitResult.retryAfter,
      }),
      {
        status: 429,
        headers: {
          ...corsHeaders,
          'Content-Type': 'application/json',
          'Retry-After': String(rateLimitResult.retryAfter),
        },
      }
    );
  }

  // Route API requests
  try {
    const response = route(request);
    
    // Add CORS headers to response
    const headers = new Headers(response.headers);
    Object.entries(corsHeaders).forEach(([key, value]) => {
      headers.set(key, value);
    });

    return new Response(response.body, {
      status: response.status,
      statusText: response.statusText,
      headers,
    });
  } catch (error) {
    console.error('Request handling error:', error);
    return new Response(
      JSON.stringify({
        error: 'Internal Server Error',
        message: error.message,
      }),
      {
        status: 500,
        headers: {
          ...corsHeaders,
          'Content-Type': 'application/json',
        },
      }
    );
  }
}

// Export for Cloudflare Workers
export default {
  fetch: handleRequest,
};
