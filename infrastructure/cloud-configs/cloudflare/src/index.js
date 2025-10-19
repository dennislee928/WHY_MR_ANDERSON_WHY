/**
 * Cloudflare Workers Entry Point for Security Platform
 * 
 * This worker handles:
 * - REST API endpoints
 * - WebSocket connections
 * - Rate limiting
 * - Caching
 * 
 * Free Tier Limits:
 * - 10M requests/month
 * - 30M CPU milliseconds
 * - 5GB D1 storage
 * - 1GB KV storage
 */

import { Router } from 'itty-router';
import { handleWebSocket } from './websocket';
import { handleApiRequest } from './api';
import { RateLimiter } from './middleware/rateLimit';
import { cacheMiddleware } from './middleware/cache';

// Create router
const router = Router();

// Rate limiter (150 requests per minute per IP)
const rateLimiter = new RateLimiter({
  limit: 150,
  window: 60 * 1000, // 1 minute
});

// API Routes
router
  .get('/api/v1/status', cacheMiddleware(60), handleApiRequest)
  .get('/api/v1/health', handleApiRequest)
  .get('/api/v1/security/threats', handleApiRequest)
  .post('/api/v1/security/threats/:id/block', handleApiRequest)
  .get('/api/v1/network/stats', cacheMiddleware(30), handleApiRequest)
  .delete('/api/v1/network/blocked-ips/:ip', handleApiRequest)
  .get('/api/v1/devices', cacheMiddleware(60), handleApiRequest)
  .post('/api/v1/ml/detect', handleApiRequest)
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
  const rateLimitResult = await rateLimiter.check(request, env);
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

  // Handle WebSocket upgrade
  if (request.headers.get('Upgrade') === 'websocket') {
    return handleWebSocket(request, env);
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

/**
 * Durable Object for WebSocket management
 */
export class WebSocketManager {
  constructor(state, env) {
    this.state = state;
    this.env = env;
    this.sessions = new Map();
  }

  async fetch(request) {
    const webSocketPair = new WebSocketPair();
    const [client, server] = Object.values(webSocketPair);

    this.handleSession(server);

    return new Response(null, {
      status: 101,
      webSocket: client,
    });
  }

  handleSession(websocket) {
    websocket.accept();
    
    const sessionId = crypto.randomUUID();
    this.sessions.set(sessionId, websocket);

    websocket.addEventListener('message', async (event) => {
      try {
        const data = JSON.parse(event.data);
        await this.handleMessage(sessionId, data);
      } catch (error) {
        console.error('WebSocket message error:', error);
      }
    });

    websocket.addEventListener('close', () => {
      this.sessions.delete(sessionId);
    });
  }

  async handleMessage(sessionId, data) {
    const websocket = this.sessions.get(sessionId);
    
    if (data.type === 'ping') {
      websocket.send(JSON.stringify({ type: 'pong', timestamp: Date.now() }));
    } else if (data.type === 'subscribe') {
      // Handle subscription
      websocket.send(JSON.stringify({
        type: 'subscribed',
        channels: data.channels || [],
      }));
    }
  }

  broadcast(message) {
    const payload = JSON.stringify(message);
    this.sessions.forEach((ws) => {
      ws.send(payload);
    });
  }
}

export default {
  fetch: handleRequest,
};

