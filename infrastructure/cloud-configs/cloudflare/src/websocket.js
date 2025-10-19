/**
 * WebSocket handler for Cloudflare Workers
 * 
 * Uses Durable Objects for connection management
 */

export async function handleWebSocket(request, env) {
  const url = new URL(request.url);
  const clientId = url.searchParams.get('client_id') || crypto.randomUUID();

  // Get Durable Object stub
  const id = env.WEBSOCKET_MANAGER.idFromName(clientId);
  const stub = env.WEBSOCKET_MANAGER.get(id);

  // Forward request to Durable Object
  return stub.fetch(request);
}

