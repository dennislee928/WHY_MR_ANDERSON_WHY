/**
 * Caching middleware for Cloudflare Workers
 * 
 * Implements response caching with TTL
 */

export function cacheMiddleware(ttl = 60) {
  return async (request, env, ctx) => {
    const cacheKey = new URL(request.url).pathname;
    
    // Try to get from cache
    const cached = await env.CACHE.get(cacheKey, 'json');
    if (cached) {
      return new Response(JSON.stringify(cached), {
        headers: {
          'Content-Type': 'application/json',
          'X-Cache': 'HIT',
        },
      });
    }

    // Not in cache, proceed with request
    return null; // Let router handle the request
  };
}

