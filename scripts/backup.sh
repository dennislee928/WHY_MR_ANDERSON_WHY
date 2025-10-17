#!/bin/bash
# ============================================
# Database Backup Script
# ============================================
# This script creates encrypted backups of the PostgreSQL database
# and optionally uploads them to cloud storage

set -e  # Exit on error

# Configuration
DATE=$(date +%Y%m%d-%H%M%S)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKUP_DIR="${BACKUP_DIR:-$PROJECT_ROOT/backups}"
RETENTION_DAYS="${BACKUP_RETENTION_DAYS:-30}"
DB_HOST="${DB_HOST:-postgres}"
DB_USER="${DB_USER:-sectools}"
DB_NAME="${DB_NAME:-security}"
ENCRYPT="${ENCRYPT_BACKUP:-false}"
UPLOAD_S3="${UPLOAD_S3:-false}"
S3_BUCKET="${S3_BUCKET:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if Docker is running
    if ! docker info &>/dev/null; then
        log_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    # Check if postgres container exists
    if ! docker ps | grep -q postgres; then
        log_error "PostgreSQL container is not running."
        exit 1
    fi
    
    # Create backup directory if not exists
    mkdir -p "$BACKUP_DIR"
    
    log_info "Prerequisites check passed ✓"
}

create_backup() {
    local backup_file="$BACKUP_DIR/security-$DATE.sql"
    
    log_info "Creating database backup..."
    log_info "Backup file: $backup_file"
    
    # Use docker-compose exec to run pg_dump
    cd "$PROJECT_ROOT/Docker/compose" || exit 1
    
    if docker-compose exec -T postgres pg_dump -U "$DB_USER" -d "$DB_NAME" > "$backup_file"; then
        log_info "Backup created successfully ✓"
        log_info "Backup size: $(du -h "$backup_file" | cut -f1)"
        echo "$backup_file"
    else
        log_error "Backup failed!"
        exit 1
    fi
}

encrypt_backup() {
    local backup_file="$1"
    local encrypted_file="${backup_file}.gpg"
    
    if [ "$ENCRYPT" = "true" ]; then
        log_info "Encrypting backup..."
        
        if command -v gpg &>/dev/null; then
            # Use passphrase from environment or prompt
            if [ -n "${BACKUP_PASSPHRASE:-}" ]; then
                echo "$BACKUP_PASSPHRASE" | gpg --batch --yes --passphrase-fd 0 \
                    --symmetric --cipher-algo AES256 --output "$encrypted_file" "$backup_file"
            else
                gpg --symmetric --cipher-algo AES256 --output "$encrypted_file" "$backup_file"
            fi
            
            # Remove unencrypted file
            rm -f "$backup_file"
            log_info "Backup encrypted ✓"
            echo "$encrypted_file"
        else
            log_warn "GPG not found. Skipping encryption."
            echo "$backup_file"
        fi
    else
        echo "$backup_file"
    fi
}

compress_backup() {
    local backup_file="$1"
    local compressed_file="${backup_file}.gz"
    
    log_info "Compressing backup..."
    
    if gzip -c "$backup_file" > "$compressed_file"; then
        rm -f "$backup_file"
        log_info "Backup compressed ✓"
        log_info "Compressed size: $(du -h "$compressed_file" | cut -f1)"
        echo "$compressed_file"
    else
        log_warn "Compression failed. Keeping original file."
        echo "$backup_file"
    fi
}

upload_to_s3() {
    local backup_file="$1"
    
    if [ "$UPLOAD_S3" = "true" ] && [ -n "$S3_BUCKET" ]; then
        log_info "Uploading to S3..."
        
        if command -v aws &>/dev/null; then
            if aws s3 cp "$backup_file" "s3://$S3_BUCKET/backups/$(basename "$backup_file")"; then
                log_info "Uploaded to S3 ✓"
            else
                log_error "S3 upload failed!"
                return 1
            fi
        else
            log_warn "AWS CLI not found. Skipping S3 upload."
        fi
    fi
}

cleanup_old_backups() {
    log_info "Cleaning up old backups (older than $RETENTION_DAYS days)..."
    
    local count=0
    while IFS= read -r -d '' file; do
        rm -f "$file"
        ((count++))
    done < <(find "$BACKUP_DIR" -name "security-*.sql*" -mtime +"$RETENTION_DAYS" -print0)
    
    if [ "$count" -gt 0 ]; then
        log_info "Removed $count old backup(s) ✓"
    else
        log_info "No old backups to remove"
    fi
}

verify_backup() {
    local backup_file="$1"
    
    log_info "Verifying backup integrity..."
    
    # Check if file exists and is not empty
    if [ -f "$backup_file" ] && [ -s "$backup_file" ]; then
        log_info "Backup verification passed ✓"
        return 0
    else
        log_error "Backup verification failed!"
        return 1
    fi
}

# Main execution
main() {
    log_info "=========================================="
    log_info "PostgreSQL Backup Script"
    log_info "=========================================="
    log_info "Started at: $(date)"
    
    check_prerequisites
    
    # Create backup
    BACKUP_FILE=$(create_backup)
    
    # Verify backup
    if ! verify_backup "$BACKUP_FILE"; then
        exit 1
    fi
    
    # Compress backup
    BACKUP_FILE=$(compress_backup "$BACKUP_FILE")
    
    # Encrypt if enabled
    BACKUP_FILE=$(encrypt_backup "$BACKUP_FILE")
    
    # Upload to S3 if enabled
    upload_to_s3 "$BACKUP_FILE"
    
    # Cleanup old backups
    cleanup_old_backups
    
    log_info "=========================================="
    log_info "Backup completed successfully!"
    log_info "Final backup: $BACKUP_FILE"
    log_info "Completed at: $(date)"
    log_info "=========================================="
}

# Run main function
main "$@"

