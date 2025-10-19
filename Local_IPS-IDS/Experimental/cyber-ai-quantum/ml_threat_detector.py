#!/usr/bin/env python3
"""
Pandora Box Console IDS-IPS - AI/ML 威脅檢測服務
使用深度學習和機器學習進行即時威脅檢測
"""

import numpy as np
import json
import logging
from datetime import datetime
from typing import Dict, List, Tuple, Optional
from dataclasses import dataclass, asdict
import asyncio
from enum import Enum

# 配置日誌
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class ThreatLevel(Enum):
    """威脅等級"""
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"
    CRITICAL = "critical"


class ThreatType(Enum):
    """威脅類型"""
    DDOS = "ddos"
    PORT_SCAN = "port_scan"
    BRUTE_FORCE = "brute_force"
    SQL_INJECTION = "sql_injection"
    XSS = "xss"
    MALWARE = "malware"
    RANSOMWARE = "ransomware"
    ZERO_DAY = "zero_day"
    APT = "advanced_persistent_threat"
    INSIDER = "insider_threat"


@dataclass
class ThreatDetection:
    """威脅檢測結果"""
    threat_id: str
    threat_type: ThreatType
    threat_level: ThreatLevel
    confidence: float
    source_ip: str
    target_ip: str
    timestamp: str
    features: Dict
    recommended_action: str
    

class DeepLearningThreatDetector:
    """深度學習威脅檢測器"""
    
    def __init__(self, model_path: Optional[str] = None):
        """初始化威脅檢測器"""
        self.model_path = model_path
        self.weights = self._initialize_weights()
        self.bias = self._initialize_bias()
        self.scaler_mean = np.array([0.5] * 20)
        self.scaler_std = np.array([0.25] * 20)
        logger.info("深度學習威脅檢測器已初始化")
    
    def _initialize_weights(self) -> List[np.ndarray]:
        """初始化神經網絡權重"""
        # 3層神經網絡: 20 -> 32 -> 16 -> 10 (威脅類型數)
        np.random.seed(42)
        weights = [
            np.random.randn(20, 32) * 0.1,  # 輸入層到隱藏層1
            np.random.randn(32, 16) * 0.1,  # 隱藏層1到隱藏層2
            np.random.randn(16, 10) * 0.1   # 隱藏層2到輸出層
        ]
        return weights
    
    def _initialize_bias(self) -> List[np.ndarray]:
        """初始化偏置"""
        return [
            np.zeros((1, 32)),
            np.zeros((1, 16)),
            np.zeros((1, 10))
        ]
    
    def _relu(self, x: np.ndarray) -> np.ndarray:
        """ReLU 激活函數"""
        return np.maximum(0, x)
    
    def _softmax(self, x: np.ndarray) -> np.ndarray:
        """Softmax 激活函數"""
        exp_x = np.exp(x - np.max(x, axis=1, keepdims=True))
        return exp_x / np.sum(exp_x, axis=1, keepdims=True)
    
    def _forward_pass(self, features: np.ndarray) -> Tuple[np.ndarray, float]:
        """前向傳播"""
        # 標準化特徵
        x = (features - self.scaler_mean) / self.scaler_std
        
        # 第一層
        z1 = np.dot(x, self.weights[0]) + self.bias[0]
        a1 = self._relu(z1)
        
        # 第二層
        z2 = np.dot(a1, self.weights[1]) + self.bias[1]
        a2 = self._relu(z2)
        
        # 輸出層
        z3 = np.dot(a2, self.weights[2]) + self.bias[2]
        output = self._softmax(z3)
        
        predicted_class = np.argmax(output)
        confidence = float(output[0, predicted_class])
        
        return predicted_class, confidence
    
    def extract_features(self, packet_data: Dict) -> np.ndarray:
        """從封包數據提取特徵"""
        features = []
        
        # 網路特徵 (0-9)
        features.append(packet_data.get('packet_size', 0) / 1500.0)
        features.append(packet_data.get('packets_per_second', 0) / 1000.0)
        features.append(packet_data.get('bytes_per_second', 0) / 1000000.0)
        features.append(packet_data.get('connection_count', 0) / 100.0)
        features.append(packet_data.get('unique_ips', 0) / 50.0)
        features.append(1.0 if packet_data.get('is_tcp', False) else 0.0)
        features.append(1.0 if packet_data.get('is_udp', False) else 0.0)
        features.append(packet_data.get('port_number', 0) / 65535.0)
        features.append(packet_data.get('ttl', 64) / 255.0)
        features.append(packet_data.get('window_size', 0) / 65535.0)
        
        # 行為特徵 (10-19)
        features.append(packet_data.get('syn_count', 0) / 100.0)
        features.append(packet_data.get('fin_count', 0) / 100.0)
        features.append(packet_data.get('rst_count', 0) / 100.0)
        features.append(packet_data.get('failed_logins', 0) / 10.0)
        features.append(packet_data.get('payload_entropy', 0))
        features.append(1.0 if packet_data.get('contains_shellcode', False) else 0.0)
        features.append(1.0 if packet_data.get('suspicious_pattern', False) else 0.0)
        features.append(packet_data.get('request_frequency', 0) / 100.0)
        features.append(packet_data.get('error_rate', 0))
        features.append(packet_data.get('anomaly_score', 0))
        
        return np.array(features).reshape(1, -1)
    
    async def detect_threat(self, packet_data: Dict) -> Optional[ThreatDetection]:
        """檢測威脅"""
        try:
            # 提取特徵
            features = self.extract_features(packet_data)
            
            # 預測
            threat_class, confidence = self._forward_pass(features)
            
            # 置信度閾值
            if confidence < 0.7:
                return None
            
            # 映射威脅類型
            threat_types = list(ThreatType)
            threat_type = threat_types[threat_class] if threat_class < len(threat_types) else ThreatType.MALWARE
            
            # 判定威脅等級
            if confidence >= 0.95:
                threat_level = ThreatLevel.CRITICAL
            elif confidence >= 0.85:
                threat_level = ThreatLevel.HIGH
            elif confidence >= 0.75:
                threat_level = ThreatLevel.MEDIUM
            else:
                threat_level = ThreatLevel.LOW
            
            # 建議操作
            actions = {
                ThreatLevel.CRITICAL: "立即阻斷並隔離",
                ThreatLevel.HIGH: "阻斷並記錄",
                ThreatLevel.MEDIUM: "監控並告警",
                ThreatLevel.LOW: "記錄觀察"
            }
            
            detection = ThreatDetection(
                threat_id=f"threat_{datetime.now().strftime('%Y%m%d%H%M%S')}_{np.random.randint(1000, 9999)}",
                threat_type=threat_type,
                threat_level=threat_level,
                confidence=confidence,
                source_ip=packet_data.get('source_ip', 'unknown'),
                target_ip=packet_data.get('target_ip', 'unknown'),
                timestamp=datetime.now().isoformat(),
                features=packet_data,
                recommended_action=actions[threat_level]
            )
            
            logger.info(f"檢測到威脅: {detection.threat_type.value} - 置信度: {confidence:.2%}")
            return detection
            
        except Exception as e:
            logger.error(f"威脅檢測錯誤: {e}")
            return None


class BehaviorAnalyzer:
    """行為分析器"""
    
    def __init__(self):
        """初始化行為分析器"""
        self.baselines = {}
        self.anomaly_threshold = 2.5
        logger.info("行為分析器已初始化")
    
    def create_baseline(self, user_id: str, behaviors: List[Dict]):
        """建立用戶行為基線"""
        if not behaviors:
            return
        
        # 計算平均值和標準差
        features = ['login_time', 'session_duration', 'data_accessed', 'commands_executed']
        baseline = {}
        
        for feature in features:
            values = [b.get(feature, 0) for b in behaviors]
            baseline[feature] = {
                'mean': np.mean(values),
                'std': np.std(values)
            }
        
        self.baselines[user_id] = baseline
        logger.info(f"已建立用戶 {user_id} 的行為基線")
    
    def detect_anomaly(self, user_id: str, behavior: Dict) -> Tuple[bool, float]:
        """檢測異常行為"""
        if user_id not in self.baselines:
            return False, 0.0
        
        baseline = self.baselines[user_id]
        anomaly_score = 0.0
        count = 0
        
        for feature, stats in baseline.items():
            value = behavior.get(feature, 0)
            z_score = abs((value - stats['mean']) / (stats['std'] + 1e-6))
            anomaly_score += z_score
            count += 1
        
        anomaly_score /= count
        is_anomaly = anomaly_score > self.anomaly_threshold
        
        if is_anomaly:
            logger.warning(f"檢測到用戶 {user_id} 的異常行為 - 分數: {anomaly_score:.2f}")
        
        return is_anomaly, anomaly_score


class AISecurityMonitor:
    """AI 安全監控器"""
    
    def __init__(self):
        """初始化 AI 安全監控器"""
        self.model_integrity_checks = []
        self.adversarial_detections = []
        logger.info("AI 安全監控器已初始化")
    
    async def check_model_integrity(self, model_hash: str) -> bool:
        """檢查模型完整性"""
        # 模擬模型完整性檢查
        expected_hash = "abc123def456"
        is_valid = model_hash == expected_hash
        
        self.model_integrity_checks.append({
            'timestamp': datetime.now().isoformat(),
            'is_valid': is_valid,
            'hash': model_hash
        })
        
        if not is_valid:
            logger.error(f"模型完整性檢查失敗: {model_hash}")
        
        return is_valid
    
    async def detect_adversarial_attack(self, input_data: np.ndarray) -> Tuple[bool, float]:
        """檢測對抗性攻擊"""
        # 簡單的對抗性攻擊檢測
        # 檢查異常的特徵值範圍和模式
        
        # 檢查範圍
        out_of_range = np.sum((input_data < -3) | (input_data > 3))
        out_of_range_ratio = out_of_range / input_data.size
        
        # 檢查突然的變化
        if input_data.ndim > 1:
            gradients = np.diff(input_data, axis=1)
            high_gradient = np.sum(np.abs(gradients) > 2)
            high_gradient_ratio = high_gradient / gradients.size if gradients.size > 0 else 0
        else:
            high_gradient_ratio = 0
        
        adversarial_score = (out_of_range_ratio + high_gradient_ratio) / 2
        is_adversarial = adversarial_score > 0.3
        
        if is_adversarial:
            self.adversarial_detections.append({
                'timestamp': datetime.now().isoformat(),
                'score': adversarial_score
            })
            logger.warning(f"檢測到對抗性攻擊 - 分數: {adversarial_score:.2%}")
        
        return is_adversarial, adversarial_score


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora AI/ML 威脅檢測系統啟動 ===")
    
    # 初始化檢測器
    threat_detector = DeepLearningThreatDetector()
    behavior_analyzer = BehaviorAnalyzer()
    security_monitor = AISecurityMonitor()
    
    # 模擬封包數據
    test_packets = [
        {
            'source_ip': '192.168.1.100',
            'target_ip': '10.0.0.1',
            'packet_size': 1200,
            'packets_per_second': 500,
            'bytes_per_second': 600000,
            'connection_count': 50,
            'unique_ips': 10,
            'is_tcp': True,
            'port_number': 80,
            'syn_count': 45,
            'payload_entropy': 0.8,
            'request_frequency': 50
        },
        {
            'source_ip': '192.168.1.101',
            'target_ip': '10.0.0.2',
            'packet_size': 64,
            'packets_per_second': 1000,
            'bytes_per_second': 64000,
            'connection_count': 100,
            'unique_ips': 50,
            'is_tcp': True,
            'port_number': 22,
            'failed_logins': 8,
            'request_frequency': 100,
            'anomaly_score': 0.9
        }
    ]
    
    # 檢測威脅
    logger.info("\n--- 威脅檢測 ---")
    for packet in test_packets:
        detection = await threat_detector.detect_threat(packet)
        if detection:
            print(json.dumps(asdict(detection), indent=2, default=str))
    
    # 行為分析
    logger.info("\n--- 行為分析 ---")
    user_behaviors = [
        {'login_time': 9, 'session_duration': 480, 'data_accessed': 100, 'commands_executed': 50},
        {'login_time': 10, 'session_duration': 460, 'data_accessed': 120, 'commands_executed': 45},
        {'login_time': 9, 'session_duration': 490, 'data_accessed': 110, 'commands_executed': 55}
    ]
    behavior_analyzer.create_baseline('user001', user_behaviors)
    
    # 測試異常行為
    anomaly_behavior = {'login_time': 2, 'session_duration': 120, 'data_accessed': 1000, 'commands_executed': 200}
    is_anomaly, score = behavior_analyzer.detect_anomaly('user001', anomaly_behavior)
    print(f"異常檢測: {is_anomaly}, 分數: {score:.2f}")
    
    # AI 安全檢查
    logger.info("\n--- AI 安全檢查 ---")
    model_valid = await security_monitor.check_model_integrity("abc123def456")
    print(f"模型完整性: {model_valid}")
    
    test_input = np.random.randn(1, 20) * 5  # 高變異輸入
    is_adv, adv_score = await security_monitor.detect_adversarial_attack(test_input)
    print(f"對抗性攻擊: {is_adv}, 分數: {adv_score:.2%}")
    
    logger.info("\n=== 系統測試完成 ===")


if __name__ == "__main__":
    asyncio.run(main())

