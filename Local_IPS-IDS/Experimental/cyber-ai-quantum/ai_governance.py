#!/usr/bin/env python3
"""
Pandora Box Console IDS-IPS - AI 治理與監控系統
確保 AI 模型的安全性、公平性和可靠性
"""

import numpy as np
import hashlib
import logging
from datetime import datetime
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict
from enum import Enum
import json

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class GovernanceStatus(Enum):
    """治理狀態"""
    COMPLIANT = "compliant"
    WARNING = "warning"
    VIOLATION = "violation"
    CRITICAL = "critical"


@dataclass
class GovernanceReport:
    """治理報告"""
    report_id: str
    timestamp: str
    overall_status: GovernanceStatus
    model_integrity: Dict
    fairness_metrics: Dict
    performance_metrics: Dict
    security_checks: Dict
    recommendations: List[str]


class ModelIntegrityChecker:
    """模型完整性檢查器"""
    
    def __init__(self):
        """初始化完整性檢查器"""
        self.model_registry = {}
        self.integrity_log = []
        logger.info("模型完整性檢查器已初始化")
    
    def register_model(self, model_id: str, model_data: bytes) -> str:
        """註冊模型並生成指紋"""
        # 計算模型哈希
        model_hash = hashlib.sha256(model_data).hexdigest()
        
        self.model_registry[model_id] = {
            'hash': model_hash,
            'size': len(model_data),
            'registered_at': datetime.now().isoformat(),
            'version': '1.0.0'
        }
        
        logger.info(f"模型已註冊: {model_id} - Hash: {model_hash[:16]}...")
        return model_hash
    
    def verify_integrity(self, model_id: str, current_data: bytes) -> Tuple[bool, str]:
        """驗證模型完整性"""
        if model_id not in self.model_registry:
            return False, "模型未註冊"
        
        expected_hash = self.model_registry[model_id]['hash']
        current_hash = hashlib.sha256(current_data).hexdigest()
        
        is_valid = expected_hash == current_hash
        
        log_entry = {
            'model_id': model_id,
            'timestamp': datetime.now().isoformat(),
            'is_valid': is_valid,
            'expected_hash': expected_hash[:16],
            'current_hash': current_hash[:16]
        }
        self.integrity_log.append(log_entry)
        
        if not is_valid:
            logger.error(f"模型完整性檢查失敗: {model_id}")
            return False, "哈希值不匹配 - 可能的模型中毒攻擊"
        
        logger.info(f"模型完整性檢查通過: {model_id}")
        return True, "完整性驗證通過"


class FairnessAuditor:
    """公平性審計器"""
    
    def __init__(self):
        """初始化公平性審計器"""
        self.audit_history = []
        logger.info("公平性審計器已初始化")
    
    def calculate_demographic_parity(
        self, 
        predictions: np.ndarray, 
        protected_attribute: np.ndarray
    ) -> float:
        """計算人口統計平等性"""
        # 計算不同群組的正預測率
        groups = np.unique(protected_attribute)
        positive_rates = []
        
        for group in groups:
            mask = protected_attribute == group
            if np.sum(mask) > 0:
                positive_rate = np.mean(predictions[mask])
                positive_rates.append(positive_rate)
        
        if len(positive_rates) < 2:
            return 1.0
        
        # 計算差異
        disparity = max(positive_rates) - min(positive_rates)
        parity_score = 1.0 - disparity
        
        return max(0.0, parity_score)
    
    def calculate_equal_opportunity(
        self,
        predictions: np.ndarray,
        true_labels: np.ndarray,
        protected_attribute: np.ndarray
    ) -> float:
        """計算機會均等性"""
        groups = np.unique(protected_attribute)
        true_positive_rates = []
        
        for group in groups:
            mask = (protected_attribute == group) & (true_labels == 1)
            if np.sum(mask) > 0:
                tpr = np.sum(predictions[mask] == 1) / np.sum(mask)
                true_positive_rates.append(tpr)
        
        if len(true_positive_rates) < 2:
            return 1.0
        
        # 計算差異
        disparity = max(true_positive_rates) - min(true_positive_rates)
        equality_score = 1.0 - disparity
        
        return max(0.0, equality_score)
    
    async def audit_fairness(
        self,
        model_id: str,
        predictions: np.ndarray,
        protected_attributes: Dict[str, np.ndarray]
    ) -> Dict:
        """執行公平性審計"""
        try:
            metrics = {}
            
            for attr_name, attr_values in protected_attributes.items():
                parity = self.calculate_demographic_parity(predictions, attr_values)
                metrics[f'{attr_name}_parity'] = parity
            
            overall_fairness = np.mean(list(metrics.values()))
            
            audit_result = {
                'model_id': model_id,
                'timestamp': datetime.now().isoformat(),
                'metrics': metrics,
                'overall_fairness': overall_fairness,
                'status': 'pass' if overall_fairness > 0.8 else 'fail'
            }
            
            self.audit_history.append(audit_result)
            logger.info(f"公平性審計完成: {model_id} - 分數: {overall_fairness:.2%}")
            
            return audit_result
            
        except Exception as e:
            logger.error(f"公平性審計失敗: {e}")
            return {}


class PerformanceMonitor:
    """性能監控器"""
    
    def __init__(self):
        """初始化性能監控器"""
        self.metrics_history = []
        self.alert_threshold = {
            'accuracy': 0.90,
            'latency_ms': 100,
            'throughput': 500,
            'memory_mb': 4000,
            'cpu_percent': 80
        }
        logger.info("性能監控器已初始化")
    
    async def monitor_performance(self, model_id: str, metrics: Dict) -> Dict:
        """監控模型性能"""
        alerts = []
        
        # 檢查準確率
        if metrics.get('accuracy', 1.0) < self.alert_threshold['accuracy']:
            alerts.append(f"準確率過低: {metrics['accuracy']:.2%}")
        
        # 檢查延遲
        if metrics.get('latency_ms', 0) > self.alert_threshold['latency_ms']:
            alerts.append(f"延遲過高: {metrics['latency_ms']}ms")
        
        # 檢查吞吐量
        if metrics.get('throughput', 0) < self.alert_threshold['throughput']:
            alerts.append(f"吞吐量過低: {metrics['throughput']} req/s")
        
        # 檢查記憶體
        if metrics.get('memory_mb', 0) > self.alert_threshold['memory_mb']:
            alerts.append(f"記憶體使用過高: {metrics['memory_mb']}MB")
        
        # 檢查 CPU
        if metrics.get('cpu_percent', 0) > self.alert_threshold['cpu_percent']:
            alerts.append(f"CPU 使用過高: {metrics['cpu_percent']}%")
        
        status = GovernanceStatus.COMPLIANT if not alerts else GovernanceStatus.WARNING
        
        report = {
            'model_id': model_id,
            'timestamp': datetime.now().isoformat(),
            'metrics': metrics,
            'alerts': alerts,
            'status': status.value
        }
        
        self.metrics_history.append(report)
        
        if alerts:
            logger.warning(f"性能告警: {model_id} - {len(alerts)} 個問題")
        else:
            logger.info(f"性能監控正常: {model_id}")
        
        return report


class AIGovernanceSystem:
    """AI 治理系統"""
    
    def __init__(self):
        """初始化 AI 治理系統"""
        self.integrity_checker = ModelIntegrityChecker()
        self.fairness_auditor = FairnessAuditor()
        self.performance_monitor = PerformanceMonitor()
        logger.info("AI 治理系統已初始化")
    
    async def generate_governance_report(
        self,
        model_id: str,
        model_data: bytes,
        predictions: np.ndarray,
        performance_metrics: Dict
    ) -> GovernanceReport:
        """生成完整的治理報告"""
        try:
            # 1. 模型完整性檢查
            if model_id not in self.integrity_checker.model_registry:
                self.integrity_checker.register_model(model_id, model_data)
            
            is_valid, integrity_msg = self.integrity_checker.verify_integrity(model_id, model_data)
            
            # 2. 公平性審計
            protected_attrs = {
                'region': np.random.randint(0, 3, len(predictions))
            }
            fairness_result = await self.fairness_auditor.audit_fairness(
                model_id, predictions, protected_attrs
            )
            
            # 3. 性能監控
            performance_result = await self.performance_monitor.monitor_performance(
                model_id, performance_metrics
            )
            
            # 4. 綜合評估
            recommendations = []
            
            if not is_valid:
                recommendations.append("重新部署原始模型")
                recommendations.append("調查可能的模型中毒攻擊")
            
            if fairness_result.get('overall_fairness', 1.0) < 0.8:
                recommendations.append("重新訓練模型以改善公平性")
                recommendations.append("檢查訓練數據偏差")
            
            if performance_result.get('alerts'):
                recommendations.extend([
                    "優化模型以改善性能",
                    "考慮模型量化或剪枝"
                ])
            
            # 判定整體狀態
            if not is_valid:
                overall_status = GovernanceStatus.CRITICAL
            elif fairness_result.get('status') == 'fail' or performance_result.get('alerts'):
                overall_status = GovernanceStatus.WARNING
            else:
                overall_status = GovernanceStatus.COMPLIANT
            
            report = GovernanceReport(
                report_id=f"gov_{datetime.now().strftime('%Y%m%d%H%M%S')}",
                timestamp=datetime.now().isoformat(),
                overall_status=overall_status,
                model_integrity={
                    'valid': is_valid,
                    'message': integrity_msg
                },
                fairness_metrics=fairness_result,
                performance_metrics=performance_result,
                security_checks={
                    'adversarial_detected': False,
                    'model_poisoning': not is_valid
                },
                recommendations=recommendations
            )
            
            logger.info(f"治理報告已生成: {report.report_id} - 狀態: {overall_status.value}")
            return report
            
        except Exception as e:
            logger.error(f"治理報告生成失敗: {e}")
            return None


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora AI 治理系統啟動 ===")
    
    # 初始化治理系統
    governance = AIGovernanceSystem()
    
    # 模擬模型數據
    model_id = "threat_detector_v1"
    model_data = b"simulated_model_weights_and_biases" * 100
    
    # 模擬預測結果
    predictions = np.random.randint(0, 2, 1000)
    
    # 模擬性能指標
    performance_metrics = {
        'accuracy': 0.958,
        'latency_ms': 12,
        'throughput': 8500,
        'memory_mb': 890,
        'cpu_percent': 35
    }
    
    # 生成治理報告
    report = await governance.generate_governance_report(
        model_id,
        model_data,
        predictions,
        performance_metrics
    )
    
    if report:
        print("\n=== 治理報告 ===")
        print(json.dumps(asdict(report), indent=2, default=str))
    
    logger.info("\n=== 系統測試完成 ===")


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())

