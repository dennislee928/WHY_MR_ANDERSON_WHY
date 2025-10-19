import React, { useState, useEffect } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { prometheusAPI, lokiAPI, quantumAPI, nginxAPI } from '../../services/axiom-api';

interface ServiceStatus {
  name: string;
  status: string;
  healthy: boolean;
  message?: string;
}

const ServicesManagement: React.FC = () => {
  const [services, setServices] = useState<ServiceStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const serviceChecks = [
    { name: 'Prometheus', check: prometheusAPI.healthCheck },
    { name: 'Loki', check: lokiAPI.healthCheck },
    { name: 'Quantum', check: quantumAPI.healthCheck },
    { name: 'Nginx', check: async () => nginxAPI.getStatus() },
  ];

  useEffect(() => {
    checkAllServices();
    const interval = setInterval(checkAllServices, 30000); // 每 30 秒刷新
    return () => clearInterval(interval);
  }, []);

  const checkAllServices = async () => {
    const results: ServiceStatus[] = [];

    for (const service of serviceChecks) {
      try {
        await service.check();
        results.push({
          name: service.name,
          status: 'healthy',
          healthy: true,
          message: '服務正常運行',
        });
      } catch (error) {
        results.push({
          name: service.name,
          status: 'unhealthy',
          healthy: false,
          message: error instanceof Error ? error.message : '服務不可用',
        });
      }
    }

    setServices(results);
    setLoading(false);
  };

  const getStatusColor = (healthy: boolean) => {
    return healthy ? 'bg-green-500' : 'bg-red-500';
  };

  const getStatusIcon = (healthy: boolean) => {
    return healthy ? '✓' : '✗';
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  const healthyCount = services.filter(s => s.healthy).length;
  const totalCount = services.length;

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">服務管理</h1>
        <Button onClick={checkAllServices}>刷新狀態</Button>
      </div>

      {/* 總覽卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總服務數</h3>
          <p className="text-3xl font-bold mt-2">{totalCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康服務</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{healthyCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康率</h3>
          <p className="text-3xl font-bold mt-2">
            {((healthyCount / totalCount) * 100).toFixed(0)}%
          </p>
        </Card>
      </div>

      {/* 服務列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">服務狀態</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {services.map((service) => (
            <div
              key={service.name}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-center space-x-4">
                <div className={`w-12 h-12 rounded-full ${getStatusColor(service.healthy)} flex items-center justify-center text-white text-xl font-bold`}>
                  {getStatusIcon(service.healthy)}
                </div>
                <div>
                  <h3 className="font-semibold text-lg">{service.name}</h3>
                  <p className="text-sm text-gray-600">{service.message}</p>
                </div>
              </div>
              <Badge variant={service.healthy ? 'default' : 'destructive'}>
                {service.status}
              </Badge>
            </div>
          ))}
        </div>
      </Card>

      {/* 快速操作 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">快速操作</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Button variant="outline" className="h-20">
            查看 Prometheus
          </Button>
          <Button variant="outline" className="h-20">
            查看 Grafana
          </Button>
          <Button variant="outline" className="h-20">
            查看 Loki 日誌
          </Button>
          <Button variant="outline" className="h-20">
            量子作業管理
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default ServicesManagement;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { prometheusAPI, lokiAPI, quantumAPI, nginxAPI } from '../../services/axiom-api';

interface ServiceStatus {
  name: string;
  status: string;
  healthy: boolean;
  message?: string;
}

const ServicesManagement: React.FC = () => {
  const [services, setServices] = useState<ServiceStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const serviceChecks = [
    { name: 'Prometheus', check: prometheusAPI.healthCheck },
    { name: 'Loki', check: lokiAPI.healthCheck },
    { name: 'Quantum', check: quantumAPI.healthCheck },
    { name: 'Nginx', check: async () => nginxAPI.getStatus() },
  ];

  useEffect(() => {
    checkAllServices();
    const interval = setInterval(checkAllServices, 30000); // 每 30 秒刷新
    return () => clearInterval(interval);
  }, []);

  const checkAllServices = async () => {
    const results: ServiceStatus[] = [];

    for (const service of serviceChecks) {
      try {
        await service.check();
        results.push({
          name: service.name,
          status: 'healthy',
          healthy: true,
          message: '服務正常運行',
        });
      } catch (error) {
        results.push({
          name: service.name,
          status: 'unhealthy',
          healthy: false,
          message: error instanceof Error ? error.message : '服務不可用',
        });
      }
    }

    setServices(results);
    setLoading(false);
  };

  const getStatusColor = (healthy: boolean) => {
    return healthy ? 'bg-green-500' : 'bg-red-500';
  };

  const getStatusIcon = (healthy: boolean) => {
    return healthy ? '✓' : '✗';
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  const healthyCount = services.filter(s => s.healthy).length;
  const totalCount = services.length;

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">服務管理</h1>
        <Button onClick={checkAllServices}>刷新狀態</Button>
      </div>

      {/* 總覽卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總服務數</h3>
          <p className="text-3xl font-bold mt-2">{totalCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康服務</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{healthyCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康率</h3>
          <p className="text-3xl font-bold mt-2">
            {((healthyCount / totalCount) * 100).toFixed(0)}%
          </p>
        </Card>
      </div>

      {/* 服務列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">服務狀態</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {services.map((service) => (
            <div
              key={service.name}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-center space-x-4">
                <div className={`w-12 h-12 rounded-full ${getStatusColor(service.healthy)} flex items-center justify-center text-white text-xl font-bold`}>
                  {getStatusIcon(service.healthy)}
                </div>
                <div>
                  <h3 className="font-semibold text-lg">{service.name}</h3>
                  <p className="text-sm text-gray-600">{service.message}</p>
                </div>
              </div>
              <Badge variant={service.healthy ? 'default' : 'destructive'}>
                {service.status}
              </Badge>
            </div>
          ))}
        </div>
      </Card>

      {/* 快速操作 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">快速操作</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Button variant="outline" className="h-20">
            查看 Prometheus
          </Button>
          <Button variant="outline" className="h-20">
            查看 Grafana
          </Button>
          <Button variant="outline" className="h-20">
            查看 Loki 日誌
          </Button>
          <Button variant="outline" className="h-20">
            量子作業管理
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default ServicesManagement;

import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { prometheusAPI, lokiAPI, quantumAPI, nginxAPI } from '../../services/axiom-api';

interface ServiceStatus {
  name: string;
  status: string;
  healthy: boolean;
  message?: string;
}

const ServicesManagement: React.FC = () => {
  const [services, setServices] = useState<ServiceStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const serviceChecks = [
    { name: 'Prometheus', check: prometheusAPI.healthCheck },
    { name: 'Loki', check: lokiAPI.healthCheck },
    { name: 'Quantum', check: quantumAPI.healthCheck },
    { name: 'Nginx', check: async () => nginxAPI.getStatus() },
  ];

  useEffect(() => {
    checkAllServices();
    const interval = setInterval(checkAllServices, 30000); // 每 30 秒刷新
    return () => clearInterval(interval);
  }, []);

  const checkAllServices = async () => {
    const results: ServiceStatus[] = [];

    for (const service of serviceChecks) {
      try {
        await service.check();
        results.push({
          name: service.name,
          status: 'healthy',
          healthy: true,
          message: '服務正常運行',
        });
      } catch (error) {
        results.push({
          name: service.name,
          status: 'unhealthy',
          healthy: false,
          message: error instanceof Error ? error.message : '服務不可用',
        });
      }
    }

    setServices(results);
    setLoading(false);
  };

  const getStatusColor = (healthy: boolean) => {
    return healthy ? 'bg-green-500' : 'bg-red-500';
  };

  const getStatusIcon = (healthy: boolean) => {
    return healthy ? '✓' : '✗';
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  const healthyCount = services.filter(s => s.healthy).length;
  const totalCount = services.length;

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">服務管理</h1>
        <Button onClick={checkAllServices}>刷新狀態</Button>
      </div>

      {/* 總覽卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總服務數</h3>
          <p className="text-3xl font-bold mt-2">{totalCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康服務</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{healthyCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康率</h3>
          <p className="text-3xl font-bold mt-2">
            {((healthyCount / totalCount) * 100).toFixed(0)}%
          </p>
        </Card>
      </div>

      {/* 服務列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">服務狀態</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {services.map((service) => (
            <div
              key={service.name}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-center space-x-4">
                <div className={`w-12 h-12 rounded-full ${getStatusColor(service.healthy)} flex items-center justify-center text-white text-xl font-bold`}>
                  {getStatusIcon(service.healthy)}
                </div>
                <div>
                  <h3 className="font-semibold text-lg">{service.name}</h3>
                  <p className="text-sm text-gray-600">{service.message}</p>
                </div>
              </div>
              <Badge variant={service.healthy ? 'default' : 'destructive'}>
                {service.status}
              </Badge>
            </div>
          ))}
        </div>
      </Card>

      {/* 快速操作 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">快速操作</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Button variant="outline" className="h-20">
            查看 Prometheus
          </Button>
          <Button variant="outline" className="h-20">
            查看 Grafana
          </Button>
          <Button variant="outline" className="h-20">
            查看 Loki 日誌
          </Button>
          <Button variant="outline" className="h-20">
            量子作業管理
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default ServicesManagement;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { prometheusAPI, lokiAPI, quantumAPI, nginxAPI } from '../../services/axiom-api';

interface ServiceStatus {
  name: string;
  status: string;
  healthy: boolean;
  message?: string;
}

const ServicesManagement: React.FC = () => {
  const [services, setServices] = useState<ServiceStatus[]>([]);
  const [loading, setLoading] = useState(true);

  const serviceChecks = [
    { name: 'Prometheus', check: prometheusAPI.healthCheck },
    { name: 'Loki', check: lokiAPI.healthCheck },
    { name: 'Quantum', check: quantumAPI.healthCheck },
    { name: 'Nginx', check: async () => nginxAPI.getStatus() },
  ];

  useEffect(() => {
    checkAllServices();
    const interval = setInterval(checkAllServices, 30000); // 每 30 秒刷新
    return () => clearInterval(interval);
  }, []);

  const checkAllServices = async () => {
    const results: ServiceStatus[] = [];

    for (const service of serviceChecks) {
      try {
        await service.check();
        results.push({
          name: service.name,
          status: 'healthy',
          healthy: true,
          message: '服務正常運行',
        });
      } catch (error) {
        results.push({
          name: service.name,
          status: 'unhealthy',
          healthy: false,
          message: error instanceof Error ? error.message : '服務不可用',
        });
      }
    }

    setServices(results);
    setLoading(false);
  };

  const getStatusColor = (healthy: boolean) => {
    return healthy ? 'bg-green-500' : 'bg-red-500';
  };

  const getStatusIcon = (healthy: boolean) => {
    return healthy ? '✓' : '✗';
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  const healthyCount = services.filter(s => s.healthy).length;
  const totalCount = services.length;

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">服務管理</h1>
        <Button onClick={checkAllServices}>刷新狀態</Button>
      </div>

      {/* 總覽卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">總服務數</h3>
          <p className="text-3xl font-bold mt-2">{totalCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康服務</h3>
          <p className="text-3xl font-bold mt-2 text-green-600">{healthyCount}</p>
        </Card>
        
        <Card className="p-4">
          <h3 className="text-sm font-medium text-gray-600">健康率</h3>
          <p className="text-3xl font-bold mt-2">
            {((healthyCount / totalCount) * 100).toFixed(0)}%
          </p>
        </Card>
      </div>

      {/* 服務列表 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">服務狀態</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {services.map((service) => (
            <div
              key={service.name}
              className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:shadow-md transition-shadow"
            >
              <div className="flex items-center space-x-4">
                <div className={`w-12 h-12 rounded-full ${getStatusColor(service.healthy)} flex items-center justify-center text-white text-xl font-bold`}>
                  {getStatusIcon(service.healthy)}
                </div>
                <div>
                  <h3 className="font-semibold text-lg">{service.name}</h3>
                  <p className="text-sm text-gray-600">{service.message}</p>
                </div>
              </div>
              <Badge variant={service.healthy ? 'default' : 'destructive'}>
                {service.status}
              </Badge>
            </div>
          ))}
        </div>
      </Card>

      {/* 快速操作 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">快速操作</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <Button variant="outline" className="h-20">
            查看 Prometheus
          </Button>
          <Button variant="outline" className="h-20">
            查看 Grafana
          </Button>
          <Button variant="outline" className="h-20">
            查看 Loki 日誌
          </Button>
          <Button variant="outline" className="h-20">
            量子作業管理
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default ServicesManagement;

