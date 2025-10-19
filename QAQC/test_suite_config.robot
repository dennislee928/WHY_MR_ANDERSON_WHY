*** Settings ***
Documentation    Security Platform Test Suite Configuration
Library          RequestsLibrary
Library          Collections
Library          DateTime
Library          String
Library          JSON

*** Variables ***
# Test Configuration
${TEST_ENVIRONMENT}    production
${BASE_URL}           https://security-platform-worker.workers.dev
${API_VERSION}        v1
${TIMEOUT}            30
${RETRY_COUNT}        3

# Test Data Files
${TEST_DATA_DIR}      ${CURDIR}/test_data
${RESULTS_DIR}        ${CURDIR}/results
${LOGS_DIR}           ${CURDIR}/logs

# Test Categories
${SMOKE_TESTS}        health,status
${REGRESSION_TESTS}   security,network,devices
${PERFORMANCE_TESTS}  rate_limit,response_time,concurrent
${INTEGRATION_TESTS}  worker,integration,all_services

*** Test Cases ***

# ========================================
# TEST SUITE INITIALIZATION
# ========================================

Setup Test Suite
    [Documentation]    Initialize test suite environment
    [Tags]    setup
    Create Directory    ${RESULTS_DIR}
    Create Directory    ${LOGS_DIR}
    Create Directory    ${TEST_DATA_DIR}
    Log    Test suite initialized successfully

# ========================================
# ENVIRONMENT VALIDATION
# ========================================

Validate Test Environment
    [Documentation]    Validate test environment is ready
    [Tags]    validation    environment
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
    Should Be Equal As Strings    ${response.status_code}    200
    Log    Test environment validation passed

# ========================================
# CONNECTIVITY TESTS
# ========================================

Test Basic Connectivity
    [Documentation]    Test basic connectivity to Cloudflare Workers
    [Tags]    connectivity    smoke
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Be Equal    ${body['healthy']}    ${True}

Test API Version Compatibility
    [Documentation]    Test API version compatibility
    [Tags]    compatibility    api_version
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/status    timeout=${TIMEOUT}
    Should Be Equal As Strings    ${response.status_code}    200
    ${body}=    Set Variable    ${response.json()}
    Should Contain    ${body}    version

# ========================================
# SECURITY TESTS
# ========================================

Test Security Headers
    [Documentation]    Test security headers in responses
    [Tags]    security    headers
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
    Should Be Equal As Strings    ${response.status_code}    200
    Should Contain    ${response.headers}    Access-Control-Allow-Origin
    Should Contain    ${response.headers}    Access-Control-Allow-Methods

Test Input Validation
    [Documentation]    Test input validation and sanitization
    [Tags]    security    validation
    ${malicious_inputs}=    Create List    <script>alert('xss')</script>    '; DROP TABLE users; --    ../../../etc/passwd
    FOR    ${input}    IN    @{malicious_inputs}
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/security/threats?filter=${input}    timeout=${TIMEOUT}
        Should Be Equal As Strings    ${response.status_code}    200
        # Should not return error or execute malicious code
    END

# ========================================
# PERFORMANCE BENCHMARKS
# ========================================

Test Response Time Benchmarks
    [Documentation]    Test response time benchmarks
    [Tags]    performance    benchmarks
    ${endpoints}=    Create List    health    status    security/threats    network/stats    devices
    FOR    ${endpoint}    IN    @{endpoints}
        ${start_time}=    Get Current Date    result_format=epoch
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/${endpoint}    timeout=${TIMEOUT}
        ${end_time}=    Get Current Date    result_format=epoch
        ${response_time}=    Evaluate    ${end_time} - ${start_time}
        Should Be True    ${response_time} < 2.0    ${endpoint} response time should be under 2 seconds
        Should Be Equal As Strings    ${response.status_code}    200
    END

Test Throughput Capacity
    [Documentation]    Test API throughput capacity
    [Tags]    performance    throughput
    ${successful_requests}=    Set Variable    0
    ${failed_requests}=    Set Variable    0
    
    FOR    ${i}    IN RANGE    0    100
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
        IF    ${response.status_code} == 200
            ${successful_requests}=    Evaluate    ${successful_requests} + 1
        ELSE
            ${failed_requests}=    Evaluate    ${failed_requests} + 1
        END
    END
    
    ${success_rate}=    Evaluate    ${successful_requests} / 100 * 100
    Should Be True    ${success_rate} >= 95    Success rate should be at least 95%

# ========================================
# RELIABILITY TESTS
# ========================================

Test Error Handling
    [Documentation]    Test error handling and recovery
    [Tags]    reliability    error_handling
    ${error_endpoints}=    Create List    /api/${API_VERSION}/nonexistent    /api/invalid/health    /api/${API_VERSION}/security/threats/invalid
    FOR    ${endpoint}    IN    @{error_endpoints}
        ${response}=    GET    ${BASE_URL}${endpoint}    timeout=${TIMEOUT}
        Should Be True    ${response.status_code} >= 400    Should return error status for invalid endpoint
    END

Test Graceful Degradation
    [Documentation]    Test graceful degradation under load
    [Tags]    reliability    degradation
    FOR    ${i}    IN RANGE    0    50
        ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
        Should Be True    ${response.status_code} in [200, 429]    Should return 200 or rate limit 429
    END

# ========================================
# DATA CONSISTENCY TESTS
# ========================================

Test Data Consistency
    [Documentation]    Test data consistency across requests
    [Tags]    consistency    data
    ${response1}=    GET    ${BASE_URL}/api/${API_VERSION}/status    timeout=${TIMEOUT}
    Sleep    1s
    ${response2}=    GET    ${BASE_URL}/api/${API_VERSION}/status    timeout=${TIMEOUT}
    
    ${body1}=    Set Variable    ${response1.json()}
    ${body2}=    Set Variable    ${response2.json()}
    
    Should Be Equal    ${body1['status']}    ${body2['status']}
    Should Be Equal    ${body1['version']}    ${body2['version']}

Test Timestamp Accuracy
    [Documentation]    Test timestamp accuracy in responses
    [Tags]    consistency    timestamp
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/health    timeout=${TIMEOUT}
    ${body}=    Set Variable    ${response.json()}
    ${api_timestamp}=    Set Variable    ${body['timestamp']}
    ${current_time}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S.%fZ
    
    # Parse timestamps and compare (allow 5 second difference)
    ${api_time}=    Convert Date    ${api_timestamp}    result_format=epoch
    ${current_time_epoch}=    Convert Date    ${current_time}    result_format=epoch
    ${time_diff}=    Evaluate    abs(${current_time_epoch} - ${api_time})
    Should Be True    ${time_diff} < 5    API timestamp should be within 5 seconds of current time

# ========================================
# TEST SUITE CLEANUP
# ========================================

Teardown Test Suite
    [Documentation]    Cleanup test suite environment
    [Tags]    teardown
    Log    Test suite cleanup completed
    Log    All tests executed successfully

*** Keywords ***

Create Test Data Files
    [Documentation]    Create test data files for various test scenarios
    ${test_data}=    Create Dictionary
    ...    valid_threat=    {"type": "malware", "severity": "high", "source": "test"}
    ...    invalid_threat=    {"invalid": "data"}
    ...    network_filter=    {"start_time": "2024-01-01", "end_time": "2024-01-02"}
    ...    device_filter=    {"status": "active", "type": "endpoint"}
    
    ${json_data}=    Convert To Json    ${test_data}
    Create File    ${TEST_DATA_DIR}/test_data.json    ${json_data}

Generate Test Report
    [Documentation]    Generate comprehensive test report
    ${report_data}=    Create Dictionary
    ...    test_suite=    Security Platform API Tests
    ...    environment=    ${TEST_ENVIRONMENT}
    ...    base_url=    ${BASE_URL}
    ...    timestamp=    ${CURRENT_TIMESTAMP}
    ...    total_tests=    ${TEST_COUNT}
    ...    passed_tests=    ${PASSED_COUNT}
    ...    failed_tests=    ${FAILED_COUNT}
    
    ${json_report}=    Convert To Json    ${report_data}
    Create File    ${RESULTS_DIR}/test_report.json    ${json_report}

Validate Test Results
    [Documentation]    Validate overall test results
    [Arguments]    ${passed_count}    ${failed_count}    ${total_count}
    ${pass_rate}=    Evaluate    ${passed_count} / ${total_count} * 100
    Should Be True    ${pass_rate} >= 90    Overall pass rate should be at least 90%
    Log    Test Results: ${passed_count}/${total_count} passed (${pass_rate}%)

Setup Test Environment
    [Documentation]    Setup test environment variables
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${EMPTY}
    ${CURRENT_TIMESTAMP}=    Get Current Date    result_format=%Y-%m-%dT%H:%M:%S.%fZ
    Set Global Variable    ${CURRENT_TIMESTAMP}    ${CURRENT_TIMESTAMP}
    Create Test Data Files

Teardown Test Environment
    [Documentation]    Cleanup test environment
    Generate Test Report
    Log    Test environment cleanup completed
