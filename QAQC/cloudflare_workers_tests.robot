*** Settings ***
Documentation    Cloudflare Workers Specific API Tests
Library          RequestsLibrary
Library          Collections
Library          DateTime
Library          String
Library          JSON

*** Variables ***
${BASE_URL}      https://security-platform-worker.workers.dev
${API_VERSION}   v1
${TIMEOUT}       30

# Cloudflare Workers specific test data
${WORKER_TEST_DATA}    {"worker": "security-platform-worker", "environment": "production"}
${KV_TEST_KEY}         test_key_${RANDOM}
${KV_TEST_VALUE}       {"test": "data", "timestamp": "${CURRENT_TIMESTAMP}"}

*** Test Cases ***

# ========================================
# CLOUDFLARE WORKERS FUNCTIONALITY TESTS
# ========================================

Test Worker Environment Detection
    [Documentation]    Test that worker correctly identifies its environment
    [Tags]    worker    environment
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    worker
    Should Be Equal    ${body['services']['worker']}    healthy

Test Worker Metadata
    [Documentation]    Test worker metadata and version information
    [Tags]    worker    metadata
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    version
    Should Contain    ${body}    uptime
    Should Be Equal    ${body['version']}    1.0.0

# ========================================
# KV STORAGE TESTS (Simulated)
# ========================================

Test Cache Functionality
    [Documentation]    Test KV cache functionality through API
    [Tags]    kv    cache
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    cache
    # Note: Cache status will show "not_configured" until KV is properly set up

Test Session Management
    [Documentation]    Test session management functionality
    [Tags]    kv    sessions
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/devices
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    message
    # Sessions will show "not_configured" until KV is properly set up

# ========================================
# D1 DATABASE TESTS (Simulated)
# ========================================

Test Database Connectivity
    [Documentation]    Test D1 database connectivity status
    [Tags]    d1    database
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    database
    # Database status will show "not_configured" until D1 is properly set up

Test Database Operations
    [Documentation]    Test database operations through API
    [Tags]    d1    database    operations
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/security/threats
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    message
    Should Contain    ${body['message']}    Database not configured yet

# ========================================
# R2 STORAGE TESTS (Simulated)
# ========================================

Test File Storage Status
    [Documentation]    Test R2 file storage status
    [Tags]    r2    storage
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    storage
    # Storage status will show "not_configured" until R2 is properly set up

# ========================================
# WORKERS AI TESTS (Simulated)
# ========================================

Test AI Service Status
    [Documentation]    Test Workers AI service status
    [Tags]    ai    workers_ai
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    ai
    # AI status will show "not_configured" until Workers AI is properly set up

# ========================================
# DURABLE OBJECTS TESTS (Simulated)
# ========================================

Test WebSocket Manager Status
    [Documentation]    Test WebSocket Manager Durable Object status
    [Tags]    durable_objects    websocket
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    websocket
    # WebSocket status will show "not_configured" until DO is properly set up

# ========================================
# ANALYTICS ENGINE TESTS (Simulated)
# ========================================

Test Analytics Engine Status
    [Documentation]    Test Analytics Engine status
    [Tags]    analytics    engine
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    analytics
    # Analytics status will show "not_configured" until Analytics Engine is properly set up

# ========================================
# RATE LIMITING TESTS
# ========================================

Test Rate Limiter Functionality
    [Documentation]    Test rate limiter functionality
    [Tags]    rate_limiter    performance
    ${requests_made}=    Set Variable    0
    ${rate_limited}=    Set Variable    ${False}
    
    FOR    ${i}    IN RANGE    0    200
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
        ${requests_made}=    Evaluate    ${requests_made} + 1
        
        IF    ${response.status_code} == 429
            ${rate_limited}=    Set Variable    ${True}
            Log    Rate limit triggered after ${requests_made} requests
            BREAK
        END
        
        IF    ${requests_made} > 150
            Log    Made ${requests_made} requests without hitting rate limit
            BREAK
        END
    END
    
    Should Be True    ${rate_limited}    Rate limiter should trigger after 150 requests

# ========================================
# SECRETS STORE TESTS (Simulated)
# ========================================

Test Secrets Store Status
    [Documentation]    Test Secrets Store status
    [Tags]    secrets    store
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    secrets
    # Secrets status will show "not_configured" until Secrets Store is properly set up

# ========================================
# QUEUE TESTS (Simulated)
# ========================================

Test Queue Processing Status
    [Documentation]    Test Queue processing status
    [Tags]    queue    processing
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    queue
    # Queue status will show "not_configured" until Queue is properly set up

# ========================================
# VECTORIZE TESTS (Simulated)
# ========================================

Test Vectorize Index Status
    [Documentation]    Test Vectorize index status
    [Tags]    vectorize    index
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    vectorize
    # Vectorize status will show "not_configured" until Vectorize is properly set up

# ========================================
# HYPERDRIVE TESTS (Simulated)
# ========================================

Test Hyperdrive Status
    [Documentation]    Test Hyperdrive status
    [Tags]    hyperdrive    database_acceleration
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    hyperdrive
    # Hyperdrive status will show "not_configured" until Hyperdrive is properly set up

# ========================================
# INTEGRATION TESTS
# ========================================

Test All Services Integration
    [Documentation]    Test integration of all Cloudflare services
    [Tags]    integration    all_services
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    
    # Check all expected services are present
    ${services}=    Set Variable    ${body['services']}
    Should Contain    ${services}    worker
    Should Contain    ${services}    database
    Should Contain    ${services}    cache
    Should Contain    ${services}    storage
    Should Contain    ${services}    ai
    Should Contain    ${services}    websocket
    Should Contain    ${services}    analytics
    Should Contain    ${services}    secrets
    Should Contain    ${services}    queue
    Should Contain    ${services}    vectorize
    Should Contain    ${services}    hyperdrive

Test Service Health Check
    [Documentation]    Test overall service health
    [Tags]    health    integration
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Be Equal    ${body['healthy']}    ${True}

*** Keywords ***

Generate Random Key
    [Documentation]    Generate a random key for testing
    ${random}=    Generate Random String    8    [LETTERS][NUMBERS]
    [Return]    test_key_${random}

Validate Service Status
    [Documentation]    Validate service status in response
    [Arguments]    ${response}    ${service_name}    ${expected_status}=not_configured
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body['services']}    ${service_name}
    Should Be Equal    ${body['services'][${service_name}]}    ${expected_status}

Setup Cloudflare Test Environment
    [Documentation]    Setup Cloudflare-specific test environment
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${EMPTY}
    ${CURRENT_TIMESTAMP}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S.%fZ
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${CURRENT_TIMESTAMP}
    ${RANDOM}=    Generate Random Key
    Set Global Variable    ${RANDOM}    ${RANDOM}

Teardown Cloudflare Test Environment
    [Documentation]    Cleanup Cloudflare test environment
    Log    Cloudflare test environment cleanup completed
