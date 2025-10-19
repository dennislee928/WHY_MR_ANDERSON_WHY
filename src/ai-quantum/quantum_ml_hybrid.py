#!/usr/bin/env python3
"""
Quantum-ML Hybrid Zero Trust Attack Predictor
Combines quantum computing advantages with classical ML
"""

import numpy as np
from typing import Dict, List, Tuple, Optional
from dataclasses import dataclass, asdict
import logging
from datetime import datetime
import secrets

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@dataclass
class ZeroTrustPrediction:
    """Zero Trust attack prediction result"""
    user_id: str
    prediction_id: str
    attack_probability: float
    trust_score: float
    risk_level: str  # LOW, MEDIUM, HIGH, CRITICAL
    contributing_factors: Dict[str, float]
    recommended_action: str
    confidence: float
    timestamp: str
    quantum_advantage: bool


class QuantumNeuralNetwork:
    """Hybrid Quantum-Classical Neural Network (支援真實量子計算)"""
    
    def __init__(self, input_dim: int = 20, quantum_dim: int = 10, use_real_quantum: bool = False):
        """初始化 Quantum Neural Network"""
        self.input_dim = input_dim
        self.quantum_dim = quantum_dim
        self.use_real_quantum = use_real_quantum
        
        # Classical layers
        self.classical_weights1 = np.random.randn(input_dim, 32) * 0.1
        self.classical_bias1 = np.zeros(32)
        
        # Quantum layer weights
        self.quantum_weights = np.random.randn(32, quantum_dim) * 0.1
        
        # Output layers
        self.output_weights = np.random.randn(quantum_dim, 1) * 0.1
        self.output_bias = np.zeros(1)
        
        # 真實量子計算支援
        self.real_quantum_classifier = None
        if use_real_quantum:
            try:
                from poc_quantum_classifier import RealQuantumNeuralNetwork
                # 使用較小的量子電路（4 qubits）
                self.real_quantum_classifier = RealQuantumNeuralNetwork(
                    num_features=4,
                    num_qubits=4,
                    feature_map_reps=2,
                    ansatz_reps=3,
                    backend_type="simulator"
                )
                logger.info("Real Quantum layer initialized (Qiskit)")
            except Exception as e:
                logger.warning(f"Failed to initialize real quantum layer: {e}")
                logger.info("Falling back to NumPy simulation")
                self.use_real_quantum = False
        
        logger.info(f"Quantum Neural Network initialized (real_quantum={self.use_real_quantum})")
    
    def _classical_layer(self, x: np.ndarray) -> np.ndarray:
        """Classical neural network layer"""
        z = np.dot(x, self.classical_weights1) + self.classical_bias1
        return np.maximum(0, z)  # ReLU
    
    def _quantum_layer(self, x: np.ndarray) -> np.ndarray:
        """Quantum layer (可使用真實量子計算或模擬)"""
        
        if self.use_real_quantum and self.real_quantum_classifier:
            # 使用真實 Qiskit 量子計算
            try:
                # 壓縮到 4 維（真實量子電路的輸入）
                compressed = x[:4] if len(x) >= 4 else np.pad(x, (0, 4-len(x)))
                
                # 執行真實量子計算
                quantum_output, metadata = self.real_quantum_classifier.forward(compressed)
                
                # 擴展回所需維度
                output = np.full(self.quantum_dim, quantum_output)
                
                logger.debug("Real quantum computation executed")
                return output
                
            except Exception as e:
                logger.warning(f"Real quantum execution failed: {e}, falling back to simulation")
                # Fall through to simulation
        
        # NumPy 模擬（原始邏輯）
        # Project to quantum dimension
        quantum_input = np.dot(x, self.quantum_weights)
        
        # Simulate quantum variational circuit
        # Using parameterized rotation gates
        theta = quantum_input
        
        # RY rotations (simulate)
        quantum_state = np.cos(theta/2) + 1j * np.sin(theta/2)
        
        # Measurement (expectation values)
        measured = np.abs(quantum_state)**2
        
        return measured.real
    
    def forward(self, features: np.ndarray) -> Tuple[float, Dict]:
        """Forward pass through hybrid network"""
        # Classical preprocessing
        classical_output = self._classical_layer(features)
        
        # Quantum processing
        quantum_output = self._quantum_layer(classical_output)
        
        # Final classical layer
        prediction = np.dot(quantum_output, self.output_weights) + self.output_bias
        prediction = 1 / (1 + np.exp(-prediction))  # Sigmoid
        
        # Feature importance via quantum amplitude
        feature_importance = {
            f'feature_{i}': float(np.abs(quantum_output[i])) 
            for i in range(min(len(quantum_output), 10))
        }
        
        return float(prediction[0]), feature_importance


class ZeroTrustQuantumPredictor:
    """Main Zero Trust attack prediction system (支援真實量子計算)"""
    
    def __init__(self, use_real_quantum: bool = False, hybrid_fallback: bool = True):
        """
        初始化 Zero Trust Quantum Predictor
        
        Args:
            use_real_quantum: 是否使用真實量子硬體
            hybrid_fallback: 是否啟用混合後備（高風險才用量子）
        """
        self.use_real_quantum = use_real_quantum
        self.hybrid_fallback = hybrid_fallback
        
        # 初始化兩個模型
        self.classical_qnn = QuantumNeuralNetwork(input_dim=20, quantum_dim=10, use_real_quantum=False)
        
        if use_real_quantum:
            self.quantum_qnn = QuantumNeuralNetwork(input_dim=20, quantum_dim=10, use_real_quantum=True)
        else:
            self.quantum_qnn = self.classical_qnn  # 共用模型
        
        # Import here to avoid circular dependency
        from quantum_crypto_sim import QuantumFeatureMapper
        self.quantum_mapper = QuantumFeatureMapper(n_qubits=8)
        
        self.prediction_history = []
        
        # Risk thresholds
        self.thresholds = {
            'CRITICAL': 0.85,
            'HIGH': 0.70,
            'MEDIUM': 0.50,
            'LOW': 0.30
        }
        
        logger.info(
            f"Zero Trust Quantum Predictor initialized "
            f"(real_quantum={use_real_quantum}, hybrid_fallback={hybrid_fallback})"
        )
    
    async def predict_zero_trust_attack(
        self,
        trust_context,  # TrustContext type
        network_features: np.ndarray,
        force_quantum: bool = False
    ) -> ZeroTrustPrediction:
        """
        Predict zero trust attack probability
        
        Args:
            trust_context: 信任上下文
            network_features: 網絡特徵
            force_quantum: 強制使用量子計算（忽略後備邏輯）
        """
        
        try:
            # Step 1: Extract features from trust context
            context_features = self._extract_context_features(trust_context)
            
            # Step 2: Combine with network features
            combined_features = np.concatenate([context_features, network_features[:10]])
            
            # Step 3: Quantum feature enhancement
            quantum_features = self.quantum_mapper.quantum_feature_encoding(
                combined_features
            )
            
            # Step 4: 混合後備邏輯
            # 策略：先用快速古典模型，高風險才用真實量子
            use_quantum_for_this = force_quantum or (not self.hybrid_fallback)
            
            if self.hybrid_fallback and not force_quantum:
                # 快速古典預測
                classical_prob, classical_importance = self.classical_qnn.forward(quantum_features)
                
                # 檢查風險等級
                if classical_prob < self.thresholds['HIGH']:
                    # 低/中風險：直接返回古典結果
                    attack_prob = classical_prob
                    feature_importance = classical_importance
                    logger.debug(f"Using classical result (prob={classical_prob:.2%})")
                else:
                    # 高/極高風險：使用真實量子計算進行更精確分析
                    use_quantum_for_this = True
                    logger.info(f"High risk detected ({classical_prob:.2%}), using quantum computation")
            
            # Step 5: 執行量子預測（如果需要）
            if use_quantum_for_this and self.use_real_quantum:
                attack_prob, feature_importance = self.quantum_qnn.forward(quantum_features)
            elif not use_quantum_for_this:
                # 已在上面計算過
                pass
            else:
                # 使用模擬量子
                attack_prob, feature_importance = self.classical_qnn.forward(quantum_features)
            
            # Step 5: Calculate trust score (inverse of attack probability)
            trust_score = 1.0 - attack_prob
            
            # Step 6: Determine risk level
            risk_level = self._determine_risk_level(attack_prob)
            
            # Step 7: Identify contributing factors
            contributing_factors = self._analyze_contributing_factors(
                trust_context,
                feature_importance,
                attack_prob
            )
            
            # Step 8: Recommend action
            action = self._recommend_action(risk_level, contributing_factors)
            
            # Step 9: Calculate confidence
            confidence = self._calculate_confidence(quantum_features, attack_prob)
            
            prediction = ZeroTrustPrediction(
                user_id=trust_context.user_id,
                prediction_id=f"ztp_{datetime.now().strftime('%Y%m%d%H%M%S')}_{np.random.randint(1000,9999)}",
                attack_probability=attack_prob,
                trust_score=trust_score,
                risk_level=risk_level,
                contributing_factors=contributing_factors,
                recommended_action=action,
                confidence=confidence,
                timestamp=datetime.now().isoformat(),
                quantum_advantage=use_quantum_for_this  # 記錄是否使用了量子計算
            )
            
            self.prediction_history.append(prediction)
            
            logger.info(
                f"Zero Trust prediction: User={trust_context.user_id}, "
                f"Risk={risk_level}, Prob={attack_prob:.2%}"
            )
            
            return prediction
            
        except Exception as e:
            logger.error(f"Zero Trust prediction failed: {e}")
            return None
    
    def _extract_context_features(self, context) -> np.ndarray:
        """Extract numerical features from trust context"""
        features = np.array([
            context.authentication_strength,
            context.credential_age_hours / 24.0,  # Normalize to days
            context.device_posture_score,
            context.device_trust_level,
            1.0 if context.network_type == 'corporate' else 0.5 if context.network_type == 'home' else 0.0,
            context.access_pattern_anomaly,
            context.resource_access_frequency / 100.0,
            context.time_of_day_anomaly,
            context.geographic_velocity,
            context.threat_intelligence_score,
            context.compliance_status,
            context.data_sensitivity,
            len(context.role_privileges) / 10.0,  # Normalize privilege count
            1.0 - context.authentication_strength,  # Inverse as risk factor
        ])
        
        return features
    
    def _determine_risk_level(self, attack_prob: float) -> str:
        """Determine risk level from attack probability"""
        if attack_prob >= self.thresholds['CRITICAL']:
            return 'CRITICAL'
        elif attack_prob >= self.thresholds['HIGH']:
            return 'HIGH'
        elif attack_prob >= self.thresholds['MEDIUM']:
            return 'MEDIUM'
        else:
            return 'LOW'
    
    def _analyze_contributing_factors(
        self,
        context,
        feature_importance: Dict,
        attack_prob: float
    ) -> Dict[str, float]:
        """Analyze which factors contribute most to risk"""
        factors = {}
        
        # Identity factors
        if context.authentication_strength < 0.6:
            factors['weak_authentication'] = 0.8
        
        if context.credential_age_hours > 72:
            factors['stale_credentials'] = 0.6
        
        # Device factors
        if context.device_posture_score < 0.5:
            factors['poor_device_posture'] = 0.7
        
        if context.network_type == 'public':
            factors['untrusted_network'] = 0.9
        
        # Behavioral factors
        if context.access_pattern_anomaly > 0.5:
            factors['anomalous_behavior'] = context.access_pattern_anomaly
        
        if context.time_of_day_anomaly > 0.7:
            factors['unusual_access_time'] = context.time_of_day_anomaly
        
        if context.geographic_velocity > 0.8:
            factors['impossible_travel'] = 1.0
        
        # Environmental factors
        if context.threat_intelligence_score > 0.5:
            factors['malicious_ip'] = context.threat_intelligence_score
        
        if context.compliance_status < 0.7:
            factors['compliance_violation'] = 1.0 - context.compliance_status
        
        # Add quantum feature importance
        for key, value in feature_importance.items():
            if value > 0.5:
                factors[f'quantum_{key}'] = value
        
        return factors
    
    def _recommend_action(self, risk_level: str, factors: Dict) -> str:
        """Recommend action based on risk level"""
        actions = {
            'CRITICAL': "DENY ACCESS - Immediate investigation required. Possible account compromise.",
            'HIGH': "REQUIRE STEP-UP AUTH - Additional verification needed before granting access.",
            'MEDIUM': "ALLOW WITH MONITORING - Grant limited access with enhanced logging.",
            'LOW': "ALLOW - Normal access granted with standard monitoring."
        }
        
        action = actions.get(risk_level, "ALLOW")
        
        # Add specific recommendations based on factors
        if 'impossible_travel' in factors:
            action += " [Geographic anomaly detected]"
        
        if 'weak_authentication' in factors:
            action += " [Enforce MFA]"
        
        if 'malicious_ip' in factors:
            action += " [Block suspicious IP]"
        
        return action
    
    def _calculate_confidence(self, quantum_features: np.ndarray, prediction: float) -> float:
        """Calculate prediction confidence using quantum metrics"""
        # Use quantum entropy as confidence measure
        probabilities = np.abs(quantum_features)**2
        probabilities = probabilities / (np.sum(probabilities) + 1e-10)
        
        # Shannon entropy
        entropy = -np.sum(probabilities * np.log2(probabilities + 1e-10))
        
        # Normalize entropy to confidence (lower entropy = higher confidence)
        max_entropy = np.log2(len(quantum_features))
        confidence = 1.0 - (entropy / max_entropy)
        
        # Adjust based on prediction certainty
        if prediction > 0.9 or prediction < 0.1:
            confidence += 0.1  # More confident at extremes
        
        return float(np.clip(confidence, 0.0, 1.0))


class QuantumEnsemblePredictor:
    """Ensemble of quantum-ML models for robust prediction"""
    
    def __init__(self, n_models: int = 5):
        """初始化 Quantum Ensemble Predictor"""
        self.models = [
            QuantumNeuralNetwork(input_dim=20, quantum_dim=10)
            for _ in range(n_models)
        ]
        self.model_weights = np.ones(n_models) / n_models
        logger.info(f"Quantum Ensemble with {n_models} models initialized")
    
    async def ensemble_predict(
        self,
        features: np.ndarray
    ) -> Tuple[float, float, Dict]:
        """Predict using ensemble voting"""
        predictions = []
        importances = []
        
        for model in self.models:
            pred, importance = model.forward(features)
            predictions.append(pred)
            importances.append(importance)
        
        # Weighted voting
        ensemble_prediction = np.average(predictions, weights=self.model_weights)
        
        # Ensemble uncertainty (variance)
        ensemble_uncertainty = np.var(predictions)
        
        # Aggregate feature importance
        avg_importance = {}
        for key in importances[0].keys():
            values = [imp.get(key, 0) for imp in importances]
            avg_importance[key] = float(np.mean(values))
        
        return float(ensemble_prediction), float(ensemble_uncertainty), avg_importance
    
    def update_model_weights(self, predictions: List[float], ground_truth: float):
        """Update model weights based on performance (online learning)"""
        # Calculate error for each model
        errors = [abs(pred - ground_truth) for pred in predictions]
        
        # Update weights (inverse of error)
        self.model_weights = 1.0 / (np.array(errors) + 1e-6)
        self.model_weights = self.model_weights / np.sum(self.model_weights)
        
        logger.info(f"Model weights updated: {self.model_weights}")


class QuantumAnomalyFusion:
    """Fuse multiple anomaly scores using quantum interference"""
    
    def __init__(self):
        """初始化 Quantum Anomaly Fusion"""
        logger.info("Quantum Anomaly Fusion initialized")
    
    def fuse_anomaly_scores(
        self,
        anomaly_scores: Dict[str, float]
    ) -> Tuple[float, Dict[str, float]]:
        """
        Fuse multiple anomaly scores using quantum interference
        
        Quantum advantage: Constructive/destructive interference
        amplifies strong signals and suppresses weak ones
        """
        
        # Convert scores to quantum amplitudes
        scores_array = np.array(list(anomaly_scores.values()))
        
        # Normalize to unit circle (quantum state)
        normalized = scores_array / (np.linalg.norm(scores_array) + 1e-6)
        
        # Apply quantum phase encoding
        phases = np.exp(1j * normalized * np.pi)
        quantum_state = normalized * phases
        
        # Quantum interference (superposition)
        interference = np.sum(quantum_state)
        
        # Measure final anomaly score
        fused_score = float(np.abs(interference))
        
        # Calculate contribution weights
        contributions = {}
        for i, (key, _) in enumerate(anomaly_scores.items()):
            contribution = abs(quantum_state[i]) / (np.sum(np.abs(quantum_state)) + 1e-10)
            contributions[key] = float(contribution)
        
        return fused_score, contributions


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora Quantum-ML Hybrid 測試 ===")
    
    # 測試 Quantum Neural Network
    logger.info("\n--- Quantum Neural Network ---")
    qnn = QuantumNeuralNetwork(input_dim=14, quantum_dim=10)
    
    test_features = np.random.randn(14)
    prediction, importance = qnn.forward(test_features)
    
    print(f"預測值: {prediction:.2%}")
    print(f"特徵重要性: {importance}")
    
    # 測試 Ensemble
    logger.info("\n--- Quantum Ensemble ---")
    ensemble = QuantumEnsemblePredictor(n_models=3)
    
    ensemble_pred, uncertainty, avg_imp = await ensemble.ensemble_predict(test_features)
    
    print(f"集成預測: {ensemble_pred:.2%}")
    print(f"不確定性: {uncertainty:.4f}")
    
    logger.info("\n=== 測試完成 ===")


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())

