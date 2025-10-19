*** Settings ***
Documentation    Security Platform API Tests - Cloudflare Workers
Library          RequestsLibrary
Library          Collections
Library          DateTime
Library          String
Library          JSON

*** Variables ***
${BASE_URL}      https://security-platform-worker.workers.dev
${API_VERSION}   v1
${TIMEOUT}       30
${RETRY_COUNT}   3

# Test Data
${VALID_THREAT_DATA}    {"type": "malware", "severity": "high", "source": "test"}
${INVALID_THREAT_DATA}  {"invalid": "data"}

*** Test Cases ***

# ========================================
# HEALTH CHECK TESTS
# ========================================

Test Health Check Endpoint
    [Documentation]    Test the basic health check endpoint
    [Tags]    health    smoke
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Be Equal    ${body['healthy']}    ${True}
    Should Contain    ${body}    timestamp
    Should Contain    ${body}    version

Test Status Endpoint
    [Documentation]    Test the system status endpoint
    [Tags]    status    smoke
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Be Equal    ${body['status']}    operational
    Should Contain    ${body}    services
    Should Contain    ${body['services']}    worker

# ========================================
# SECURITY ENDPOINT TESTS
# ========================================

Test Security Threats Endpoint
    [Documentation]    Test the security threats listing endpoint
    [Tags]    security    threats
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/security/threats
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    threats
    Should Contain    ${body}    total
    Should Contain    ${body}    message

Test Security Threats with Parameters
    [Documentation]    Test security threats endpoint with query parameters
    [Tags]    security    threats    parameters
    ${params}=    Create Dictionary    severity=high    limit=10
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/security/threats    params=${params}
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    threats

# ========================================
# NETWORK ENDPOINT TESTS
# ========================================

Test Network Stats Endpoint
    [Documentation]    Test the network statistics endpoint
    [Tags]    network    stats
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/network/stats
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    totalRequests
    Should Contain    ${body}    blockedRequests
    Should Contain    ${body}    activeConnections
    Should Contain    ${body}    message

Test Network Stats with Time Range
    [Documentation]    Test network stats with time range parameters
    [Tags]    network    stats    parameters
    ${params}=    Create Dictionary    start_time=${CURRENT_TIMESTAMP}    end_time=${CURRENT_TIMESTAMP}
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/network/stats    params=${params}
    Should Be Equal As Strings    ${response.status_code}    200

# ========================================
# DEVICE ENDPOINT TESTS
# ========================================

Test Devices Endpoint
    [Documentation]    Test the devices listing endpoint
    [Tags]    devices    management
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/devices
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    devices
    Should Contain    ${body}    total
    Should Contain    ${body}    message

Test Devices with Filter
    [Documentation]    Test devices endpoint with filtering parameters
    [Tags]    devices    filter
    ${params}=    Create Dictionary    status=active    type=endpoint
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/devices    params=${params}
    Should Be Equal As Strings    ${response.status_code}    200

# ========================================
# ERROR HANDLING TESTS
# ========================================

Test Non-existent Endpoint
    [Documentation]    Test handling of non-existent endpoints
    [Tags]    error    not_found
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/nonexistent
    Should Be Equal As Strings    ${response.status_code}    404

Test Invalid API Version
    [Documentation]    Test handling of invalid API version
    [Tags]    error    invalid_version
    ${response}=    GET    ${BASE_URL}/api/invalid/health
    Should Be Equal As Strings    ${response.status_code}    404

# ========================================
# RATE LIMITING TESTS
# ========================================

Test Rate Limiting
    [Documentation]    Test API rate limiting functionality
    [Tags]    rate_limit    performance
    FOR    ${i}    IN RANGE    0    160
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
        IF    ${response.status_code} == 429
            Log    Rate limit triggered after ${i} requests
            BREAK
        END
    END
    # Should eventually hit rate limit (150 requests per minute)

# ========================================
# CORS TESTS
# ========================================

Test CORS Preflight Request
    [Documentation]    Test CORS preflight request handling
    [Tags]    cors    preflight
    ${headers}=    Create Dictionary    Origin=https://example.com
    ${response}=    OPTIONS    ${BASE_URL}/api/${API_VERSION}/health    headers=${headers}
    Should Be Equal As Strings    ${response.status_code}    200
    Should Contain    ${response.headers}    Access-Control-Allow-Origin

Test CORS Headers in Response
    [Documentation]    Test CORS headers in API responses
    [Tags]    cors    headers
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    Should Be Equal As Strings    ${response.status_code}    200
    Should Contain    ${response.headers}    Access-Control-Allow-Origin
    Should Contain    ${response.headers}    Access-Control-Allow-Methods

# ========================================
# PERFORMANCE TESTS
# ========================================

Test Response Time
    [Documentation]    Test API response time performance
    [Tags]    performance    response_time
    ${start_time}=    Get Current Date    result_format=epoch
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    ${end_time}=    Get Current Date    result_format=epoch
    ${response_time}=    Evaluate    ${end_time} - ${start_time}
    Should Be True    ${response_time} < 2.0    Response time should be under 2 seconds
    Should Be Equal As Strings    ${response.status_code}    200

Test Concurrent Requests
    [Documentation]    Test handling of concurrent requests
    [Tags]    performance    concurrent
    ${responses}=    Create List
    FOR    ${i}    IN RANGE    0    10
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
        Append To List    ${responses}    ${response.status_code}
    END
    FOR    ${status_code}    IN    @{responses}
        Should Be Equal As Strings    ${status_code}    200
    END

# ========================================
# DATA VALIDATION TESTS
# ========================================

Test JSON Response Format
    [Documentation]    Test that all responses return valid JSON
    [Tags]    validation    json
    ${endpoints}=    Create List    health    status    security/threats    network/stats    devices
    FOR    ${endpoint}    IN    @{endpoints}
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/${endpoint}
        Should Be Equal As Strings    ${response.status_code}    200
        ${body}=    Set Variable    ${response.json()}
        Should Not Be Empty    ${body}
    END

Test Response Schema Validation
    [Documentation]    Test response schema validation
    [Tags]    validation    schema
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    healthy
    Should Contain    ${body}    timestamp
    Should Contain    ${body}    version
    Should Contain    ${body}    uptime

*** Keywords ***

Get Current Timestamp
    [Documentation]    Get current timestamp in ISO format
    ${timestamp}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S.%fZ
    [Return]    ${timestamp}

Validate API Response
    [Documentation]    Validate basic API response structure
    [Arguments]    ${response}    ${expected_status}=200
    Should Be Equal As Strings    ${response.status_code}    ${expected_status}
    Should Not Be Empty    ${response.content}
    ${body}=    Set Variable    ${response.json()}
    Should Not Be Empty    ${body}

Setup Test Environment
    [Documentation]    Setup test environment and variables
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${EMPTY}
    ${CURRENT_TIMESTAMP}=    Get Current Timestamp
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${CURRENT_TIMESTAMP}

Teardown Test Environment
    [Documentation]    Cleanup test environment
    Log    Test environment cleanup completed
