/**
 * Axiom Backend V3 API Client
 */

const BASE_URL = process.env.NEXT_PUBLIC_AXIOM_BACKEND_URL || 'http://localhost:3001';

// 通用 API 請求函數
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error?.message || `HTTP ${response.status}`);
  }

  const data = await response.json();
  return data.data as T;
}

// Prometheus APIs
export const prometheusAPI = {
  healthCheck: () => apiRequest('/api/v2/prometheus/health'),
  
  query: (query: string, time?: string) =>
    apiRequest('/api/v2/prometheus/query', {
      method: 'POST',
      body: JSON.stringify({ query, time }),
    }),
  
  queryRange: (query: string, start: string, end: string, step?: string) =>
    apiRequest('/api/v2/prometheus/query-range', {
      method: 'POST',
      body: JSON.stringify({ query, start, end, step }),
    }),
  
  getRules: () => apiRequest('/api/v2/prometheus/rules'),
  
  getTargets: () => apiRequest('/api/v2/prometheus/targets'),
  
  getStatus: () => apiRequest('/api/v2/prometheus/status'),
};

// Loki APIs
export const lokiAPI = {
  healthCheck: () => apiRequest('/api/v2/loki/health'),
  
  query: (query: string, limit?: number, start?: string, end?: string) => {
    const params = new URLSearchParams({ query });
    if (limit) params.append('limit', limit.toString());
    if (start) params.append('start', start);
    if (end) params.append('end', end);
    return apiRequest(`/api/v2/loki/query?${params.toString()}`);
  },
  
  getLabels: () => apiRequest('/api/v2/loki/labels'),
  
  getLabelValues: (label: string) =>
    apiRequest(`/api/v2/loki/labels/${label}/values`),
};

// Quantum APIs
export const quantumAPI = {
  healthCheck: () => apiRequest('/api/v2/quantum/health'),
  
  generateQKD: (keyLength: number, backend?: string, shots?: number) =>
    apiRequest('/api/v2/quantum/qkd/generate', {
      method: 'POST',
      body: JSON.stringify({ key_length: keyLength, backend, shots }),
    }),
  
  classifyQSVM: (features: number[], backend?: string, featureDim?: number, shots?: number) =>
    apiRequest('/api/v2/quantum/qsvm/classify', {
      method: 'POST',
      body: JSON.stringify({ features, backend, feature_dim: featureDim, shots }),
    }),
  
  predictZeroTrust: (userId: string, ipAddress: string, deviceId?: string, features?: any, useQuantum?: boolean) =>
    apiRequest('/api/v2/quantum/zerotrust/predict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        ip_address: ipAddress,
        device_id: deviceId,
        features,
        use_quantum: useQuantum,
      }),
    }),
  
  getJob: (jobId: string) =>
    apiRequest(`/api/v2/quantum/jobs/${jobId}`),
  
  listJobs: (params?: { type?: string; status?: string; page?: number; page_size?: number }) => {
    const query = new URLSearchParams();
    if (params?.type) query.append('type', params.type);
    if (params?.status) query.append('status', params.status);
    if (params?.page) query.append('page', params.page.toString());
    if (params?.page_size) query.append('page_size', params.page_size.toString());
    return apiRequest(`/api/v2/quantum/jobs?${query.toString()}`);
  },
  
  getStats: () => apiRequest('/api/v2/quantum/stats'),
};

// Nginx APIs
export const nginxAPI = {
  getStatus: () => apiRequest('/api/v2/nginx/status'),
  
  getConfig: () => apiRequest('/api/v2/nginx/config'),
  
  updateConfig: (config: string, validate?: boolean, backup?: boolean) =>
    apiRequest('/api/v2/nginx/config', {
      method: 'PUT',
      body: JSON.stringify({ config, validate, backup }),
    }),
  
  reload: () =>
    apiRequest('/api/v2/nginx/reload', {
      method: 'POST',
    }),
};

// Windows Logs APIs
export const windowsLogsAPI = {
  batchReceive: (agentId: string, computer: string, logs: any[]) =>
    apiRequest('/api/v2/logs/windows/batch', {
      method: 'POST',
      body: JSON.stringify({ agent_id: agentId, computer, logs }),
    }),
  
  query: (params?: {
    agent_id?: string;
    log_type?: string;
    level?: string;
    start_time?: string;
    end_time?: string;
    keyword?: string;
    page?: number;
    page_size?: number;
  }) => {
    const query = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value) query.append(key, value.toString());
      });
    }
    return apiRequest(`/api/v2/logs/windows?${query.toString()}`);
  },
  
  getStats: (timeRange?: string) => {
    const query = timeRange ? `?time_range=${timeRange}` : '';
    return apiRequest(`/api/v2/logs/windows/stats${query}`);
  },
};

// 通用健康檢查
export const systemAPI = {
  healthCheck: () => apiRequest('/health'),
};

export default {
  prometheus: prometheusAPI,
  loki: lokiAPI,
  quantum: quantumAPI,
  nginx: nginxAPI,
  windowsLogs: windowsLogsAPI,
  system: systemAPI,
};


 * Axiom Backend V3 API Client
 */

const BASE_URL = process.env.NEXT_PUBLIC_AXIOM_BACKEND_URL || 'http://localhost:3001';

// 通用 API 請求函數
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error?.message || `HTTP ${response.status}`);
  }

  const data = await response.json();
  return data.data as T;
}

// Prometheus APIs
export const prometheusAPI = {
  healthCheck: () => apiRequest('/api/v2/prometheus/health'),
  
  query: (query: string, time?: string) =>
    apiRequest('/api/v2/prometheus/query', {
      method: 'POST',
      body: JSON.stringify({ query, time }),
    }),
  
  queryRange: (query: string, start: string, end: string, step?: string) =>
    apiRequest('/api/v2/prometheus/query-range', {
      method: 'POST',
      body: JSON.stringify({ query, start, end, step }),
    }),
  
  getRules: () => apiRequest('/api/v2/prometheus/rules'),
  
  getTargets: () => apiRequest('/api/v2/prometheus/targets'),
  
  getStatus: () => apiRequest('/api/v2/prometheus/status'),
};

// Loki APIs
export const lokiAPI = {
  healthCheck: () => apiRequest('/api/v2/loki/health'),
  
  query: (query: string, limit?: number, start?: string, end?: string) => {
    const params = new URLSearchParams({ query });
    if (limit) params.append('limit', limit.toString());
    if (start) params.append('start', start);
    if (end) params.append('end', end);
    return apiRequest(`/api/v2/loki/query?${params.toString()}`);
  },
  
  getLabels: () => apiRequest('/api/v2/loki/labels'),
  
  getLabelValues: (label: string) =>
    apiRequest(`/api/v2/loki/labels/${label}/values`),
};

// Quantum APIs
export const quantumAPI = {
  healthCheck: () => apiRequest('/api/v2/quantum/health'),
  
  generateQKD: (keyLength: number, backend?: string, shots?: number) =>
    apiRequest('/api/v2/quantum/qkd/generate', {
      method: 'POST',
      body: JSON.stringify({ key_length: keyLength, backend, shots }),
    }),
  
  classifyQSVM: (features: number[], backend?: string, featureDim?: number, shots?: number) =>
    apiRequest('/api/v2/quantum/qsvm/classify', {
      method: 'POST',
      body: JSON.stringify({ features, backend, feature_dim: featureDim, shots }),
    }),
  
  predictZeroTrust: (userId: string, ipAddress: string, deviceId?: string, features?: any, useQuantum?: boolean) =>
    apiRequest('/api/v2/quantum/zerotrust/predict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        ip_address: ipAddress,
        device_id: deviceId,
        features,
        use_quantum: useQuantum,
      }),
    }),
  
  getJob: (jobId: string) =>
    apiRequest(`/api/v2/quantum/jobs/${jobId}`),
  
  listJobs: (params?: { type?: string; status?: string; page?: number; page_size?: number }) => {
    const query = new URLSearchParams();
    if (params?.type) query.append('type', params.type);
    if (params?.status) query.append('status', params.status);
    if (params?.page) query.append('page', params.page.toString());
    if (params?.page_size) query.append('page_size', params.page_size.toString());
    return apiRequest(`/api/v2/quantum/jobs?${query.toString()}`);
  },
  
  getStats: () => apiRequest('/api/v2/quantum/stats'),
};

// Nginx APIs
export const nginxAPI = {
  getStatus: () => apiRequest('/api/v2/nginx/status'),
  
  getConfig: () => apiRequest('/api/v2/nginx/config'),
  
  updateConfig: (config: string, validate?: boolean, backup?: boolean) =>
    apiRequest('/api/v2/nginx/config', {
      method: 'PUT',
      body: JSON.stringify({ config, validate, backup }),
    }),
  
  reload: () =>
    apiRequest('/api/v2/nginx/reload', {
      method: 'POST',
    }),
};

// Windows Logs APIs
export const windowsLogsAPI = {
  batchReceive: (agentId: string, computer: string, logs: any[]) =>
    apiRequest('/api/v2/logs/windows/batch', {
      method: 'POST',
      body: JSON.stringify({ agent_id: agentId, computer, logs }),
    }),
  
  query: (params?: {
    agent_id?: string;
    log_type?: string;
    level?: string;
    start_time?: string;
    end_time?: string;
    keyword?: string;
    page?: number;
    page_size?: number;
  }) => {
    const query = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value) query.append(key, value.toString());
      });
    }
    return apiRequest(`/api/v2/logs/windows?${query.toString()}`);
  },
  
  getStats: (timeRange?: string) => {
    const query = timeRange ? `?time_range=${timeRange}` : '';
    return apiRequest(`/api/v2/logs/windows/stats${query}`);
  },
};

// 通用健康檢查
export const systemAPI = {
  healthCheck: () => apiRequest('/health'),
};

export default {
  prometheus: prometheusAPI,
  loki: lokiAPI,
  quantum: quantumAPI,
  nginx: nginxAPI,
  windowsLogs: windowsLogsAPI,
  system: systemAPI,
};

 * Axiom Backend V3 API Client
 */

const BASE_URL = process.env.NEXT_PUBLIC_AXIOM_BACKEND_URL || 'http://localhost:3001';

// 通用 API 請求函數
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error?.message || `HTTP ${response.status}`);
  }

  const data = await response.json();
  return data.data as T;
}

// Prometheus APIs
export const prometheusAPI = {
  healthCheck: () => apiRequest('/api/v2/prometheus/health'),
  
  query: (query: string, time?: string) =>
    apiRequest('/api/v2/prometheus/query', {
      method: 'POST',
      body: JSON.stringify({ query, time }),
    }),
  
  queryRange: (query: string, start: string, end: string, step?: string) =>
    apiRequest('/api/v2/prometheus/query-range', {
      method: 'POST',
      body: JSON.stringify({ query, start, end, step }),
    }),
  
  getRules: () => apiRequest('/api/v2/prometheus/rules'),
  
  getTargets: () => apiRequest('/api/v2/prometheus/targets'),
  
  getStatus: () => apiRequest('/api/v2/prometheus/status'),
};

// Loki APIs
export const lokiAPI = {
  healthCheck: () => apiRequest('/api/v2/loki/health'),
  
  query: (query: string, limit?: number, start?: string, end?: string) => {
    const params = new URLSearchParams({ query });
    if (limit) params.append('limit', limit.toString());
    if (start) params.append('start', start);
    if (end) params.append('end', end);
    return apiRequest(`/api/v2/loki/query?${params.toString()}`);
  },
  
  getLabels: () => apiRequest('/api/v2/loki/labels'),
  
  getLabelValues: (label: string) =>
    apiRequest(`/api/v2/loki/labels/${label}/values`),
};

// Quantum APIs
export const quantumAPI = {
  healthCheck: () => apiRequest('/api/v2/quantum/health'),
  
  generateQKD: (keyLength: number, backend?: string, shots?: number) =>
    apiRequest('/api/v2/quantum/qkd/generate', {
      method: 'POST',
      body: JSON.stringify({ key_length: keyLength, backend, shots }),
    }),
  
  classifyQSVM: (features: number[], backend?: string, featureDim?: number, shots?: number) =>
    apiRequest('/api/v2/quantum/qsvm/classify', {
      method: 'POST',
      body: JSON.stringify({ features, backend, feature_dim: featureDim, shots }),
    }),
  
  predictZeroTrust: (userId: string, ipAddress: string, deviceId?: string, features?: any, useQuantum?: boolean) =>
    apiRequest('/api/v2/quantum/zerotrust/predict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        ip_address: ipAddress,
        device_id: deviceId,
        features,
        use_quantum: useQuantum,
      }),
    }),
  
  getJob: (jobId: string) =>
    apiRequest(`/api/v2/quantum/jobs/${jobId}`),
  
  listJobs: (params?: { type?: string; status?: string; page?: number; page_size?: number }) => {
    const query = new URLSearchParams();
    if (params?.type) query.append('type', params.type);
    if (params?.status) query.append('status', params.status);
    if (params?.page) query.append('page', params.page.toString());
    if (params?.page_size) query.append('page_size', params.page_size.toString());
    return apiRequest(`/api/v2/quantum/jobs?${query.toString()}`);
  },
  
  getStats: () => apiRequest('/api/v2/quantum/stats'),
};

// Nginx APIs
export const nginxAPI = {
  getStatus: () => apiRequest('/api/v2/nginx/status'),
  
  getConfig: () => apiRequest('/api/v2/nginx/config'),
  
  updateConfig: (config: string, validate?: boolean, backup?: boolean) =>
    apiRequest('/api/v2/nginx/config', {
      method: 'PUT',
      body: JSON.stringify({ config, validate, backup }),
    }),
  
  reload: () =>
    apiRequest('/api/v2/nginx/reload', {
      method: 'POST',
    }),
};

// Windows Logs APIs
export const windowsLogsAPI = {
  batchReceive: (agentId: string, computer: string, logs: any[]) =>
    apiRequest('/api/v2/logs/windows/batch', {
      method: 'POST',
      body: JSON.stringify({ agent_id: agentId, computer, logs }),
    }),
  
  query: (params?: {
    agent_id?: string;
    log_type?: string;
    level?: string;
    start_time?: string;
    end_time?: string;
    keyword?: string;
    page?: number;
    page_size?: number;
  }) => {
    const query = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value) query.append(key, value.toString());
      });
    }
    return apiRequest(`/api/v2/logs/windows?${query.toString()}`);
  },
  
  getStats: (timeRange?: string) => {
    const query = timeRange ? `?time_range=${timeRange}` : '';
    return apiRequest(`/api/v2/logs/windows/stats${query}`);
  },
};

// 通用健康檢查
export const systemAPI = {
  healthCheck: () => apiRequest('/health'),
};

export default {
  prometheus: prometheusAPI,
  loki: lokiAPI,
  quantum: quantumAPI,
  nginx: nginxAPI,
  windowsLogs: windowsLogsAPI,
  system: systemAPI,
};


 * Axiom Backend V3 API Client
 */

const BASE_URL = process.env.NEXT_PUBLIC_AXIOM_BACKEND_URL || 'http://localhost:3001';

// 通用 API 請求函數
async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${BASE_URL}${endpoint}`;
  
  const response = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({}));
    throw new Error(error.error?.message || `HTTP ${response.status}`);
  }

  const data = await response.json();
  return data.data as T;
}

// Prometheus APIs
export const prometheusAPI = {
  healthCheck: () => apiRequest('/api/v2/prometheus/health'),
  
  query: (query: string, time?: string) =>
    apiRequest('/api/v2/prometheus/query', {
      method: 'POST',
      body: JSON.stringify({ query, time }),
    }),
  
  queryRange: (query: string, start: string, end: string, step?: string) =>
    apiRequest('/api/v2/prometheus/query-range', {
      method: 'POST',
      body: JSON.stringify({ query, start, end, step }),
    }),
  
  getRules: () => apiRequest('/api/v2/prometheus/rules'),
  
  getTargets: () => apiRequest('/api/v2/prometheus/targets'),
  
  getStatus: () => apiRequest('/api/v2/prometheus/status'),
};

// Loki APIs
export const lokiAPI = {
  healthCheck: () => apiRequest('/api/v2/loki/health'),
  
  query: (query: string, limit?: number, start?: string, end?: string) => {
    const params = new URLSearchParams({ query });
    if (limit) params.append('limit', limit.toString());
    if (start) params.append('start', start);
    if (end) params.append('end', end);
    return apiRequest(`/api/v2/loki/query?${params.toString()}`);
  },
  
  getLabels: () => apiRequest('/api/v2/loki/labels'),
  
  getLabelValues: (label: string) =>
    apiRequest(`/api/v2/loki/labels/${label}/values`),
};

// Quantum APIs
export const quantumAPI = {
  healthCheck: () => apiRequest('/api/v2/quantum/health'),
  
  generateQKD: (keyLength: number, backend?: string, shots?: number) =>
    apiRequest('/api/v2/quantum/qkd/generate', {
      method: 'POST',
      body: JSON.stringify({ key_length: keyLength, backend, shots }),
    }),
  
  classifyQSVM: (features: number[], backend?: string, featureDim?: number, shots?: number) =>
    apiRequest('/api/v2/quantum/qsvm/classify', {
      method: 'POST',
      body: JSON.stringify({ features, backend, feature_dim: featureDim, shots }),
    }),
  
  predictZeroTrust: (userId: string, ipAddress: string, deviceId?: string, features?: any, useQuantum?: boolean) =>
    apiRequest('/api/v2/quantum/zerotrust/predict', {
      method: 'POST',
      body: JSON.stringify({
        user_id: userId,
        ip_address: ipAddress,
        device_id: deviceId,
        features,
        use_quantum: useQuantum,
      }),
    }),
  
  getJob: (jobId: string) =>
    apiRequest(`/api/v2/quantum/jobs/${jobId}`),
  
  listJobs: (params?: { type?: string; status?: string; page?: number; page_size?: number }) => {
    const query = new URLSearchParams();
    if (params?.type) query.append('type', params.type);
    if (params?.status) query.append('status', params.status);
    if (params?.page) query.append('page', params.page.toString());
    if (params?.page_size) query.append('page_size', params.page_size.toString());
    return apiRequest(`/api/v2/quantum/jobs?${query.toString()}`);
  },
  
  getStats: () => apiRequest('/api/v2/quantum/stats'),
};

// Nginx APIs
export const nginxAPI = {
  getStatus: () => apiRequest('/api/v2/nginx/status'),
  
  getConfig: () => apiRequest('/api/v2/nginx/config'),
  
  updateConfig: (config: string, validate?: boolean, backup?: boolean) =>
    apiRequest('/api/v2/nginx/config', {
      method: 'PUT',
      body: JSON.stringify({ config, validate, backup }),
    }),
  
  reload: () =>
    apiRequest('/api/v2/nginx/reload', {
      method: 'POST',
    }),
};

// Windows Logs APIs
export const windowsLogsAPI = {
  batchReceive: (agentId: string, computer: string, logs: any[]) =>
    apiRequest('/api/v2/logs/windows/batch', {
      method: 'POST',
      body: JSON.stringify({ agent_id: agentId, computer, logs }),
    }),
  
  query: (params?: {
    agent_id?: string;
    log_type?: string;
    level?: string;
    start_time?: string;
    end_time?: string;
    keyword?: string;
    page?: number;
    page_size?: number;
  }) => {
    const query = new URLSearchParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value) query.append(key, value.toString());
      });
    }
    return apiRequest(`/api/v2/logs/windows?${query.toString()}`);
  },
  
  getStats: (timeRange?: string) => {
    const query = timeRange ? `?time_range=${timeRange}` : '';
    return apiRequest(`/api/v2/logs/windows/stats${query}`);
  },
};

// 通用健康檢查
export const systemAPI = {
  healthCheck: () => apiRequest('/health'),
};

export default {
  prometheus: prometheusAPI,
  loki: lokiAPI,
  quantum: quantumAPI,
  nginx: nginxAPI,
  windowsLogs: windowsLogsAPI,
  system: systemAPI,
};

