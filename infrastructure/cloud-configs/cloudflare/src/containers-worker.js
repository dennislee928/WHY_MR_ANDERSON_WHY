/**
 * Security Platform Containers Worker
 * 
 * This worker orchestrates multiple containers for the security platform:
 * - Backend API Container
 * - AI/Quantum Processing Container
 * - Security Tools Container
 * - Database Container
 * - Monitoring Container
 * 
 * Features:
 * - Container orchestration
 * - Load balancing between containers
 * - Health monitoring
 * - Auto-scaling
 * - Service discovery
 */

import { Router } from 'itty-router';

// Create router
const router = Router();

// Container registry for service discovery
const CONTAINER_REGISTRY = {
  backend: 'BACKEND_API',
  ai: 'AI_QUANTUM',
  security: 'SECURITY_TOOLS',
  database: 'DATABASE',
  monitoring: 'MONITORING'
};

// Container health status
const containerHealth = new Map();

/**
 * Container Health Check
 */
async function checkContainerHealth(env, containerName) {
  try {
    const container = env[CONTAINER_REGISTRY[containerName]];
    if (!container) {
      return { healthy: false, error: 'Container not found' };
    }

    const response = await container.fetch('/health', {
      method: 'GET',
      timeout: 5000
    });

    if (response.ok) {
      const health = await response.json();
      containerHealth.set(containerName, {
        healthy: true,
        lastCheck: Date.now(),
        ...health
      });
      return { healthy: true, ...health };
    } else {
      containerHealth.set(containerName, {
        healthy: false,
        lastCheck: Date.now(),
        error: `HTTP ${response.status}`
      });
      return { healthy: false, error: `HTTP ${response.status}` };
    }
  } catch (error) {
    containerHealth.set(containerName, {
      healthy: false,
      lastCheck: Date.now(),
      error: error.message
    });
    return { healthy: false, error: error.message };
  }
}

/**
 * Load Balancer for Container Requests
 */
async function loadBalanceRequest(env, serviceName, path, options = {}) {
  const containerName = CONTAINER_REGISTRY[serviceName];
  const container = env[containerName];
  
  if (!container) {
    throw new Error(`Container ${serviceName} not found`);
  }

  // Check container health before routing
  const health = await checkContainerHealth(env, serviceName);
  if (!health.healthy) {
    throw new Error(`Container ${serviceName} is unhealthy: ${health.error}`);
  }

  // Route request to container
  const url = new URL(path, `https://${serviceName}.internal`);
  return await container.fetch(url.toString(), {
    ...options,
    headers: {
      'X-Forwarded-For': options.headers?.['CF-Connecting-IP'] || 'unknown',
      'X-Service-Name': serviceName,
      ...options.headers
    }
  });
}

/**
 * Container Orchestration Routes
 */

// Health check for all containers
router.get('/api/v1/containers/health', async (request, env) => {
  const healthChecks = {};
  
  for (const [serviceName] of Object.entries(CONTAINER_REGISTRY)) {
    healthChecks[serviceName] = await checkContainerHealth(env, serviceName);
  }

  const allHealthy = Object.values(healthChecks).every(h => h.healthy);
  
  return new Response(JSON.stringify({
    overall: allHealthy ? 'healthy' : 'degraded',
    containers: healthChecks,
    timestamp: new Date().toISOString()
  }), {
    headers: { 'Content-Type': 'application/json' }
  });
});

// Service discovery endpoint
router.get('/api/v1/services', async (request, env) => {
  const services = {};
  
  for (const [serviceName, containerBinding] of Object.entries(CONTAINER_REGISTRY)) {
    const health = containerHealth.get(serviceName) || { healthy: false };
    services[serviceName] = {
      binding: containerBinding,
      healthy: health.healthy,
      lastCheck: health.lastCheck,
      endpoints: [
        `/${serviceName}/api/v1/*`,
        `/${serviceName}/health`,
        `/${serviceName}/metrics`
      ]
    };
  }

  return new Response(JSON.stringify({
    services,
    timestamp: new Date().toISOString()
  }), {
    headers: { 'Content-Type': 'application/json' }
  });
});

// Backend API routes
router.all('/api/v1/backend/*', async (request, env) => {
  try {
    const path = request.url.replace(/.*\/api\/v1\/backend/, '/api/v1');
    const response = await loadBalanceRequest(env, 'backend', path, {
      method: request.method,
      headers: request.headers,
      body: request.body
    });
    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Backend service unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// AI/Quantum routes
router.all('/api/v1/ai/*', async (request, env) => {
  try {
    const path = request.url.replace(/.*\/api\/v1\/ai/, '/api/v1');
    const response = await loadBalanceRequest(env, 'ai', path, {
      method: request.method,
      headers: request.headers,
      body: request.body
    });
    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'AI service unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Security Tools routes
router.all('/api/v1/security/*', async (request, env) => {
  try {
    const path = request.url.replace(/.*\/api\/v1\/security/, '/api/v1');
    const response = await loadBalanceRequest(env, 'security', path, {
      method: request.method,
      headers: request.headers,
      body: request.body
    });
    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Security service unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Database routes
router.all('/api/v1/database/*', async (request, env) => {
  try {
    const path = request.url.replace(/.*\/api\/v1\/database/, '/api/v1');
    const response = await loadBalanceRequest(env, 'database', path, {
      method: request.method,
      headers: request.headers,
      body: request.body
    });
    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Database service unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Monitoring routes
router.all('/api/v1/monitoring/*', async (request, env) => {
  try {
    const path = request.url.replace(/.*\/api\/v1\/monitoring/, '/api/v1');
    const response = await loadBalanceRequest(env, 'monitoring', path, {
      method: request.method,
      headers: request.headers,
      body: request.body
    });
    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Monitoring service unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Container management routes
router.post('/api/v1/containers/:serviceName/scale', async (request, env) => {
  try {
    const { serviceName } = request.params;
    const { replicas } = await request.json();
    
    if (!CONTAINER_REGISTRY[serviceName]) {
      return new Response(JSON.stringify({
        error: 'Service not found'
      }), {
        status: 404,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    const response = await loadBalanceRequest(env, serviceName, '/api/v1/scale', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ replicas })
    });

    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Scaling failed',
      message: error.message
    }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Container logs
router.get('/api/v1/containers/:serviceName/logs', async (request, env) => {
  try {
    const { serviceName } = request.params;
    
    if (!CONTAINER_REGISTRY[serviceName]) {
      return new Response(JSON.stringify({
        error: 'Service not found'
      }), {
        status: 404,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    const response = await loadBalanceRequest(env, serviceName, '/api/v1/logs', {
      method: 'GET',
      headers: request.headers
    });

    return response;
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Logs unavailable',
      message: error.message
    }), {
      status: 503,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Metrics aggregation
router.get('/api/v1/metrics', async (request, env) => {
  try {
    const metrics = {};
    
    for (const [serviceName] of Object.entries(CONTAINER_REGISTRY)) {
      try {
        const response = await loadBalanceRequest(env, serviceName, '/api/v1/metrics', {
          method: 'GET'
        });
        if (response.ok) {
          metrics[serviceName] = await response.json();
        }
      } catch (error) {
        metrics[serviceName] = { error: error.message };
      }
    }

    return new Response(JSON.stringify({
      aggregated_metrics: metrics,
      timestamp: new Date().toISOString()
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  } catch (error) {
    return new Response(JSON.stringify({
      error: 'Metrics aggregation failed',
      message: error.message
    }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});

// Main request handler
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

  // Route requests
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
    return new Response(JSON.stringify({
      error: 'Internal Server Error',
      message: error.message,
    }), {
      status: 500,
      headers: {
        ...corsHeaders,
        'Content-Type': 'application/json',
      },
    });
  }
}

export default {
  fetch: handleRequest,
};
