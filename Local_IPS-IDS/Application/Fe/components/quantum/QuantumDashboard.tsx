import React, { useState, useEffect } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { quantumAPI } from '../../services/axiom-api';

interface QuantumJob {
  job_id: string;
  type: string;
  status: string;
  backend?: string;
  submitted_at: string;
  completed_at?: string;
}

interface QuantumStats {
  total_jobs: number;
  completed_jobs: number;
  failed_jobs: number;
  running_jobs: number;
  success_rate: number;
  jobs_by_type: Record<string, number>;
}

const QuantumDashboard: React.FC = () => {
  const [stats, setStats] = useState<QuantumStats | null>(null);
  const [jobs, setJobs] = useState<QuantumJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // QKD 表單狀態
  const [qkdKeyLength, setQkdKeyLength] = useState(256);
  const [qkdBackend, setQkdBackend] = useState('simulator');
  const [qkdSubmitting, setQkdSubmitting] = useState(false);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 10000); // 每 10 秒刷新
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsData, jobsData] = await Promise.all([
        quantumAPI.getStats(),
        quantumAPI.listJobs({ page: 1, page_size: 10 }),
      ]);
      
      setStats(statsData as QuantumStats);
      setJobs((jobsData as any).jobs || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateQKD = async () => {
    setQkdSubmitting(true);
    try {
      const result = await quantumAPI.generateQKD(qkdKeyLength, qkdBackend, 1024);
      alert(`QKD 作業已提交！作業 ID: ${(result as any).job_id}`);
      loadData(); // 刷新數據
    } catch (err) {
      alert('QKD 生成失敗: ' + (err instanceof Error ? err.message : 'Unknown error'));
    } finally {
      setQkdSubmitting(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'running':
        return 'bg-blue-500';
      case 'pending':
        return 'bg-yellow-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <h1 className="text-3xl font-bold">量子功能控制中心</h1>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* 統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總作業數</h3>
          <p className="text-3xl font-bold mt-2">{stats?.total_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">已完成</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{stats?.completed_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">運行中</h3>
          <p className="text-3xl font-bold mt-2 text-blue-600">{stats?.running_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">成功率</h3>
          <p className="text-3xl font-bold mt-2">
            {((stats?.success_rate || 0) * 100).toFixed(1)}%
          </p>
        </Card>
      </div>

      {/* QKD 生成表單 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">生成量子密鑰 (QKD)</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2">密鑰長度 (bits)</label>
            <input
              type="number"
              value={qkdKeyLength}
              onChange={(e) => setQkdKeyLength(Number(e.target.value))}
              min={64}
              max={2048}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium mb-2">後端</label>
            <select
              value={qkdBackend}
              onChange={(e) => setQkdBackend(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="simulator">模擬器</option>
              <option value="ibm_quantum">IBM Quantum</option>
            </select>
          </div>
          
          <div className="flex items-end">
            <Button
              onClick={handleGenerateQKD}
              disabled={qkdSubmitting}
              className="w-full"
            >
              {qkdSubmitting ? '提交中...' : '生成密鑰'}
            </Button>
          </div>
        </div>
      </Card>

      {/* 最近作業列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">最近量子作業</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">作業 ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">類型</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">後端</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">提交時間</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {jobs.map((job) => (
                <tr key={job.job_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-mono">{job.job_id}</td>
                  <td className="px-6 py-4 text-sm">
                    <Badge>{job.type}</Badge>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white ${getStatusColor(job.status)}`}>
                      {job.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">{job.backend || 'N/A'}</td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(job.submitted_at).toLocaleString('zh-TW')}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      {/* 作業類型分布 */}
      {stats && stats.jobs_by_type && (
        <Card className="p-6">
          <h2 className="text-xl font-bold mb-4">作業類型分布</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {Object.entries(stats.jobs_by_type).map(([type, count]) => (
              <div key={type} className="text-center">
                <p className="text-2xl font-bold">{count}</p>
                <p className="text-sm text-gray-600">{type.toUpperCase()}</p>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};

export default QuantumDashboard;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { quantumAPI } from '../../services/axiom-api';

interface QuantumJob {
  job_id: string;
  type: string;
  status: string;
  backend?: string;
  submitted_at: string;
  completed_at?: string;
}

interface QuantumStats {
  total_jobs: number;
  completed_jobs: number;
  failed_jobs: number;
  running_jobs: number;
  success_rate: number;
  jobs_by_type: Record<string, number>;
}

const QuantumDashboard: React.FC = () => {
  const [stats, setStats] = useState<QuantumStats | null>(null);
  const [jobs, setJobs] = useState<QuantumJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // QKD 表單狀態
  const [qkdKeyLength, setQkdKeyLength] = useState(256);
  const [qkdBackend, setQkdBackend] = useState('simulator');
  const [qkdSubmitting, setQkdSubmitting] = useState(false);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 10000); // 每 10 秒刷新
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsData, jobsData] = await Promise.all([
        quantumAPI.getStats(),
        quantumAPI.listJobs({ page: 1, page_size: 10 }),
      ]);
      
      setStats(statsData as QuantumStats);
      setJobs((jobsData as any).jobs || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateQKD = async () => {
    setQkdSubmitting(true);
    try {
      const result = await quantumAPI.generateQKD(qkdKeyLength, qkdBackend, 1024);
      alert(`QKD 作業已提交！作業 ID: ${(result as any).job_id}`);
      loadData(); // 刷新數據
    } catch (err) {
      alert('QKD 生成失敗: ' + (err instanceof Error ? err.message : 'Unknown error'));
    } finally {
      setQkdSubmitting(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'running':
        return 'bg-blue-500';
      case 'pending':
        return 'bg-yellow-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <h1 className="text-3xl font-bold">量子功能控制中心</h1>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* 統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總作業數</h3>
          <p className="text-3xl font-bold mt-2">{stats?.total_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">已完成</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{stats?.completed_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">運行中</h3>
          <p className="text-3xl font-bold mt-2 text-blue-600">{stats?.running_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">成功率</h3>
          <p className="text-3xl font-bold mt-2">
            {((stats?.success_rate || 0) * 100).toFixed(1)}%
          </p>
        </Card>
      </div>

      {/* QKD 生成表單 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">生成量子密鑰 (QKD)</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2">密鑰長度 (bits)</label>
            <input
              type="number"
              value={qkdKeyLength}
              onChange={(e) => setQkdKeyLength(Number(e.target.value))}
              min={64}
              max={2048}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium mb-2">後端</label>
            <select
              value={qkdBackend}
              onChange={(e) => setQkdBackend(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="simulator">模擬器</option>
              <option value="ibm_quantum">IBM Quantum</option>
            </select>
          </div>
          
          <div className="flex items-end">
            <Button
              onClick={handleGenerateQKD}
              disabled={qkdSubmitting}
              className="w-full"
            >
              {qkdSubmitting ? '提交中...' : '生成密鑰'}
            </Button>
          </div>
        </div>
      </Card>

      {/* 最近作業列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">最近量子作業</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">作業 ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">類型</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">後端</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">提交時間</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {jobs.map((job) => (
                <tr key={job.job_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-mono">{job.job_id}</td>
                  <td className="px-6 py-4 text-sm">
                    <Badge>{job.type}</Badge>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white ${getStatusColor(job.status)}`}>
                      {job.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">{job.backend || 'N/A'}</td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(job.submitted_at).toLocaleString('zh-TW')}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      {/* 作業類型分布 */}
      {stats && stats.jobs_by_type && (
        <Card className="p-6">
          <h2 className="text-xl font-bold mb-4">作業類型分布</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {Object.entries(stats.jobs_by_type).map(([type, count]) => (
              <div key={type} className="text-center">
                <p className="text-2xl font-bold">{count}</p>
                <p className="text-sm text-gray-600">{type.toUpperCase()}</p>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};

export default QuantumDashboard;

import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { quantumAPI } from '../../services/axiom-api';

interface QuantumJob {
  job_id: string;
  type: string;
  status: string;
  backend?: string;
  submitted_at: string;
  completed_at?: string;
}

interface QuantumStats {
  total_jobs: number;
  completed_jobs: number;
  failed_jobs: number;
  running_jobs: number;
  success_rate: number;
  jobs_by_type: Record<string, number>;
}

const QuantumDashboard: React.FC = () => {
  const [stats, setStats] = useState<QuantumStats | null>(null);
  const [jobs, setJobs] = useState<QuantumJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // QKD 表單狀態
  const [qkdKeyLength, setQkdKeyLength] = useState(256);
  const [qkdBackend, setQkdBackend] = useState('simulator');
  const [qkdSubmitting, setQkdSubmitting] = useState(false);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 10000); // 每 10 秒刷新
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsData, jobsData] = await Promise.all([
        quantumAPI.getStats(),
        quantumAPI.listJobs({ page: 1, page_size: 10 }),
      ]);
      
      setStats(statsData as QuantumStats);
      setJobs((jobsData as any).jobs || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateQKD = async () => {
    setQkdSubmitting(true);
    try {
      const result = await quantumAPI.generateQKD(qkdKeyLength, qkdBackend, 1024);
      alert(`QKD 作業已提交！作業 ID: ${(result as any).job_id}`);
      loadData(); // 刷新數據
    } catch (err) {
      alert('QKD 生成失敗: ' + (err instanceof Error ? err.message : 'Unknown error'));
    } finally {
      setQkdSubmitting(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'running':
        return 'bg-blue-500';
      case 'pending':
        return 'bg-yellow-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <h1 className="text-3xl font-bold">量子功能控制中心</h1>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* 統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總作業數</h3>
          <p className="text-3xl font-bold mt-2">{stats?.total_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">已完成</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{stats?.completed_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">運行中</h3>
          <p className="text-3xl font-bold mt-2 text-blue-600">{stats?.running_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">成功率</h3>
          <p className="text-3xl font-bold mt-2">
            {((stats?.success_rate || 0) * 100).toFixed(1)}%
          </p>
        </Card>
      </div>

      {/* QKD 生成表單 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">生成量子密鑰 (QKD)</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2">密鑰長度 (bits)</label>
            <input
              type="number"
              value={qkdKeyLength}
              onChange={(e) => setQkdKeyLength(Number(e.target.value))}
              min={64}
              max={2048}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium mb-2">後端</label>
            <select
              value={qkdBackend}
              onChange={(e) => setQkdBackend(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="simulator">模擬器</option>
              <option value="ibm_quantum">IBM Quantum</option>
            </select>
          </div>
          
          <div className="flex items-end">
            <Button
              onClick={handleGenerateQKD}
              disabled={qkdSubmitting}
              className="w-full"
            >
              {qkdSubmitting ? '提交中...' : '生成密鑰'}
            </Button>
          </div>
        </div>
      </Card>

      {/* 最近作業列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">最近量子作業</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">作業 ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">類型</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">後端</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">提交時間</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {jobs.map((job) => (
                <tr key={job.job_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-mono">{job.job_id}</td>
                  <td className="px-6 py-4 text-sm">
                    <Badge>{job.type}</Badge>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white ${getStatusColor(job.status)}`}>
                      {job.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">{job.backend || 'N/A'}</td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(job.submitted_at).toLocaleString('zh-TW')}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      {/* 作業類型分布 */}
      {stats && stats.jobs_by_type && (
        <Card className="p-6">
          <h2 className="text-xl font-bold mb-4">作業類型分布</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {Object.entries(stats.jobs_by_type).map(([type, count]) => (
              <div key={type} className="text-center">
                <p className="text-2xl font-bold">{count}</p>
                <p className="text-sm text-gray-600">{type.toUpperCase()}</p>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};

export default QuantumDashboard;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { quantumAPI } from '../../services/axiom-api';

interface QuantumJob {
  job_id: string;
  type: string;
  status: string;
  backend?: string;
  submitted_at: string;
  completed_at?: string;
}

interface QuantumStats {
  total_jobs: number;
  completed_jobs: number;
  failed_jobs: number;
  running_jobs: number;
  success_rate: number;
  jobs_by_type: Record<string, number>;
}

const QuantumDashboard: React.FC = () => {
  const [stats, setStats] = useState<QuantumStats | null>(null);
  const [jobs, setJobs] = useState<QuantumJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // QKD 表單狀態
  const [qkdKeyLength, setQkdKeyLength] = useState(256);
  const [qkdBackend, setQkdBackend] = useState('simulator');
  const [qkdSubmitting, setQkdSubmitting] = useState(false);

  useEffect(() => {
    loadData();
    const interval = setInterval(loadData, 10000); // 每 10 秒刷新
    return () => clearInterval(interval);
  }, []);

  const loadData = async () => {
    try {
      const [statsData, jobsData] = await Promise.all([
        quantumAPI.getStats(),
        quantumAPI.listJobs({ page: 1, page_size: 10 }),
      ]);
      
      setStats(statsData as QuantumStats);
      setJobs((jobsData as any).jobs || []);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load data');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateQKD = async () => {
    setQkdSubmitting(true);
    try {
      const result = await quantumAPI.generateQKD(qkdKeyLength, qkdBackend, 1024);
      alert(`QKD 作業已提交！作業 ID: ${(result as any).job_id}`);
      loadData(); // 刷新數據
    } catch (err) {
      alert('QKD 生成失敗: ' + (err instanceof Error ? err.message : 'Unknown error'));
    } finally {
      setQkdSubmitting(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'bg-green-500';
      case 'running':
        return 'bg-blue-500';
      case 'pending':
        return 'bg-yellow-500';
      case 'failed':
        return 'bg-red-500';
      default:
        return 'bg-gray-500';
    }
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <h1 className="text-3xl font-bold">量子功能控制中心</h1>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      )}

      {/* 統計卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總作業數</h3>
          <p className="text-3xl font-bold mt-2">{stats?.total_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">已完成</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{stats?.completed_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">運行中</h3>
          <p className="text-3xl font-bold mt-2 text-blue-600">{stats?.running_jobs || 0}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">成功率</h3>
          <p className="text-3xl font-bold mt-2">
            {((stats?.success_rate || 0) * 100).toFixed(1)}%
          </p>
        </Card>
      </div>

      {/* QKD 生成表單 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">生成量子密鑰 (QKD)</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2">密鑰長度 (bits)</label>
            <input
              type="number"
              value={qkdKeyLength}
              onChange={(e) => setQkdKeyLength(Number(e.target.value))}
              min={64}
              max={2048}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium mb-2">後端</label>
            <select
              value={qkdBackend}
              onChange={(e) => setQkdBackend(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="simulator">模擬器</option>
              <option value="ibm_quantum">IBM Quantum</option>
            </select>
          </div>
          
          <div className="flex items-end">
            <Button
              onClick={handleGenerateQKD}
              disabled={qkdSubmitting}
              className="w-full"
            >
              {qkdSubmitting ? '提交中...' : '生成密鑰'}
            </Button>
          </div>
        </div>
      </Card>

      {/* 最近作業列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">最近量子作業</h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">作業 ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">類型</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">狀態</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">後端</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">提交時間</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {jobs.map((job) => (
                <tr key={job.job_id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-mono">{job.job_id}</td>
                  <td className="px-6 py-4 text-sm">
                    <Badge>{job.type}</Badge>
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white ${getStatusColor(job.status)}`}>
                      {job.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm">{job.backend || 'N/A'}</td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(job.submitted_at).toLocaleString('zh-TW')}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>

      {/* 作業類型分布 */}
      {stats && stats.jobs_by_type && (
        <Card className="p-6">
          <h2 className="text-xl font-bold mb-4">作業類型分布</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {Object.entries(stats.jobs_by_type).map(([type, count]) => (
              <div key={type} className="text-center">
                <p className="text-2xl font-bold">{count}</p>
                <p className="text-sm text-gray-600">{type.toUpperCase()}</p>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
};

export default QuantumDashboard;

