import React, { useState, useEffect } from 'react';
import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { nginxAPI } from '../../services/axiom-api';

interface NginxConfig {
  config: string;
  config_path: string;
  last_modified: string;
  size: number;
  valid: boolean;
}

const NginxConfigEditor: React.FC = () => {
  const [config, setConfig] = useState<NginxConfig | null>(null);
  const [editedConfig, setEditedConfig] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [reloading, setReloading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    setLoading(true);
    try {
      const result = await nginxAPI.getConfig() as NginxConfig;
      setConfig(result);
      setEditedConfig(result.config);
    } catch (error) {
      alert('載入配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!editedConfig.trim()) {
      alert('配置不能為空');
      return;
    }

    setSaving(true);
    try {
      await nginxAPI.updateConfig(editedConfig, true, true);
      alert('配置已保存並驗證成功');
      setIsEditing(false);
      loadConfig();
    } catch (error) {
      alert('保存配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setSaving(false);
    }
  };

  const handleReload = async () => {
    if (!confirm('確定要重載 Nginx 配置嗎？')) {
      return;
    }

    setReloading(true);
    try {
      await nginxAPI.reload();
      alert('Nginx 配置已重載');
    } catch (error) {
      alert('重載失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setReloading(false);
    }
  };

  const handleCancel = () => {
    setEditedConfig(config?.config || '');
    setIsEditing(false);
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Nginx 配置管理</h1>
        <div className="flex space-x-2">
          {!isEditing ? (
            <>
              <Button onClick={() => setIsEditing(true)}>編輯配置</Button>
              <Button onClick={handleReload} disabled={reloading} variant="outline">
                {reloading ? '重載中...' : '重載 Nginx'}
              </Button>
            </>
          ) : (
            <>
              <Button onClick={handleSave} disabled={saving}>
                {saving ? '保存中...' : '保存配置'}
              </Button>
              <Button onClick={handleCancel} variant="outline">取消</Button>
            </>
          )}
        </div>
      </div>

      {/* 配置信息 */}
      {config && (
        <Card className="p-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-600">配置路徑:</span>
              <p className="font-mono">{config.config_path}</p>
            </div>
            <div>
              <span className="text-gray-600">最後修改:</span>
              <p>{new Date(config.last_modified).toLocaleString('zh-TW')}</p>
            </div>
            <div>
              <span className="text-gray-600">文件大小:</span>
              <p>{(config.size / 1024).toFixed(2)} KB</p>
            </div>
            <div>
              <span className="text-gray-600">狀態:</span>
              <p>
                <Badge variant={config.valid ? 'default' : 'destructive'}>
                  {config.valid ? '有效' : '無效'}
                </Badge>
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* 配置編輯器 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">配置內容</h2>
        <div className="relative">
          <textarea
            value={editedConfig}
            onChange={(e) => setEditedConfig(e.target.value)}
            readOnly={!isEditing}
            className={`w-full h-[600px] font-mono text-sm p-4 border border-gray-300 rounded-md ${
              !isEditing ? 'bg-gray-50' : 'bg-white'
            }`}
            style={{ tabSize: 2 }}
          />
          {!isEditing && (
            <div className="absolute inset-0 bg-gray-100 bg-opacity-50 rounded-md" />
          )}
        </div>

        {isEditing && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              ⚠️ 配置將在保存時自動驗證。無效的配置將被拒絕。建議在修改前備份。
            </p>
          </div>
        )}
      </Card>
    </div>
  );
};

export default NginxConfigEditor;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { nginxAPI } from '../../services/axiom-api';

interface NginxConfig {
  config: string;
  config_path: string;
  last_modified: string;
  size: number;
  valid: boolean;
}

const NginxConfigEditor: React.FC = () => {
  const [config, setConfig] = useState<NginxConfig | null>(null);
  const [editedConfig, setEditedConfig] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [reloading, setReloading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    setLoading(true);
    try {
      const result = await nginxAPI.getConfig() as NginxConfig;
      setConfig(result);
      setEditedConfig(result.config);
    } catch (error) {
      alert('載入配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!editedConfig.trim()) {
      alert('配置不能為空');
      return;
    }

    setSaving(true);
    try {
      await nginxAPI.updateConfig(editedConfig, true, true);
      alert('配置已保存並驗證成功');
      setIsEditing(false);
      loadConfig();
    } catch (error) {
      alert('保存配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setSaving(false);
    }
  };

  const handleReload = async () => {
    if (!confirm('確定要重載 Nginx 配置嗎？')) {
      return;
    }

    setReloading(true);
    try {
      await nginxAPI.reload();
      alert('Nginx 配置已重載');
    } catch (error) {
      alert('重載失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setReloading(false);
    }
  };

  const handleCancel = () => {
    setEditedConfig(config?.config || '');
    setIsEditing(false);
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Nginx 配置管理</h1>
        <div className="flex space-x-2">
          {!isEditing ? (
            <>
              <Button onClick={() => setIsEditing(true)}>編輯配置</Button>
              <Button onClick={handleReload} disabled={reloading} variant="outline">
                {reloading ? '重載中...' : '重載 Nginx'}
              </Button>
            </>
          ) : (
            <>
              <Button onClick={handleSave} disabled={saving}>
                {saving ? '保存中...' : '保存配置'}
              </Button>
              <Button onClick={handleCancel} variant="outline">取消</Button>
            </>
          )}
        </div>
      </div>

      {/* 配置信息 */}
      {config && (
        <Card className="p-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-600">配置路徑:</span>
              <p className="font-mono">{config.config_path}</p>
            </div>
            <div>
              <span className="text-gray-600">最後修改:</span>
              <p>{new Date(config.last_modified).toLocaleString('zh-TW')}</p>
            </div>
            <div>
              <span className="text-gray-600">文件大小:</span>
              <p>{(config.size / 1024).toFixed(2)} KB</p>
            </div>
            <div>
              <span className="text-gray-600">狀態:</span>
              <p>
                <Badge variant={config.valid ? 'default' : 'destructive'}>
                  {config.valid ? '有效' : '無效'}
                </Badge>
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* 配置編輯器 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">配置內容</h2>
        <div className="relative">
          <textarea
            value={editedConfig}
            onChange={(e) => setEditedConfig(e.target.value)}
            readOnly={!isEditing}
            className={`w-full h-[600px] font-mono text-sm p-4 border border-gray-300 rounded-md ${
              !isEditing ? 'bg-gray-50' : 'bg-white'
            }`}
            style={{ tabSize: 2 }}
          />
          {!isEditing && (
            <div className="absolute inset-0 bg-gray-100 bg-opacity-50 rounded-md" />
          )}
        </div>

        {isEditing && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              ⚠️ 配置將在保存時自動驗證。無效的配置將被拒絕。建議在修改前備份。
            </p>
          </div>
        )}
      </Card>
    </div>
  );
};

export default NginxConfigEditor;

import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { nginxAPI } from '../../services/axiom-api';

interface NginxConfig {
  config: string;
  config_path: string;
  last_modified: string;
  size: number;
  valid: boolean;
}

const NginxConfigEditor: React.FC = () => {
  const [config, setConfig] = useState<NginxConfig | null>(null);
  const [editedConfig, setEditedConfig] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [reloading, setReloading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    setLoading(true);
    try {
      const result = await nginxAPI.getConfig() as NginxConfig;
      setConfig(result);
      setEditedConfig(result.config);
    } catch (error) {
      alert('載入配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!editedConfig.trim()) {
      alert('配置不能為空');
      return;
    }

    setSaving(true);
    try {
      await nginxAPI.updateConfig(editedConfig, true, true);
      alert('配置已保存並驗證成功');
      setIsEditing(false);
      loadConfig();
    } catch (error) {
      alert('保存配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setSaving(false);
    }
  };

  const handleReload = async () => {
    if (!confirm('確定要重載 Nginx 配置嗎？')) {
      return;
    }

    setReloading(true);
    try {
      await nginxAPI.reload();
      alert('Nginx 配置已重載');
    } catch (error) {
      alert('重載失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setReloading(false);
    }
  };

  const handleCancel = () => {
    setEditedConfig(config?.config || '');
    setIsEditing(false);
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Nginx 配置管理</h1>
        <div className="flex space-x-2">
          {!isEditing ? (
            <>
              <Button onClick={() => setIsEditing(true)}>編輯配置</Button>
              <Button onClick={handleReload} disabled={reloading} variant="outline">
                {reloading ? '重載中...' : '重載 Nginx'}
              </Button>
            </>
          ) : (
            <>
              <Button onClick={handleSave} disabled={saving}>
                {saving ? '保存中...' : '保存配置'}
              </Button>
              <Button onClick={handleCancel} variant="outline">取消</Button>
            </>
          )}
        </div>
      </div>

      {/* 配置信息 */}
      {config && (
        <Card className="p-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-600">配置路徑:</span>
              <p className="font-mono">{config.config_path}</p>
            </div>
            <div>
              <span className="text-gray-600">最後修改:</span>
              <p>{new Date(config.last_modified).toLocaleString('zh-TW')}</p>
            </div>
            <div>
              <span className="text-gray-600">文件大小:</span>
              <p>{(config.size / 1024).toFixed(2)} KB</p>
            </div>
            <div>
              <span className="text-gray-600">狀態:</span>
              <p>
                <Badge variant={config.valid ? 'default' : 'destructive'}>
                  {config.valid ? '有效' : '無效'}
                </Badge>
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* 配置編輯器 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">配置內容</h2>
        <div className="relative">
          <textarea
            value={editedConfig}
            onChange={(e) => setEditedConfig(e.target.value)}
            readOnly={!isEditing}
            className={`w-full h-[600px] font-mono text-sm p-4 border border-gray-300 rounded-md ${
              !isEditing ? 'bg-gray-50' : 'bg-white'
            }`}
            style={{ tabSize: 2 }}
          />
          {!isEditing && (
            <div className="absolute inset-0 bg-gray-100 bg-opacity-50 rounded-md" />
          )}
        </div>

        {isEditing && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              ⚠️ 配置將在保存時自動驗證。無效的配置將被拒絕。建議在修改前備份。
            </p>
          </div>
        )}
      </Card>
    </div>
  );
};

export default NginxConfigEditor;


import { Card } from '../ui/card';
import { Button } from '../ui/button';
import { nginxAPI } from '../../services/axiom-api';

interface NginxConfig {
  config: string;
  config_path: string;
  last_modified: string;
  size: number;
  valid: boolean;
}

const NginxConfigEditor: React.FC = () => {
  const [config, setConfig] = useState<NginxConfig | null>(null);
  const [editedConfig, setEditedConfig] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [reloading, setReloading] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    setLoading(true);
    try {
      const result = await nginxAPI.getConfig() as NginxConfig;
      setConfig(result);
      setEditedConfig(result.config);
    } catch (error) {
      alert('載入配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!editedConfig.trim()) {
      alert('配置不能為空');
      return;
    }

    setSaving(true);
    try {
      await nginxAPI.updateConfig(editedConfig, true, true);
      alert('配置已保存並驗證成功');
      setIsEditing(false);
      loadConfig();
    } catch (error) {
      alert('保存配置失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setSaving(false);
    }
  };

  const handleReload = async () => {
    if (!confirm('確定要重載 Nginx 配置嗎？')) {
      return;
    }

    setReloading(true);
    try {
      await nginxAPI.reload();
      alert('Nginx 配置已重載');
    } catch (error) {
      alert('重載失敗: ' + (error instanceof Error ? error.message : 'Unknown error'));
    } finally {
      setReloading(false);
    }
  };

  const handleCancel = () => {
    setEditedConfig(config?.config || '');
    setIsEditing(false);
  };

  if (loading) {
    return <div className="flex items-center justify-center h-64">載入中...</div>;
  }

  return (
    <div className="space-y-6 p-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Nginx 配置管理</h1>
        <div className="flex space-x-2">
          {!isEditing ? (
            <>
              <Button onClick={() => setIsEditing(true)}>編輯配置</Button>
              <Button onClick={handleReload} disabled={reloading} variant="outline">
                {reloading ? '重載中...' : '重載 Nginx'}
              </Button>
            </>
          ) : (
            <>
              <Button onClick={handleSave} disabled={saving}>
                {saving ? '保存中...' : '保存配置'}
              </Button>
              <Button onClick={handleCancel} variant="outline">取消</Button>
            </>
          )}
        </div>
      </div>

      {/* 配置信息 */}
      {config && (
        <Card className="p-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-600">配置路徑:</span>
              <p className="font-mono">{config.config_path}</p>
            </div>
            <div>
              <span className="text-gray-600">最後修改:</span>
              <p>{new Date(config.last_modified).toLocaleString('zh-TW')}</p>
            </div>
            <div>
              <span className="text-gray-600">文件大小:</span>
              <p>{(config.size / 1024).toFixed(2)} KB</p>
            </div>
            <div>
              <span className="text-gray-600">狀態:</span>
              <p>
                <Badge variant={config.valid ? 'default' : 'destructive'}>
                  {config.valid ? '有效' : '無效'}
                </Badge>
              </p>
            </div>
          </div>
        </Card>
      )}

      {/* 配置編輯器 */}
      <Card className="p-6">
        <h2 className="text-xl font-bold mb-4">配置內容</h2>
        <div className="relative">
          <textarea
            value={editedConfig}
            onChange={(e) => setEditedConfig(e.target.value)}
            readOnly={!isEditing}
            className={`w-full h-[600px] font-mono text-sm p-4 border border-gray-300 rounded-md ${
              !isEditing ? 'bg-gray-50' : 'bg-white'
            }`}
            style={{ tabSize: 2 }}
          />
          {!isEditing && (
            <div className="absolute inset-0 bg-gray-100 bg-opacity-50 rounded-md" />
          )}
        </div>

        {isEditing && (
          <div className="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-md">
            <p className="text-sm text-yellow-800">
              ⚠️ 配置將在保存時自動驗證。無效的配置將被拒絕。建議在修改前備份。
            </p>
          </div>
        )}
      </Card>
    </div>
  );
};

export default NginxConfigEditor;

