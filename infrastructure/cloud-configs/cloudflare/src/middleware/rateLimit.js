/**
 * Rate limiting middleware for Cloudflare Workers
 * Uses KV for distributed rate limiting
 */

export class RateLimiter {
  constructor(options = {}) {
    this.limit = options.limit || 100;
    this.window = options.window || 60 * 1000; // 1 minute default
  }

  async check(request, env) {
    const ip = request.headers.get('CF-Connecting-IP') || 'unknown';
    const key = `ratelimit:${ip}`;

    // Get current count from KV
    const data = await env.CACHE.get(key, 'json');
    const now = Date.now();

    if (!data) {
      // First request
      await env.CACHE.put(
        key,
        JSON.stringify({
          count: 1,
          resetAt: now + this.window,
        }),
        {
          expirationTtl: Math.ceil(this.window / 1000),
        }
      );

      return { allowed: true };
    }

    // Check if window has expired
    if (now > data.resetAt) {
      await env.CACHE.put(
        key,
        JSON.stringify({
          count: 1,
          resetAt: now + this.window,
        }),
        {
          expirationTtl: Math.ceil(this.window / 1000),
        }
      );

      return { allowed: true };
    }

    // Check if limit exceeded
    if (data.count >= this.limit) {
      return {
        allowed: false,
        retryAfter: Math.ceil((data.resetAt - now) / 1000),
      };
    }

    // Increment count
    await env.CACHE.put(
      key,
      JSON.stringify({
        count: data.count + 1,
        resetAt: data.resetAt,
      }),
      {
        expirationTtl: Math.ceil((data.resetAt - now) / 1000),
      }
    );

    return { allowed: true };
  }
}

