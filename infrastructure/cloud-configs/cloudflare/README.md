# Cloudflare Workers Deployment

Deploy the Security Platform to Cloudflare Workers using the free tier.

## Free Tier Limits

Based on `constraints.md`, the Cloudflare Workers Paid plan includes:

- **10 million** Workers requests/month
- **30 million** CPU milliseconds
- **5GB** D1 database storage
- **1GB** KV storage
- **1 million** Durable Objects requests

## Prerequisites

1. Cloudflare account
2. Wrangler CLI installed:
   ```bash
   npm install -g wrangler
   ```
3. Authentication:
   ```bash
   wrangler login
   ```

## Setup

### 1. Create D1 Database

```bash
cd infrastructure/cloud-configs/cloudflare
npm install
npm run d1:create
```

Note the database ID and add it to `wrangler.toml`:
```toml
[[d1_databases]]
database_id = "your-database-id-here"
```

### 2. Initialize Database Schema

```bash
npm run d1:init
```

### 3. Create KV Namespaces

```bash
npm run kv:create:cache
npm run kv:create:sessions
```

Add the namespace IDs to `wrangler.toml`:
```toml
[[kv_namespaces]]
binding = "CACHE"
id = "your-cache-namespace-id"

[[kv_namespaces]]
binding = "SESSIONS"
id = "your-sessions-namespace-id"
```

### 4. Update Configuration

Edit `wrangler.toml` and add your account ID:
```toml
account_id = "your-cloudflare-account-id"
```

## Deployment

### Development

```bash
npm run dev
```

Access at: http://localhost:8787

### Staging

```bash
npm run deploy:staging
```

### Production

```bash
npm run deploy:production
```

## API Endpoints

Once deployed, your Worker will handle:

### REST API
- `GET /api/v1/status` - System status
- `GET /api/v1/health` - Health check
- `GET /api/v1/security/threats` - List threats
- `POST /api/v1/security/threats/:id/block` - Block threat
- `GET /api/v1/network/stats` - Network statistics
- `DELETE /api/v1/network/blocked-ips/:ip` - Unblock IP
- `GET /api/v1/devices` - List devices
- `POST /api/v1/ml/detect` - ML threat detection

### WebSocket
- `ws://your-worker.workers.dev/ws?client_id=xxx` - Real-time updates

## Database Queries

Query your D1 database:

```bash
# List all threats
wrangler d1 execute security_platform_db --command "SELECT * FROM threats"

# Count active threats
wrangler d1 execute security_platform_db --command "SELECT COUNT(*) FROM threats WHERE status='active'"

# View recent API logs
wrangler d1 execute security_platform_db --command "SELECT * FROM api_logs ORDER BY timestamp DESC LIMIT 10"
```

## Monitoring

### View Logs

```bash
npm run tail
```

### Analytics

Visit Cloudflare Dashboard > Workers & Pages > your-worker > Analytics

Monitor:
- Request count
- CPU time used
- Error rate
- Response time

## Cost Management

### Stay Within Free Tier

1. **Request Limit (10M/month)**
   - ~13,888 requests/hour average
   - Implement caching (already included)
   - Monitor in Dashboard

2. **CPU Time (30M milliseconds)**
   - Keep handlers lightweight
   - Use caching for expensive operations
   - Offload heavy processing to external services

3. **D1 Storage (5GB)**
   - Implement data retention policies
   - Archive old records
   - Clean up test data

### Monitor Usage

```bash
# Check current usage
wrangler dash
```

Or visit: https://dash.cloudflare.com/

## Scaling Beyond Free Tier

If you exceed free tier limits:

1. **Cloudflare Workers Paid**: $5/month
   - Additional requests: $0.30 per 1M
   - Additional CPU: $0.02 per million milliseconds

2. **Optimization Strategies**:
   - Implement aggressive caching
   - Use Cloudflare Cache API
   - Reduce D1 queries with KV caching
   - Batch operations

## Troubleshooting

### Worker Not Responding

```bash
# Check deployment status
wrangler deployments list

# View recent errors
wrangler tail --format=pretty
```

### Database Errors

```bash
# Verify database exists
wrangler d1 list

# Test connection
wrangler d1 execute security_platform_db --command "SELECT 1"
```

### KV Issues

```bash
# List namespaces
wrangler kv:namespace list

# Get a key
wrangler kv:key get --namespace-id=xxx "your-key"
```

## Custom Domain

### Add Custom Domain

1. Go to Cloudflare Dashboard > Workers & Pages
2. Select your Worker
3. Click "Triggers" > "Custom Domains"
4. Add your domain (must be on Cloudflare)

### Update wrangler.toml

```toml
routes = [
  { pattern = "api.example.com/*", zone_name = "example.com" }
]
```

## Environment Variables

Set secrets:

```bash
# Set API keys
wrangler secret put IBM_QUANTUM_TOKEN
wrangler secret put DATABASE_ENCRYPTION_KEY

# List secrets
wrangler secret list
```

## Rollback

```bash
# List deployments
wrangler deployments list

# Rollback to previous
wrangler rollback
```

## Performance Tips

1. **Caching**
   - Use KV for frequently accessed data
   - Implement TTL-based cache invalidation
   - Cache API responses

2. **Rate Limiting**
   - Already implemented (150 req/min per IP)
   - Adjust in `src/middleware/rateLimit.js`

3. **Optimize Queries**
   - Use indexes (already in schema.sql)
   - Limit result sets
   - Use prepared statements

## Architecture

```
┌──────────────────────────────────────────────┐
│         Cloudflare Global Network            │
└──────────────────┬───────────────────────────┘
                   │
        ┌──────────▼──────────┐
        │   Workers Runtime    │
        │   (Edge Compute)     │
        └──────────┬───────────┘
                   │
     ┌─────────────┼─────────────┐
     │             │             │
┌────▼────┐  ┌────▼────┐  ┌─────▼──────┐
│   D1    │  │   KV    │  │  Durable   │
│Database │  │ Storage │  │  Objects   │
└─────────┘  └─────────┘  └────────────┘
```

## Security

- All connections use HTTPS by default
- Rate limiting prevents abuse
- D1 data encrypted at rest
- Secrets managed via Wrangler

## Further Reading

- [Cloudflare Workers Docs](https://developers.cloudflare.com/workers/)
- [D1 Database](https://developers.cloudflare.com/d1/)
- [Workers KV](https://developers.cloudflare.com/workers-kv/)
- [Durable Objects](https://developers.cloudflare.com/workers/runtime-apis/durable-objects/)

