#!/usr/bin/env python3
"""
Pandora Box Console IDS-IPS - AI 資料流監控與分析
使用 AI 監控和分析網路資料流
"""

import numpy as np
import logging
from datetime import datetime, timedelta
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict
from collections import defaultdict
import asyncio
import json

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@dataclass
class DataFlowMetrics:
    """資料流指標"""
    timestamp: str
    total_bytes: int
    packet_count: int
    unique_sources: int
    unique_destinations: int
    protocols: Dict[str, int]
    avg_packet_size: float
    flow_rate: float


@dataclass
class AnomalyAlert:
    """異常告警"""
    alert_id: str
    timestamp: str
    anomaly_type: str
    severity: str
    source_ip: str
    anomaly_score: float
    baseline_deviation: float
    recommended_action: str


class NetworkFlowAnalyzer:
    """網路流量分析器"""
    
    def __init__(self, window_size: int = 60):
        """初始化流量分析器"""
        self.window_size = window_size  # 秒
        self.flow_buffer = []
        self.stats_history = []
        logger.info(f"網路流量分析器已初始化 (窗口: {window_size}s)")
    
    def process_packet(self, packet: Dict):
        """處理單個封包"""
        packet['processed_at'] = datetime.now()
        self.flow_buffer.append(packet)
        
        # 清理過期數據
        cutoff_time = datetime.now() - timedelta(seconds=self.window_size)
        self.flow_buffer = [
            p for p in self.flow_buffer 
            if p['processed_at'] > cutoff_time
        ]
    
    def calculate_metrics(self) -> Optional[DataFlowMetrics]:
        """計算當前窗口的流量指標"""
        if not self.flow_buffer:
            return None
        
        total_bytes = sum(p.get('size', 0) for p in self.flow_buffer)
        packet_count = len(self.flow_buffer)
        
        sources = set(p.get('source_ip', '') for p in self.flow_buffer)
        destinations = set(p.get('dest_ip', '') for p in self.flow_buffer)
        
        protocols = defaultdict(int)
        for p in self.flow_buffer:
            proto = p.get('protocol', 'unknown')
            protocols[proto] += 1
        
        avg_packet_size = total_bytes / packet_count if packet_count > 0 else 0
        flow_rate = total_bytes / self.window_size  # bytes/second
        
        metrics = DataFlowMetrics(
            timestamp=datetime.now().isoformat(),
            total_bytes=total_bytes,
            packet_count=packet_count,
            unique_sources=len(sources),
            unique_destinations=len(destinations),
            protocols=dict(protocols),
            avg_packet_size=avg_packet_size,
            flow_rate=flow_rate
        )
        
        self.stats_history.append(metrics)
        
        # 保留最近1小時的歷史
        cutoff_time = datetime.now() - timedelta(hours=1)
        self.stats_history = [
            m for m in self.stats_history
            if datetime.fromisoformat(m.timestamp) > cutoff_time
        ]
        
        return metrics


class DataFlowAnomalyDetector:
    """資料流異常檢測器"""
    
    def __init__(self):
        """初始化異常檢測器"""
        self.baselines = {}
        self.anomalies = []
        self.threshold = 3.0  # 3個標準差
        logger.info("資料流異常檢測器已初始化")
    
    def update_baseline(self, metrics_history: List[DataFlowMetrics]):
        """更新基線"""
        if len(metrics_history) < 10:
            logger.warning("數據不足，無法更新基線")
            return
        
        # 提取數值特徵
        flow_rates = [m.flow_rate for m in metrics_history]
        packet_counts = [m.packet_count for m in metrics_history]
        avg_sizes = [m.avg_packet_size for m in metrics_history]
        
        self.baselines = {
            'flow_rate': {
                'mean': np.mean(flow_rates),
                'std': np.std(flow_rates)
            },
            'packet_count': {
                'mean': np.mean(packet_counts),
                'std': np.std(packet_counts)
            },
            'avg_packet_size': {
                'mean': np.mean(avg_sizes),
                'std': np.std(avg_sizes)
            }
        }
        
        logger.info("基線已更新")
    
    def detect_anomaly(self, current_metrics: DataFlowMetrics) -> Optional[AnomalyAlert]:
        """檢測異常"""
        if not self.baselines:
            return None
        
        anomalies_detected = []
        max_deviation = 0.0
        
        # 檢查流量率異常
        flow_rate = current_metrics.flow_rate
        baseline = self.baselines['flow_rate']
        z_score = abs((flow_rate - baseline['mean']) / (baseline['std'] + 1e-6))
        
        if z_score > self.threshold:
            anomalies_detected.append('flow_rate')
            max_deviation = max(max_deviation, z_score)
        
        # 檢查封包數異常
        packet_count = current_metrics.packet_count
        baseline = self.baselines['packet_count']
        z_score = abs((packet_count - baseline['mean']) / (baseline['std'] + 1e-6))
        
        if z_score > self.threshold:
            anomalies_detected.append('packet_count')
            max_deviation = max(max_deviation, z_score)
        
        # 檢查平均封包大小異常
        avg_size = current_metrics.avg_packet_size
        baseline = self.baselines['avg_packet_size']
        z_score = abs((avg_size - baseline['mean']) / (baseline['std'] + 1e-6))
        
        if z_score > self.threshold:
            anomalies_detected.append('avg_packet_size')
            max_deviation = max(max_deviation, z_score)
        
        if not anomalies_detected:
            return None
        
        # 判定嚴重程度
        if max_deviation > 5.0:
            severity = 'critical'
            action = '立即阻斷並調查'
        elif max_deviation > 4.0:
            severity = 'high'
            action = '告警並監控'
        elif max_deviation > 3.0:
            severity = 'medium'
            action = '記錄並觀察'
        else:
            severity = 'low'
            action = '持續監控'
        
        alert = AnomalyAlert(
            alert_id=f"anomaly_{datetime.now().strftime('%Y%m%d%H%M%S')}",
            timestamp=datetime.now().isoformat(),
            anomaly_type=', '.join(anomalies_detected),
            severity=severity,
            source_ip='multiple',
            anomaly_score=max_deviation,
            baseline_deviation=max_deviation,
            recommended_action=action
        )
        
        self.anomalies.append(alert)
        logger.warning(f"檢測到異常: {alert.anomaly_type} - 偏差: {max_deviation:.2f}σ")
        
        return alert


class DataFlowMonitor:
    """資料流監控器"""
    
    def __init__(self):
        """初始化資料流監控器"""
        self.analyzer = NetworkFlowAnalyzer(window_size=60)
        self.anomaly_detector = DataFlowAnomalyDetector()
        self.running = False
        logger.info("資料流監控器已初始化")
    
    async def start_monitoring(self):
        """啟動監控"""
        self.running = True
        logger.info("資料流監控已啟動")
        
        while self.running:
            # 計算當前指標
            metrics = self.analyzer.calculate_metrics()
            
            if metrics:
                # 檢測異常
                anomaly = self.anomaly_detector.detect_anomaly(metrics)
                
                if anomaly:
                    logger.warning(f"異常告警: {anomaly.alert_id}")
                    # 這裡可以發送到 RabbitMQ 或告警系統
            
            # 每10秒更新基線
            if len(self.analyzer.stats_history) >= 10:
                self.anomaly_detector.update_baseline(self.analyzer.stats_history)
            
            await asyncio.sleep(10)
    
    def stop_monitoring(self):
        """停止監控"""
        self.running = False
        logger.info("資料流監控已停止")


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora 資料流監控系統啟動 ===")
    
    # 初始化監控器
    monitor = DataFlowMonitor()
    
    # 模擬一些封包
    logger.info("\n--- 模擬網路流量 ---")
    for i in range(50):
        packet = {
            'source_ip': f'192.168.1.{i % 20}',
            'dest_ip': f'10.0.0.{i % 5}',
            'size': np.random.randint(64, 1500),
            'protocol': np.random.choice(['TCP', 'UDP', 'ICMP'])
        }
        monitor.analyzer.process_packet(packet)
    
    # 計算指標
    metrics = monitor.analyzer.calculate_metrics()
    if metrics:
        print("\n當前資料流指標:")
        print(json.dumps(asdict(metrics), indent=2, default=str))
    
    # 更新基線
    logger.info("\n--- 建立基線 ---")
    # 添加更多歷史數據
    for _ in range(3):
        for i in range(50):
            packet = {
                'source_ip': f'192.168.1.{i % 20}',
                'dest_ip': f'10.0.0.{i % 5}',
                'size': np.random.randint(64, 1500),
                'protocol': np.random.choice(['TCP', 'UDP', 'ICMP'])
            }
            monitor.analyzer.process_packet(packet)
        
        m = monitor.analyzer.calculate_metrics()
        if m:
            monitor.analyzer.stats_history.append(m)
    
    monitor.anomaly_detector.update_baseline(monitor.analyzer.stats_history)
    
    # 模擬異常流量
    logger.info("\n--- 注入異常流量 ---")
    for i in range(500):  # 大量流量
        packet = {
            'source_ip': '192.168.1.100',
            'dest_ip': '10.0.0.1',
            'size': 64,  # 小封包
            'protocol': 'TCP'
        }
        monitor.analyzer.process_packet(packet)
    
    # 檢測異常
    metrics = monitor.analyzer.calculate_metrics()
    if metrics:
        anomaly = monitor.anomaly_detector.detect_anomaly(metrics)
        if anomaly:
            print("\n異常告警:")
            print(json.dumps(asdict(anomaly), indent=2, default=str))
    
    logger.info("\n=== 系統測試完成 ===")


if __name__ == "__main__":
    asyncio.run(main())

