/**
 * Cloudflare Workers Entry Point for Security Platform
 * 
 * This worker handles:
 * - REST API endpoints
 * - Basic rate limiting
 * - CORS handling
 * 
 * Free Tier Limits:
 * - 10M requests/month
 * - 30M CPU milliseconds
 * - 5GB D1 storage
 * - 1GB KV storage
 */

import { Router } from 'itty-router';

// Create router
const router = Router();

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

// API Routes
router
  .get('/api/v1/status', async () => {
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
  })
  .get('/api/v1/health', async () => {
    return new Response(JSON.stringify({
      healthy: true,
      timestamp: new Date().toISOString(),
      version: '1.0.0',
      uptime: 'running'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  })
  .get('/api/v1/security/threats', async () => {
    return new Response(JSON.stringify({
      threats: [],
      total: 0,
      message: 'Database not configured yet'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  })
  .get('/api/v1/network/stats', async () => {
    return new Response(JSON.stringify({
      totalRequests: 0,
      blockedRequests: 0,
      activeConnections: 0,
      message: 'Database not configured yet'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  })
  .get('/api/v1/devices', async () => {
    return new Response(JSON.stringify({
      devices: [],
      total: 0,
      message: 'Database not configured yet'
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  })
  .all('*', () => new Response('Not Found', { status: 404 }));

/**
 * Main request handler
 */
async function handleRequest(request, env, ctx) {
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
    const response = await router.handle(request, env, ctx);
    
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

export default {
  fetch: handleRequest,
};

