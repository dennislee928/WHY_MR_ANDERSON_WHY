# Database Container
FROM postgres:15-alpine

# Environment variables
ENV POSTGRES_DB=security_platform
ENV POSTGRES_USER=security_user
ENV POSTGRES_PASSWORD=secure_password_change_me

# Copy initialization scripts
COPY init-scripts/ /docker-entrypoint-initdb.d/

# Create custom configuration
COPY postgresql.conf /etc/postgresql/postgresql.conf

# Expose port
EXPOSE 5432

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD pg_isready -U $POSTGRES_USER -d $POSTGRES_DB || exit 1

# Use custom configuration
CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]
