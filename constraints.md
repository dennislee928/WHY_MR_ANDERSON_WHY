cloudlfare workers paid:
Workers Paid Plan


Includes first 10 million Workers and Pages Functions requests and 30 million CPU milliseconds

Includes 20 million Workers Logs observability events per month with 7 days retention

Includes first 5GB of D1 storage, 25 billion rows read, and 50 million rows written

Includes first 1GB of KV storage, 10 million read operations, and 1 million write, delete, and list operations

Includes first 1 million requests, 400,000 GB-second duration, 1 GB of stored data, 1 million read units, 1 million write units, and 1 million delete operations of Durable Objects

Includes first 10 million Workers Trace event logs

Includes first 1 million standard operations of Queues

Includes first 200,000 AI Gateway logs stored

Workers and Pages Functions


+ $0.30 per additional 1 million billable requests

+ $0.02 per additional million CPU milliseconds

D1


+ $0.75 per additional 1GB of storage

+ $0.001 per additional 1 million rows read

+ $1.00 per additional 1 million rows written

Durable Objects


+ $0.15 per additional 1 million requests

+ $12.50 per additional 1 million GB-second duration

+ $0.20 per additional 1GB of storage

+ $0.00 per additional 1 million rows read

+ $1.00 per additional 1 million rows written

KV Storage


+ $0.50 per additional 1GB of storage

+ $0.50 per additional 1 million read operations

+ $5.00 per additional 1 million write, delete, and list operations

Logpush

+ $0.05 per additional 1 million Workers Trace event logs

Queues


+ $0.40 per additional 1 million standard operations
___
Project Integration Strategy
a) Merge both projects into one unified system (recommended for simpler structure)
b) Keep them separate but standardize their structure and CI/CD
c) Focus only on Local_IPS-IDS (main project)
CI/CD Platform Choice
a) Buddy CI (easiest setup, good for Docker)
b) Argo CD (GitOps, best for Kubernetes)
c) Harness (enterprise-grade, complex)
d) Provide configs for all three
Primary Deployment Target
a) Cloudflare Workers (per constraints.md, serverless)
b) OCI Free Tier (VM/container based)
c) IBM Cloud Free Tier (VM/container based)
d) Multi-cloud strategy with all options
Documentation Language
a) Translate all Chinese docs to English only
b) Keep bilingual (English + Traditional Chinese)
Project Structure Simplification
a) Consolidate duplicate configs/scripts across projects
b) Flatten deep directory structures
c) Remove experimental/archive folders
d) All of the above