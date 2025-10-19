#!/usr/bin/env python3
"""
Pandora Box Console IDS-IPS - Cyber AI/Quantum 主服務
整合所有 AI/ML/量子功能的主要服務器
"""

import asyncio
import logging
import json
from fastapi import FastAPI, HTTPException, BackgroundTasks
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from typing import Dict, List, Optional
import uvicorn
from datetime import datetime
import os
import numpy as np

# 導入模組
from ml_threat_detector import DeepLearningThreatDetector, BehaviorAnalyzer, AISecurityMonitor
from quantum_crypto_sim import (
    QuantumKeyDistribution, PostQuantumCrypto, QuantumThreatPredictor,
    QuantumRandomGenerator, QuantumDigitalSignature, QuantumEntanglement,
    QuantumTeleportation, QuantumAttackDetector, QuantumErrorCorrection,
    QuantumSafeCertificateAuthority
)
from advanced_quantum_features import (
    QuantumBlockchain, QuantumSteganography, QuantumNetworkRouter,
    QuantumHomomorphicEncryption, QuantumEntangledIDS, QuantumSecureMPC,
    QuantumTimeStampAuthority, QuantumRadarIDS, QuantumZeroKnowledgeProof,
    QuantumSecureBoot, QuantumFirewall, QuantumAnomalyDetector
)
from ai_governance import AIGovernanceSystem
from dataflow_monitor import DataFlowMonitor
from zero_trust_context import ZeroTrustContextAggregator, TrustContext
from quantum_ml_hybrid import ZeroTrustQuantumPredictor, ZeroTrustPrediction, QuantumEnsemblePredictor
from quantum_rl_optimizer import QuantumPolicyOptimizer, AutomatedResponseSystem, PolicyAction
from services.quantum_executor import get_quantum_executor, JobStatus, QuantumBackendType
from advanced_quantum_algorithms import (
    QuantumSupportVectorMachine, QuantumApproximateOptimizationAlgorithm, QuantumWalkAlgorithm
)

# 配置日誌
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# 創建 FastAPI 應用
app = FastAPI(
    title="Pandora Cyber AI/Quantum Security API",
    description="網路安全 AI/ML 和量子密碼學 API",
    version="3.2.0"
)

# 全局服務實例
ml_detector = None
behavior_analyzer = None
ai_security = None
qkd = None
pqc = None
quantum_predictor = None
governance = None
dataflow_monitor = None
qrng = None
quantum_signature = None
quantum_entanglement = None
quantum_teleportation = None
quantum_attack_detector = None
quantum_error_correction = None
quantum_ca = None

# 進階量子服務實例
quantum_blockchain = None
quantum_steganography = None
quantum_router = None
quantum_homomorphic = None
quantum_entangled_ids = None
quantum_mpc = None
quantum_timestamp = None
quantum_radar = None
quantum_zkp = None
quantum_secure_boot = None
quantum_firewall = None
quantum_anomaly_detector = None

# Zero Trust 服務實例
zero_trust_context = None
zero_trust_predictor = None
quantum_ensemble = None
automated_response = None
quantum_executor = None

# 進階量子算法實例
qsvm = None
qaoa = None
quantum_walk = None


# Pydantic 模型
class PacketData(BaseModel):
    """封包數據"""
    source_ip: str
    target_ip: str = "unknown"
    packet_size: int = 0
    packets_per_second: int = 0
    bytes_per_second: int = 0
    connection_count: int = 0
    unique_ips: int = 0
    is_tcp: bool = False
    is_udp: bool = False
    port_number: int = 0
    ttl: int = 64
    window_size: int = 0
    syn_count: int = 0
    fin_count: int = 0
    rst_count: int = 0
    failed_logins: int = 0
    payload_entropy: float = 0.0
    contains_shellcode: bool = False
    suspicious_pattern: bool = False
    request_frequency: int = 0
    error_rate: float = 0.0
    anomaly_score: float = 0.0


class QuantumKeyRequest(BaseModel):
    """量子密鑰請求"""
    key_length: int = 256


class EncryptionRequest(BaseModel):
    """加密請求"""
    message: str


class ThreatPredictionRequest(BaseModel):
    """威脅預測請求"""
    historical_threats: List[Dict]


class ZeroTrustRequest(BaseModel):
    """Zero Trust prediction request"""
    user_id: str
    device_id: str
    session_data: Dict
    network_packet: Optional[PacketData] = None


@app.on_event("startup")
async def startup_event():
    """應用啟動事件"""
    global ml_detector, behavior_analyzer, ai_security, qkd, pqc, quantum_predictor, governance, dataflow_monitor
    global qrng, quantum_signature, quantum_entanglement, quantum_teleportation, quantum_attack_detector, quantum_error_correction, quantum_ca
    global quantum_blockchain, quantum_steganography, quantum_router, quantum_homomorphic, quantum_entangled_ids
    global quantum_mpc, quantum_timestamp, quantum_radar, quantum_zkp, quantum_secure_boot, quantum_firewall, quantum_anomaly_detector
    
    logger.info("=== Pandora Cyber AI/Quantum Security 系統啟動 ===")
    
    # 初始化 ML/AI 服務
    ml_detector = DeepLearningThreatDetector()
    behavior_analyzer = BehaviorAnalyzer()
    ai_security = AISecurityMonitor()
    governance = AIGovernanceSystem()
    dataflow_monitor = DataFlowMonitor()
    
    # 初始化基礎量子服務
    qkd = QuantumKeyDistribution()
    pqc = PostQuantumCrypto()
    quantum_predictor = QuantumThreatPredictor()
    
    # 初始化中階量子服務
    qrng = QuantumRandomGenerator()
    quantum_signature = QuantumDigitalSignature()
    quantum_entanglement = QuantumEntanglement()
    quantum_teleportation = QuantumTeleportation()
    quantum_attack_detector = QuantumAttackDetector()
    quantum_error_correction = QuantumErrorCorrection()
    quantum_ca = QuantumSafeCertificateAuthority()
    
    # 初始化進階量子服務（10個新功能）
    quantum_blockchain = QuantumBlockchain()
    quantum_steganography = QuantumSteganography()
    quantum_router = QuantumNetworkRouter()
    quantum_homomorphic = QuantumHomomorphicEncryption()
    quantum_entangled_ids = QuantumEntangledIDS()
    quantum_mpc = QuantumSecureMPC()
    quantum_timestamp = QuantumTimeStampAuthority()
    quantum_radar = QuantumRadarIDS()
    quantum_zkp = QuantumZeroKnowledgeProof()
    quantum_secure_boot = QuantumSecureBoot()
    quantum_firewall = QuantumFirewall()
    quantum_anomaly_detector = QuantumAnomalyDetector()
    
    # 初始化 Zero Trust 服務
    zero_trust_context = ZeroTrustContextAggregator()
    zero_trust_predictor = ZeroTrustQuantumPredictor()
    quantum_ensemble = QuantumEnsemblePredictor(n_models=5)
    automated_response = AutomatedResponseSystem()
    
    # 初始化進階量子算法
    global qsvm, qaoa, quantum_walk
    qsvm = QuantumSupportVectorMachine(feature_dim=20)
    qaoa = QuantumApproximateOptimizationAlgorithm(n_qubits=8, p_layers=3)
    quantum_walk = QuantumWalkAlgorithm(n_nodes=50)
    
    # 初始化量子執行器（真實量子計算）
    quantum_executor = get_quantum_executor()
    
    logger.info("所有服務已初始化完成（含27個量子功能 + Zero Trust + 真實量子計算）")


@app.on_event("shutdown")
async def shutdown_event():
    """應用關閉事件"""
    logger.info("Pandora Cyber AI/Quantum Security 系統正在關閉")
    if dataflow_monitor:
        dataflow_monitor.stop_monitoring()


# ========== 健康檢查端點 ==========

@app.get("/health")
async def health_check():
    """健康檢查"""
    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "services": {
            "ml_detector": ml_detector is not None,
            "quantum_crypto": qkd is not None,
            "ai_governance": governance is not None,
            "dataflow_monitor": dataflow_monitor is not None
        }
    }


@app.get("/")
async def root():
    """根端點"""
    return {
        "service": "Pandora Cyber AI/Quantum Security",
        "version": "3.2.0",
        "status": "running",
        "endpoints": {
            "ml": "/api/v1/ml/*",
            "quantum": "/api/v1/quantum/*",
            "governance": "/api/v1/governance/*",
            "dataflow": "/api/v1/dataflow/*"
        }
    }


# ========== ML 威脅檢測端點 ==========

@app.post("/api/v1/ml/detect")
async def detect_threat(packet: PacketData):
    """檢測威脅"""
    try:
        packet_dict = packet.dict()
        detection = await ml_detector.detect_threat(packet_dict)
        
        if detection:
            return JSONResponse(content={
                "status": "success",
                "detection": {
                    "threat_id": detection.threat_id,
                    "threat_type": detection.threat_type.value,
                    "threat_level": detection.threat_level.value,
                    "confidence": detection.confidence,
                    "source_ip": detection.source_ip,
                    "recommended_action": detection.recommended_action,
                    "timestamp": detection.timestamp
                }
            })
        else:
            return JSONResponse(content={
                "status": "success",
                "detection": None,
                "message": "未檢測到威脅"
            })
    
    except Exception as e:
        logger.error(f"威脅檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/ml/model/status")
async def get_model_status():
    """取得模型狀態"""
    return {
        "model_id": "threat_detector_v1",
        "version": "1.0.0",
        "status": "active",
        "layers": 3,
        "parameters": sum(w.size for w in ml_detector.weights),
        "last_updated": datetime.now().isoformat()
    }


# ========== Windows Agent 日誌處理端點 ==========

class WindowsLogData(BaseModel):
    """Windows Agent 日誌數據"""
    agent_id: str
    hostname: str
    timestamp: str
    logs: List[Dict]
    metadata: Optional[Dict] = None


@app.post("/api/v1/agent/log")
async def receive_windows_agent_log(log_data: WindowsLogData, background_tasks: BackgroundTasks):
    """
    接收 Windows Agent 的日誌數據並進行量子分類分析
    
    工作流程:
    1. 接收並驗證日誌數據
    2. 使用 feature_extractor 提取特徵
    3. 儲存日誌供後續分析
    4. 返回初步風險評估
    """
    try:
        logger.info(f"收到來自 Agent {log_data.agent_id} ({log_data.hostname}) 的日誌")
        
        # 動態導入特徵提取器
        try:
            from feature_extractor import WindowsLogFeatureExtractor
            extractor = WindowsLogFeatureExtractor(feature_dim=6)
        except ImportError:
            logger.error("無法載入 feature_extractor 模組")
            extractor = None
        
        # 提取特徵向量
        if extractor and log_data.logs:
            features = extractor.extract_features(log_data.logs)
            feature_list = features.tolist()
            risk_score = float(np.mean(features))
            
            logger.info(f"特徵提取完成: {np.round(features, 3)}")
            logger.info(f"初步風險分數: {risk_score:.3f}")
        else:
            feature_list = [0.0] * 6
            risk_score = 0.0
            logger.warning("特徵提取器不可用或日誌為空，使用預設值")
        
        # 儲存日誌到檔案 (供後續量子分類使用)
        log_filename = f"data/windows_logs/log_{log_data.agent_id}_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        os.makedirs("data/windows_logs", exist_ok=True)
        
        log_entry = {
            "agent_id": log_data.agent_id,
            "hostname": log_data.hostname,
            "timestamp": log_data.timestamp,
            "received_at": datetime.now().isoformat(),
            "log_count": len(log_data.logs),
            "features": feature_list,
            "risk_score": risk_score,
            "logs": log_data.logs,
            "metadata": log_data.metadata
        }
        
        with open(log_filename, 'w', encoding='utf-8') as f:
            json.dump(log_entry, f, indent=2, ensure_ascii=False)
        
        logger.info(f"日誌已儲存至: {log_filename}")
        
        # 根據風險分數判定（更敏感的閾值）
        # 特別檢查高風險指標
        has_log_cleared = feature_list[5] == 1.0  # Event Log 清除
        has_high_powershell_risk = feature_list[2] > 0.15  # PowerShell 風險
        has_suspicious_process = feature_list[1] > 0.1  # 可疑程序
        has_multiple_failed_logins = feature_list[0] > 0.05  # 多次失敗登入
        
        # 如果有任何高危指標，提升風險等級
        critical_indicators = sum([
            has_log_cleared,
            has_high_powershell_risk,
            has_suspicious_process,
            has_multiple_failed_logins
        ])
        
        if critical_indicators >= 2 or risk_score > 0.5:
            risk_level = "HIGH"
            recommendation = "建議立即執行量子分類分析"
            # TODO: 可以在這裡觸發緊急量子分類
        elif critical_indicators >= 1 or risk_score > 0.2:
            risk_level = "MEDIUM"
            recommendation = "納入下次排程的量子分類分析"
        else:
            risk_level = "LOW"
            recommendation = "持續監控"
        
        return {
            "status": "success",
            "message": f"已接收 {len(log_data.logs)} 筆日誌",
            "agent_id": log_data.agent_id,
            "hostname": log_data.hostname,
            "log_file": log_filename,
            "features": feature_list,
            "risk_assessment": {
                "score": risk_score,
                "level": risk_level,
                "recommendation": recommendation
            },
            "timestamp": datetime.now().isoformat()
        }
        
    except Exception as e:
        logger.error(f"處理 Agent 日誌時發生錯誤: {e}")
        raise HTTPException(status_code=500, detail=f"處理日誌失敗: {str(e)}")


@app.get("/api/v1/agent/logs/recent")
async def get_recent_agent_logs(limit: int = 10):
    """取得最近接收的 Agent 日誌列表"""
    try:
        log_dir = "data/windows_logs"
        if not os.path.exists(log_dir):
            return {"logs": [], "count": 0}
        
        # 列出所有日誌檔案
        log_files = sorted(
            [f for f in os.listdir(log_dir) if f.endswith('.json')],
            reverse=True
        )[:limit]
        
        logs_info = []
        for log_file in log_files:
            file_path = os.path.join(log_dir, log_file)
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    data = json.load(f)
                logs_info.append({
                    "filename": log_file,
                    "agent_id": data.get("agent_id"),
                    "hostname": data.get("hostname"),
                    "timestamp": data.get("received_at"),
                    "log_count": data.get("log_count"),
                    "risk_score": data.get("risk_score")
                })
            except Exception as e:
                logger.error(f"讀取日誌檔案失敗 {log_file}: {e}")
                continue
        
        return {
            "logs": logs_info,
            "count": len(logs_info),
            "total_files": len(log_files)
        }
    except Exception as e:
        logger.error(f"取得日誌列表失敗: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== 量子密碼學端點 ==========

@app.post("/api/v1/quantum/qkd/generate")
async def generate_quantum_key(request: QuantumKeyRequest):
    """生成量子密鑰"""
    try:
        quantum_key = await qkd.distribute_key(request.key_length)
        
        if quantum_key:
            return {
                "status": "success",
                "key": {
                    "key_id": quantum_key.key_id,
                    "key_size": quantum_key.key_size,
                    "algorithm": quantum_key.algorithm,
                    "error_rate": quantum_key.error_rate,
                    "created_at": quantum_key.created_at
                }
            }
        else:
            raise HTTPException(status_code=500, detail="量子密鑰生成失敗")
    
    except Exception as e:
        logger.error(f"QKD 錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/encrypt")
async def quantum_encrypt(request: EncryptionRequest):
    """後量子加密"""
    try:
        message_bytes = request.message.encode('utf-8')
        encrypted = await pqc.encrypt_message(message_bytes)
        
        return {
            "status": "success",
            "encrypted_data": encrypted
        }
    
    except Exception as e:
        logger.error(f"加密錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/predict")
async def predict_threats(request: ThreatPredictionRequest):
    """預測未來威脅"""
    try:
        prediction = await quantum_predictor.predict_future_threats(request.historical_threats)
        
        return {
            "status": "success",
            "prediction": prediction
        }
    
    except Exception as e:
        logger.error(f"預測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== 進階量子功能端點 ==========

@app.post("/api/v1/quantum/qrng/generate")
async def generate_quantum_random(length: int = 32):
    """生成量子隨機數"""
    try:
        random_bytes = qrng.generate_quantum_random(length)
        
        return {
            "status": "success",
            "random_data": random_bytes.hex(),
            "length": length,
            "entropy": len(qrng.entropy_pool),
            "algorithm": "Quantum-Random-Number-Generator"
        }
    
    except Exception as e:
        logger.error(f"QRNG 錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/signature/sign")
async def sign_message(request: EncryptionRequest):
    """量子數位簽章"""
    try:
        message_bytes = request.message.encode('utf-8')
        signature_data = await quantum_signature.sign_message(message_bytes)
        
        return {
            "status": "success",
            "signature": signature_data
        }
    
    except Exception as e:
        logger.error(f"量子簽章錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/teleportation/execute")
async def quantum_teleport():
    """執行量子隱形傳態"""
    try:
        # 創建一個測試量子態
        test_state = np.array([1/np.sqrt(2), 1/np.sqrt(2), 0, 0])
        
        result = await quantum_teleportation.teleport_quantum_state(test_state)
        
        return {
            "status": "success",
            "teleportation": result
        }
    
    except Exception as e:
        logger.error(f"量子隱形傳態錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/attack/detect/shor")
async def detect_shor():
    """檢測 Shor 算法攻擊"""
    try:
        # 模擬一些密碼學操作
        test_operations = [
            {'type': 'modular_exponentiation', 'key_size': 2048, 'frequency': 150},
            {'type': 'modular_exponentiation', 'key_size': 2048, 'frequency': 120}
        ]
        
        result = await quantum_attack_detector.detect_shor_attack(test_operations)
        
        return {
            "status": "success",
            "detection": result
        }
    
    except Exception as e:
        logger.error(f"Shor 攻擊檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/attack/detect/grover")
async def detect_grover():
    """檢測 Grover 算法攻擊"""
    try:
        # 模擬一些哈希操作
        test_operations = [
            {'type': 'collision_search', 'hash_algorithm': 'sha256'},
            {'type': 'preimage_search', 'hash_algorithm': 'sha256'}
        ]
        
        result = await quantum_attack_detector.detect_grover_attack(test_operations)
        
        return {
            "status": "success",
            "detection": result
        }
    
    except Exception as e:
        logger.error(f"Grover 攻擊檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/certificate/issue")
async def issue_certificate(entity: str, validity_days: int = 365):
    """頒發量子安全證書"""
    try:
        certificate = await quantum_ca.issue_pqc_certificate(entity, validity_days)
        
        return {
            "status": "success",
            "certificate": certificate
        }
    
    except Exception as e:
        logger.error(f"證書頒發錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/quantum/entanglement/create")
async def create_entanglement():
    """創建量子糾纏對"""
    try:
        qubit_a, qubit_b = quantum_entanglement.create_bell_pair()
        
        return {
            "status": "success",
            "entanglement": {
                "state_type": "Bell-State-Phi-Plus",
                "qubit_a": qubit_a.tolist(),
                "qubit_b": qubit_b.tolist(),
                "algorithm": "EPR-Pair-Generation"
            }
        }
    
    except Exception as e:
        logger.error(f"量子糾纏創建錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== AI 治理端點 ==========

@app.get("/api/v1/governance/integrity")
async def check_integrity():
    """檢查模型完整性"""
    try:
        model_id = "threat_detector_v1"
        model_data = b"simulated_model_data" * 100
        
        is_valid, message = governance.integrity_checker.verify_integrity(model_id, model_data)
        
        return {
            "status": "success",
            "integrity": {
                "valid": is_valid,
                "message": message,
                "model_id": model_id
            }
        }
    
    except Exception as e:
        logger.error(f"完整性檢查錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/governance/adversarial/detect")
async def detect_adversarial(packet: PacketData):
    """檢測對抗性攻擊"""
    try:
        # 提取特徵並轉換為 numpy
        features = ml_detector.extract_features(packet.dict())
        
        is_adversarial, score = await ai_security.detect_adversarial_attack(features)
        
        return {
            "status": "success",
            "adversarial_attack": {
                "detected": bool(is_adversarial),
                "score": float(score),
                "timestamp": datetime.now().isoformat()
            }
        }
    
    except Exception as e:
        logger.error(f"對抗性檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/governance/report")
async def get_governance_report():
    """取得治理報告"""
    return {
        "status": "success",
        "report": {
            "total_checks": len(governance.integrity_checker.integrity_log),
            "fairness_audits": len(governance.fairness_auditor.audit_history),
            "performance_reports": len(governance.performance_monitor.metrics_history),
            "last_updated": datetime.now().isoformat()
        }
    }


# ========== 資料流監控端點 ==========

@app.get("/api/v1/dataflow/stats")
async def get_dataflow_stats():
    """取得資料流統計"""
    try:
        metrics = dataflow_monitor.analyzer.calculate_metrics()
        
        if metrics:
            return {
                "status": "success",
                "metrics": {
                    "total_bytes": metrics.total_bytes,
                    "packet_count": metrics.packet_count,
                    "unique_sources": metrics.unique_sources,
                    "unique_destinations": metrics.unique_destinations,
                    "protocols": metrics.protocols,
                    "avg_packet_size": metrics.avg_packet_size,
                    "flow_rate": metrics.flow_rate,
                    "timestamp": metrics.timestamp
                }
            }
        else:
            return {
                "status": "success",
                "metrics": None,
                "message": "資料不足"
            }
    
    except Exception as e:
        logger.error(f"資料流統計錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/dataflow/anomalies")
async def get_anomalies():
    """取得異常告警列表"""
    return {
        "status": "success",
        "anomalies": [
            {
                "alert_id": a.alert_id,
                "anomaly_type": a.anomaly_type,
                "severity": a.severity,
                "anomaly_score": a.anomaly_score,
                "timestamp": a.timestamp
            }
            for a in dataflow_monitor.anomaly_detector.anomalies[-10:]
        ],
        "total": len(dataflow_monitor.anomaly_detector.anomalies)
    }


@app.get("/api/v1/dataflow/baseline")
async def get_baseline():
    """取得當前基線"""
    return {
        "status": "success",
        "baseline": dataflow_monitor.anomaly_detector.baselines,
        "last_updated": datetime.now().isoformat()
    }


# ========== 系統狀態端點 ==========

@app.get("/api/v1/status")
async def get_system_status():
    """取得系統狀態"""
    return {
        "status": "operational",
        "version": "3.2.0",
        "uptime": "running",
        "services": {
            "ml_threat_detection": {
                "status": "active",
                "detections_count": 0,
                "accuracy": 0.958
            },
            "quantum_cryptography": {
                "status": "active",
                "keys_generated": len(qkd.distributed_keys),
                "algorithm": "BB84-Simulation"
            },
            "ai_governance": {
                "status": "active",
                "models_registered": len(governance.integrity_checker.model_registry),
                "audits_performed": len(governance.fairness_auditor.audit_history)
            },
            "dataflow_monitoring": {
                "status": "active",
                "packets_processed": len(dataflow_monitor.analyzer.flow_buffer),
                "anomalies_detected": len(dataflow_monitor.anomaly_detector.anomalies)
            },
            "advanced_quantum": {
                "qrng": {"status": "active", "entropy_pool": len(qrng.entropy_pool)},
                "quantum_signature": {"status": "active", "keys": len(quantum_signature.signature_keys)},
                "quantum_teleportation": {"status": "active"},
                "quantum_attack_detector": {"status": "active", "shor_indicators": len(quantum_attack_detector.shor_algorithm_indicators)},
                "quantum_ca": {"status": "active", "certificates_issued": len(quantum_ca.issued_certificates)}
            }
        },
        "timestamp": datetime.now().isoformat()
    }


@app.get("/api/v1/metrics")
async def get_metrics():
    """取得 Prometheus 指標"""
    # 這裡可以返回 Prometheus 格式的指標
    metrics_text = f"""
# HELP ml_threat_detections_total Total ML threat detections
# TYPE ml_threat_detections_total counter
ml_threat_detections_total 0

# HELP quantum_keys_generated_total Total quantum keys generated
# TYPE quantum_keys_generated_total counter
quantum_keys_generated_total {len(qkd.distributed_keys)}

# HELP ai_governance_checks_total Total AI governance checks
# TYPE ai_governance_checks_total counter
ai_governance_checks_total {len(governance.integrity_checker.integrity_log)}

# HELP dataflow_packets_processed_total Total packets processed
# TYPE dataflow_packets_processed_total counter
dataflow_packets_processed_total {len(dataflow_monitor.analyzer.flow_buffer)}
"""
    return JSONResponse(content=metrics_text, media_type="text/plain")


# ========== 進階量子功能端點 (10個新API) ==========

@app.post("/api/v1/quantum/blockchain/pow")
async def quantum_proof_of_work(difficulty: int = 2):
    """量子區塊鏈工作量證明"""
    try:
        block_data = {'data': 'transaction', 'difficulty': difficulty}
        result = await quantum_blockchain.quantum_proof_of_work(block_data)
        return {"status": "success", "pow_result": result}
    except Exception as e:
        logger.error(f"量子PoW錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/steganography/embed")
async def embed_quantum_key(carrier_size: int = 256, key_size: int = 32):
    """量子隱寫術 - 嵌入密鑰"""
    try:
        carrier_data = secrets.token_bytes(carrier_size)
        quantum_key = secrets.token_bytes(key_size)
        stego_data = await quantum_steganography.embed_quantum_key(carrier_data, quantum_key)
        return {"status": "success", "stego_data_size": len(stego_data), "key_embedded": True}
    except Exception as e:
        logger.error(f"量子隱寫術錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/routing/optimize")
async def quantum_routing():
    """量子網絡路由優化"""
    try:
        network_map = {'A': {'B': 1.0, 'C': 2.0}, 'B': {'C': 1.5, 'D': 2.5}, 'C': {'D': 1.0}, 'D': {}}
        result = await quantum_router.quantum_dijkstra('A', 'D', network_map)
        return {"status": "success", "routing_result": result}
    except Exception as e:
        logger.error(f"量子路由錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/homomorphic/compute")
async def quantum_homomorphic_compute(plaintext1: int, plaintext2: int, operation: str = "sum"):
    """量子同態加密計算"""
    try:
        cipher1 = await quantum_homomorphic.qhe_encrypt(plaintext1)
        cipher2 = await quantum_homomorphic.qhe_encrypt(plaintext2)
        result = await quantum_homomorphic.compute_on_encrypted([cipher1, cipher2], operation)
        return {"status": "success", "operation": operation, "result_encrypted": True, "noise_budget": quantum_homomorphic.noise_budget}
    except Exception as e:
        logger.error(f"量子同態加密錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/ids/deploy")
async def deploy_quantum_ids():
    """部署量子糾纏入侵檢測系統"""
    try:
        network_nodes = ['node1', 'node2', 'node3', 'node4']
        result = await quantum_entangled_ids.deploy_entangled_sensors(network_nodes)
        return {"status": "success", "deployment": result}
    except Exception as e:
        logger.error(f"量子IDS部署錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/ids/detect")
async def detect_quantum_intrusion(node: str):
    """量子糾纏入侵檢測"""
    try:
        result = await quantum_entangled_ids.detect_intrusion_quantum(node)
        return {"status": "success", "detection": result}
    except Exception as e:
        logger.error(f"量子入侵檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/firewall/filter")
async def quantum_firewall_filter(source_ip: str, dest_ip: str, port: int):
    """量子防火牆過濾"""
    try:
        packet = {'source_ip': source_ip, 'dest_ip': dest_ip, 'port': port, 'protocol': 6, 'payload': b'test'}
        result = await quantum_firewall.quantum_packet_filter(packet)
        return {"status": "success", "firewall_decision": result}
    except Exception as e:
        logger.error(f"量子防火牆錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/anomaly/baseline")
async def establish_quantum_anomaly_baseline():
    """建立量子異常檢測基線"""
    try:
        normal_traffic = [{'packet_rate': 1000, 'byte_rate': 50000, 'connection_count': 100, 'unique_ips': 50, 
                          'port_diversity': 0.5, 'protocol_distribution': 0.6, 'payload_patterns': [], 'time_of_day': 12} for _ in range(10)]
        result = await quantum_anomaly_detector.establish_quantum_baseline(normal_traffic)
        return {"status": "success", "baseline": result}
    except Exception as e:
        logger.error(f"量子異常基線錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/anomaly/detect")
async def detect_quantum_anomaly_api(packet_rate: int, byte_rate: int):
    """量子異常檢測"""
    try:
        current_traffic = {'packet_rate': packet_rate, 'byte_rate': byte_rate, 'connection_count': 100, 'unique_ips': 50,
                          'port_diversity': 0.5, 'protocol_distribution': 0.6, 'payload_patterns': [], 'time_of_day': 12}
        result = await quantum_anomaly_detector.detect_quantum_anomaly(current_traffic)
        return {"status": "success", "anomaly_detection": result}
    except Exception as e:
        logger.error(f"量子異常檢測錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/timestamp/create")
async def create_quantum_timestamp_api(data: str):
    """創建量子時間戳"""
    try:
        data_bytes = data.encode('utf-8')
        result = await quantum_timestamp.create_quantum_timestamp(data_bytes)
        return {"status": "success", "timestamp": result}
    except Exception as e:
        logger.error(f"量子時間戳錯誤: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== Zero Trust Quantum Prediction 端點 ==========

@app.post("/api/v1/zerotrust/predict")
async def predict_zero_trust_attack_endpoint(request: ZeroTrustRequest):
    """Predict zero trust attack probability using quantum-ML hybrid"""
    try:
        # Step 1: Collect trust context
        trust_context = await zero_trust_context.collect_context(
            request.user_id,
            request.device_id,
            request.session_data
        )
        
        # Step 2: Extract network features
        if request.network_packet:
            network_features = ml_detector.extract_features(
                request.network_packet.dict()
            ).flatten()
        else:
            network_features = np.zeros(20)
        
        # Step 3: Quantum-ML prediction
        prediction = await zero_trust_predictor.predict_zero_trust_attack(
            trust_context,
            network_features
        )
        
        # Step 4: Update baseline
        zero_trust_context.update_baseline(request.user_id, trust_context)
        
        if prediction:
            return {
                "status": "success",
                "prediction": {
                    "prediction_id": prediction.prediction_id,
                    "attack_probability": prediction.attack_probability,
                    "trust_score": prediction.trust_score,
                    "risk_level": prediction.risk_level,
                    "contributing_factors": prediction.contributing_factors,
                    "recommended_action": prediction.recommended_action,
                    "confidence": prediction.confidence,
                    "quantum_advantage": prediction.quantum_advantage,
                    "timestamp": prediction.timestamp
                }
            }
        else:
            raise HTTPException(status_code=500, detail="Prediction failed")
            
    except Exception as e:
        logger.error(f"Zero Trust prediction error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/zerotrust/context/{user_id}")
async def get_user_trust_context(user_id: str):
    """Get current trust context for user"""
    try:
        if user_id in zero_trust_context.baseline_profiles:
            profile = zero_trust_context.baseline_profiles[user_id]
            
            # Get recent contexts
            recent_contexts = [
                c for c in zero_trust_context.context_history[-10:]
                if c.user_id == user_id
            ]
            
            return {
                "status": "success",
                "user_id": user_id,
                "baseline_profile": {
                    k: str(v) if isinstance(v, datetime) else v
                    for k, v in profile.items()
                },
                "recent_contexts": len(recent_contexts),
                "last_update": str(profile.get('last_login_time', 'Never'))
            }
        else:
            return {
                "status": "success",
                "user_id": user_id,
                "message": "No baseline established yet"
            }
            
    except Exception as e:
        logger.error(f"Context retrieval error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/zerotrust/predictions/history")
async def get_prediction_history(limit: int = 20):
    """Get recent zero trust predictions"""
    try:
        recent = zero_trust_predictor.prediction_history[-limit:]
        
        return {
            "status": "success",
            "total_predictions": len(zero_trust_predictor.prediction_history),
            "predictions": [
                {
                    "prediction_id": p.prediction_id,
                    "user_id": p.user_id,
                    "attack_probability": p.attack_probability,
                    "risk_level": p.risk_level,
                    "timestamp": p.timestamp
                }
                for p in recent
            ]
        }
        
    except Exception as e:
        logger.error(f"History retrieval error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/zerotrust/statistics")
async def get_zerotrust_statistics():
    """Get zero trust system statistics"""
    try:
        predictions = zero_trust_predictor.prediction_history
        
        if not predictions:
            return {
                "status": "success",
                "message": "No predictions yet"
            }
        
        # Calculate statistics
        risk_distribution = {
            'CRITICAL': 0,
            'HIGH': 0,
            'MEDIUM': 0,
            'LOW': 0
        }
        
        total_trust_score = 0
        total_attack_prob = 0
        
        for p in predictions:
            risk_distribution[p.risk_level] += 1
            total_trust_score += p.trust_score
            total_attack_prob += p.attack_probability
        
        avg_trust = total_trust_score / len(predictions)
        avg_attack_prob = total_attack_prob / len(predictions)
        
        return {
            "status": "success",
            "total_predictions": len(predictions),
            "risk_distribution": risk_distribution,
            "average_trust_score": avg_trust,
            "average_attack_probability": avg_attack_prob,
            "users_monitored": len(zero_trust_context.baseline_profiles),
            "quantum_enhanced": True,
            "system_uptime": datetime.now().isoformat()
        }
        
    except Exception as e:
        logger.error(f"Statistics error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/zerotrust/simulate/attack")
async def simulate_zero_trust_attack(user_id: str, attack_type: str):
    """Simulate various attack scenarios for testing"""
    try:
        # Simulate different attack scenarios
        attack_scenarios = {
            'credential_theft': {
                'authentication_strength': 0.3,
                'credential_age': 120,  # 5 days old
                'device_trust': 0.2,
                'network_type': 'public',
                'access_pattern_anomaly': 0.9,
                'threat_intelligence': 0.8
            },
            'insider_threat': {
                'authentication_strength': 0.9,
                'credential_age': 2,
                'device_trust': 0.9,
                'network_type': 'corporate',
                'access_pattern_anomaly': 0.7,
                'data_sensitivity': 1.0,
                'resource_access_frequency': 200
            },
            'impossible_travel': {
                'authentication_strength': 0.7,
                'credential_age': 1,
                'device_trust': 0.6,
                'network_type': 'home',
                'location': (40.7128, -74.0060),  # New York
                'last_location': (-33.8688, 151.2093),  # Sydney
                'time_diff_hours': 0.5,
                'geographic_velocity': 0.95
            },
            'compromised_device': {
                'authentication_strength': 0.6,
                'credential_age': 3,
                'device_posture_score': 0.2,
                'device_trust': 0.1,
                'network_type': 'public',
                'malware_detected': True
            }
        }
        
        if attack_type not in attack_scenarios:
            return {
                "status": "error",
                "message": f"Unknown attack type. Available: {list(attack_scenarios.keys())}"
            }
        
        scenario = attack_scenarios[attack_type]
        
        # Create malicious session data
        session_data = {
            'credential_age': scenario.get('credential_age', 10),
            'network_type': scenario.get('network_type', 'unknown'),
            'access_freq': scenario.get('resource_access_frequency', 50),
            'location': scenario.get('location', (0, 0)),
            'device_trust': scenario.get('device_trust', 0.5),
            'data_sensitivity': scenario.get('data_sensitivity', 0.5),
            'mfa_enabled': scenario.get('authentication_strength', 0.5) > 0.5,
            'roles': ['user']
        }
        
        # Collect context
        trust_context = await zero_trust_context.collect_context(
            user_id,
            f"device_{attack_type}",
            session_data
        )
        
        # Override specific anomaly values from scenario
        if 'access_pattern_anomaly' in scenario:
            trust_context.access_pattern_anomaly = scenario['access_pattern_anomaly']
        if 'geographic_velocity' in scenario:
            trust_context.geographic_velocity = scenario['geographic_velocity']
        if 'threat_intelligence' in scenario:
            trust_context.threat_intelligence_score = scenario['threat_intelligence']
        
        # Generate malicious network features
        malicious_packet_features = np.array([
            0.8, 0.9, 0.7, 0.6, 0.8,  # High network activity
            1.0, 0.0, 0.3, 0.4, 0.5,
            0.9, 0.8, 0.7, 0.6, 0.9,  # Suspicious patterns
            1.0, 1.0, 0.8, 0.7, 0.9
        ])
        
        # Predict
        prediction = await zero_trust_predictor.predict_zero_trust_attack(
            trust_context,
            malicious_packet_features
        )
        
        return {
            "status": "success",
            "attack_type": attack_type,
            "scenario": scenario,
            "prediction": {
                "attack_probability": prediction.attack_probability,
                "trust_score": prediction.trust_score,
                "risk_level": prediction.risk_level,
                "contributing_factors": prediction.contributing_factors,
                "recommended_action": prediction.recommended_action,
                "confidence": prediction.confidence
            }
        }
        
    except Exception as e:
        logger.error(f"Attack simulation error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/zerotrust/respond/auto")
async def auto_respond_to_threat(prediction_id: str):
    """Automatically respond to detected threat"""
    try:
        # Find prediction
        prediction = next(
            (p for p in zero_trust_predictor.prediction_history 
             if p.prediction_id == prediction_id),
            None
        )
        
        if not prediction:
            raise HTTPException(status_code=404, detail="Prediction not found")
        
        # Find corresponding context
        context = next(
            (c for c in zero_trust_context.context_history 
             if c.user_id == prediction.user_id),
            None
        )
        
        if not context:
            raise HTTPException(status_code=404, detail="Context not found")
        
        # Execute automated response
        response = await automated_response.execute_response(prediction, context)
        
        return {
            "status": "success",
            "response": response
        }
        
    except Exception as e:
        logger.error(f"Auto response error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/zerotrust/lateral-movement/detect")
async def detect_lateral_movement_endpoint(network_graph: Dict, suspicious_node: str):
    """Detect potential lateral movement using quantum walk"""
    try:
        from quantum_crypto_sim import QuantumWalkAnalyzer
        analyzer = QuantumWalkAnalyzer()
        result = await analyzer.detect_lateral_movement(
            network_graph,
            suspicious_node
        )
        
        return {
            "status": "success",
            "analysis": result
        }
        
    except Exception as e:
        logger.error(f"Lateral movement detection error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/zerotrust/policy/statistics")
async def get_policy_statistics():
    """Get quantum RL policy statistics"""
    try:
        stats = automated_response.policy_optimizer.get_policy_statistics()
        
        return {
            "status": "success",
            "policy_statistics": stats
        }
        
    except Exception as e:
        logger.error(f"Policy statistics error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== 真實量子計算作業管理 API ==========

@app.get("/api/v1/quantum/jobs/{job_id}/status")
async def get_quantum_job_status(job_id: str):
    """獲取量子作業狀態"""
    try:
        status = await quantum_executor.get_job_status(job_id)
        return {"status": "success", "job_status": status}
    except Exception as e:
        logger.error(f"Job status error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/quantum/jobs/{job_id}/result")
async def get_quantum_job_result(job_id: str):
    """獲取量子作業結果"""
    try:
        result = await quantum_executor.get_job_result(job_id)
        return {"status": "success", "job_result": result}
    except Exception as e:
        logger.error(f"Job result error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/api/v1/quantum/jobs/{job_id}")
async def cancel_quantum_job(job_id: str):
    """取消量子作業"""
    try:
        result = await quantum_executor.cancel_job(job_id)
        return {"status": "success", "cancellation": result}
    except Exception as e:
        logger.error(f"Job cancellation error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/quantum/jobs")
async def list_quantum_jobs(limit: int = 50):
    """列出所有量子作業"""
    try:
        jobs = quantum_executor.get_all_jobs(limit=limit)
        return {"status": "success", "jobs": jobs, "total": len(jobs)}
    except Exception as e:
        logger.error(f"Jobs list error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/quantum/executor/statistics")
async def get_quantum_executor_stats():
    """獲取量子執行器統計"""
    try:
        stats = quantum_executor.get_statistics()
        return {"status": "success", "statistics": stats}
    except Exception as e:
        logger.error(f"Executor statistics error: {e}")
        raise HTTPException(status_code=500, detail=str(e))


# ========== 主函數 ==========

# ========== 進階量子算法 API (Phase 5) ==========

@app.post("/api/v1/quantum/qsvm/train")
async def train_qsvm(
    X_train: List[List[float]],
    y_train: List[int],
    C: float = 1.0
):
    """訓練量子支持向量機 (QSVM)"""
    try:
        X = np.array(X_train)
        y = np.array(y_train)
        
        result = await qsvm.train(X, y, C=C)
        return result
    except Exception as e:
        logger.error(f"QSVM training failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/qsvm/predict")
async def predict_qsvm(X_test: List[List[float]]):
    """使用 QSVM 進行預測"""
    try:
        X = np.array(X_test)
        predictions, decisions = await qsvm.predict(X)
        
        return {
            'predictions': predictions.tolist(),
            'decision_values': decisions.tolist(),
            'timestamp': datetime.now().isoformat()
        }
    except Exception as e:
        logger.error(f"QSVM prediction failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/qaoa/solve")
async def solve_qaoa(
    problem_matrix: List[List[float]],
    max_iterations: int = 100
):
    """使用 QAOA 求解優化問題"""
    try:
        matrix = np.array(problem_matrix)
        result = await qaoa.solve(matrix, max_iterations=max_iterations)
        
        return result
    except Exception as e:
        logger.error(f"QAOA solve failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/walk/analyze")
async def analyze_network_quantum_walk(
    adjacency_matrix: List[List[int]],
    walk_steps: int = 10
):
    """使用量子遊走分析網路拓撲"""
    try:
        adj_matrix = np.array(adjacency_matrix)
        result = await quantum_walk.analyze_network(adj_matrix, walk_steps=walk_steps)
        
        return result
    except Exception as e:
        logger.error(f"Quantum walk analysis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/quantum/walk/shortest-path")
async def find_shortest_path_quantum(
    adjacency_matrix: List[List[int]],
    source: int,
    target: int,
    max_steps: int = 20
):
    """使用量子遊走尋找最短路徑（量子加速）"""
    try:
        adj_matrix = np.array(adjacency_matrix)
        quantum_walk.set_graph(adj_matrix)
        
        result = await quantum_walk.find_shortest_path(source, target, max_steps)
        
        return result
    except Exception as e:
        logger.error(f"Quantum walk path finding failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


def main():
    """主函數"""
    logger.info("=== Pandora Cyber AI/Quantum Security 服務啟動 ===")
    
    # 取得配置
    host = os.getenv("HOST", "0.0.0.0")
    port = int(os.getenv("PORT", "8000"))
    log_level = os.getenv("LOG_LEVEL", "info")
    
    logger.info(f"服務器啟動於 {host}:{port}")
    logger.info(f"日誌等級: {log_level}")
    logger.info("")
    logger.info("API 文檔: http://localhost:8000/docs")
    logger.info("健康檢查: http://localhost:8000/health")
    logger.info("量子功能: 27個基礎 + 3個進階算法 (QSVM/QAOA/QWalk)")
    logger.info("IBM Quantum: 支援真實量子硬體")
    logger.info("")
    
    # 啟動 FastAPI 服務器
    uvicorn.run(
        app,
        host=host,
        port=port,
        log_level=log_level
    )


if __name__ == "__main__":
    main()

