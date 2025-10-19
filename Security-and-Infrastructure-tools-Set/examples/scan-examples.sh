#!/bin/bash
# ============================================
# Scan Examples
# ============================================
# This file provides practical scanning examples

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_example() {
    echo -e "${BLUE}$1${NC}"
    echo -e "${GREEN}$2${NC}"
    echo ""
}

echo "============================================"
echo "  Security Scanning Examples"
echo "============================================"
echo ""

# Nuclei Examples
echo "=== Nuclei Scanner ==="
echo ""

print_example "Example 1: Basic scan of a single target" \
    "make scan-nuclei TARGET=https://example.com"

print_example "Example 2: Scan with specific severity" \
    "docker-compose run --rm scanner-nuclei nuclei -u https://example.com -severity critical,high -o /results/critical.json"

print_example "Example 3: Scan with specific template tags" \
    "docker-compose run --rm scanner-nuclei nuclei -u https://example.com -tags cve,exposure,misconfig -o /results/tagged.json"

print_example "Example 4: Scan multiple targets from file" \
    "echo 'https://example1.com' > targets.txt
echo 'https://example2.com' >> targets.txt
docker-compose run --rm -v ./targets.txt:/targets.txt scanner-nuclei nuclei -l /targets.txt -o /results/batch.json"

print_example "Example 5: Scan with custom templates" \
    "docker-compose run --rm -v ./custom-templates:/custom scanner-nuclei nuclei -u https://example.com -t /custom -o /results/custom.json"

print_example "Example 6: Update templates" \
    "docker-compose run --rm scanner-nuclei nuclei -update-templates"

# Nmap Examples
echo "=== Nmap Scanner ==="
echo ""

print_example "Example 1: Basic port scan" \
    "make scan-nmap TARGET=192.168.1.1"

print_example "Example 2: Scan entire subnet" \
    "make scan-nmap TARGET=192.168.1.0/24"

print_example "Example 3: Service version detection" \
    "docker-compose run --rm nmap nmap -sV 192.168.1.1 -oX /results/nmap-version.xml"

print_example "Example 4: OS detection" \
    "docker-compose run --rm nmap nmap -O 192.168.1.1 -oX /results/nmap-os.xml"

print_example "Example 5: Full scan (aggressive)" \
    "docker-compose run --rm nmap nmap -A -T4 192.168.1.1 -oX /results/nmap-full.xml"

print_example "Example 6: Scan specific ports" \
    "docker-compose run --rm nmap nmap -p 80,443,8080 192.168.1.1 -oX /results/nmap-ports.xml"

print_example "Example 7: NSE script scan" \
    "docker-compose run --rm nmap nmap --script vuln 192.168.1.1 -oX /results/nmap-vuln.xml"

# AMASS Examples
echo "=== AMASS Scanner ==="
echo ""

print_example "Example 1: Subdomain enumeration" \
    "docker-compose run --rm scanner-amass amass enum -d example.com -o /results/amass-subs.txt"

print_example "Example 2: Passive mode (no active probing)" \
    "docker-compose run --rm scanner-amass amass enum -passive -d example.com -o /results/amass-passive.txt"

print_example "Example 3: With configuration file (API keys)" \
    "docker-compose run --rm -v ./amass-config.ini:/config.ini scanner-amass amass enum -config /config.ini -d example.com"

print_example "Example 4: JSON output" \
    "docker-compose run --rm scanner-amass amass enum -d example.com -json /results/amass.json"

# Advanced Workflows
echo "=== Advanced Workflows ==="
echo ""

print_example "Example 1: Complete asset discovery and scanning" \
    "# Step 1: Discover subdomains
docker-compose run --rm scanner-amass amass enum -d example.com -o /results/subs.txt

# Step 2: Scan each subdomain
while read sub; do
    make scan-nuclei TARGET=https://\$sub
done < /results/subs.txt"

print_example "Example 2: Scheduled daily scan (cron)" \
    "# Add to crontab
0 2 * * * cd /path/to/project/Make_Files && make scan-nuclei TARGET=https://example.com"

print_example "Example 3: Scan and alert on critical findings" \
    "# Run scan
make scan-nuclei TARGET=https://example.com

# Check for critical findings
CRITICAL=\$(docker exec postgres psql -U sectools -d security -t -c \"SELECT COUNT(*) FROM scan_findings WHERE severity='critical' AND discovered_at > NOW() - INTERVAL '1 day'\")

if [ \$CRITICAL -gt 0 ]; then
    echo \"Found \$CRITICAL critical vulnerabilities!\"
    # Send alert (e.g., via Slack webhook)
    curl -X POST -H 'Content-type: application/json' --data '{\"text\":\"Found '\$CRITICAL' critical vulns!\"}' \$SLACK_WEBHOOK_URL
fi"

echo "============================================"
echo "For more examples, see docs/ and README.md"
echo "============================================"

