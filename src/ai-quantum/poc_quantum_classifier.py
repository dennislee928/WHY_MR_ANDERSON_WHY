#!/usr/bin/env python3
"""
Pandora Real Quantum Classifier - Proof of Concept
使用 Qiskit 和 IBM Quantum 實現真實量子威脅分類

這是從模擬到真實量子計算的概念驗證實現
"""

import numpy as np
import logging
from typing import Tuple, Dict, List, Optional
from datetime import datetime
import os

# Qiskit imports
try:
    from qiskit import QuantumCircuit, transpile
    from qiskit.circuit.library import RealAmplitudes, ZZFeatureMap, PauliFeatureMap
    from qiskit_aer.primitives import Sampler as AerSampler
    from qiskit_machine_learning.neural_networks import SamplerQNN
    from qiskit_machine_learning.algorithms import VQC
    from qiskit.algorithms.optimizers import COBYLA, SPSA, L_BFGS_B
    from qiskit.primitives import Sampler
    
    QISKIT_AVAILABLE = True
except ImportError:
    QISKIT_AVAILABLE = False
    print("⚠️  Qiskit 未安裝。請運行: pip install qiskit qiskit-aer qiskit-machine-learning")

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class RealQuantumNeuralNetwork:
    """
    真實量子神經網絡 (使用 Qiskit)
    
    架構:
    1. 古典預處理層
    2. 量子特徵映射 (ZZFeatureMap)
    3. 變分量子電路 (RealAmplitudes)
    4. 量子測量
    5. 古典後處理層
    """
    
    def __init__(
        self,
        num_features: int = 4,
        num_qubits: int = 4,
        feature_map_reps: int = 2,
        ansatz_reps: int = 3,
        backend_type: str = "simulator"
    ):
        """
        初始化真實量子神經網絡
        
        Args:
            num_features: 輸入特徵數量
            num_qubits: 量子位元數量
            feature_map_reps: 特徵映射重複次數
            ansatz_reps: Ansatz 重複次數
            backend_type: 'simulator' 或 'hardware'
        """
        if not QISKIT_AVAILABLE:
            raise ImportError("Qiskit is required for real quantum computation")
        
        self.num_features = num_features
        self.num_qubits = num_qubits
        self.feature_map_reps = feature_map_reps
        self.ansatz_reps = ansatz_reps
        self.backend_type = backend_type
        
        # 構建量子電路
        self.feature_map = self._build_feature_map()
        self.ansatz = self._build_ansatz()
        self.quantum_circuit = self._build_quantum_circuit()
        
        # 初始化變分參數（隨機）
        num_params = self.ansatz.num_parameters
        self.quantum_weights = np.random.randn(num_params) * 0.1
        
        # 選擇後端
        self.sampler = self._initialize_backend()
        
        # 構建 Sampler QNN
        self.qnn = self._build_sampler_qnn()
        
        logger.info(
            f"Real Quantum Neural Network initialized: "
            f"{num_qubits} qubits, {num_params} parameters, backend={backend_type}"
        )
    
    def _build_feature_map(self) -> QuantumCircuit:
        """構建量子特徵映射"""
        # ZZFeatureMap: 使用 ZZ 交互作用編碼特徵
        feature_map = ZZFeatureMap(
            feature_dimension=self.num_features,
            reps=self.feature_map_reps,
            entanglement='linear'  # 或 'full' for more entanglement
        )
        
        logger.info(f"Feature map created: {self.num_features} features, {self.feature_map_reps} reps")
        return feature_map
    
    def _build_ansatz(self) -> QuantumCircuit:
        """構建變分Ansatz"""
        # RealAmplitudes: 常用的變分形式
        ansatz = RealAmplitudes(
            num_qubits=self.num_qubits,
            reps=self.ansatz_reps,
            entanglement='linear'
        )
        
        logger.info(
            f"Ansatz created: {self.num_qubits} qubits, "
            f"{ansatz.num_parameters} parameters, {self.ansatz_reps} reps"
        )
        return ansatz
    
    def _build_quantum_circuit(self) -> QuantumCircuit:
        """組合完整的量子電路"""
        qc = QuantumCircuit(self.num_qubits)
        
        # 組合特徵映射和Ansatz
        qc.compose(self.feature_map, inplace=True)
        qc.compose(self.ansatz, inplace=True)
        
        # 添加測量（可選，取決於使用方式）
        # qc.measure_all()
        
        logger.info(f"Quantum circuit composed: depth={qc.depth()}")
        return qc
    
    def _initialize_backend(self):
        """初始化量子後端"""
        if self.backend_type == "simulator":
            # 本地模擬器（快速，用於開發）
            sampler = AerSampler()
            logger.info("Using local Aer Sampler for simulation")
        else:
            # IBM Quantum 硬體或雲端模擬器
            from qiskit_ibm_runtime import QiskitRuntimeService, Sampler as RuntimeSampler
            
            # 載入 IBM Quantum Token
            token = os.getenv('IBM_QUANTUM_TOKEN')
            if not token:
                raise ValueError("IBM_QUANTUM_TOKEN not found in environment")
            
            service = QiskitRuntimeService(
                channel='ibm_quantum',
                token=token
            )
            
            # 選擇後端
            backend = service.least_busy(operational=True, simulator=False)
            sampler = RuntimeSampler(backend=backend)
            
            logger.info(f"Using IBM Quantum backend: {backend.name}")
        
        return sampler
    
    def _build_sampler_qnn(self) -> SamplerQNN:
        """構建 Sampler Quantum Neural Network"""
        # 創建 SamplerQNN
        qnn = SamplerQNN(
            circuit=self.quantum_circuit,
            input_params=self.feature_map.parameters,
            weight_params=self.ansatz.parameters,
            sampler=self.sampler
        )
        
        return qnn
    
    def forward(self, features: np.ndarray) -> Tuple[float, Dict]:
        """
        前向傳播
        
        Args:
            features: 輸入特徵向量 (應為 num_features 維度)
        
        Returns:
            prediction: 預測值 (0-1)
            metadata: 額外資訊
        """
        try:
            # 確保特徵維度匹配
            if len(features) != self.num_features:
                # 填充或截斷
                if len(features) < self.num_features:
                    features = np.pad(features, (0, self.num_features - len(features)))
                else:
                    features = features[:self.num_features]
            
            # 正規化特徵到 [0, π]
            normalized_features = (features - np.min(features)) / (np.max(features) - np.min(features) + 1e-10)
            normalized_features = normalized_features * np.pi
            
            # 執行量子電路
            # SamplerQNN.forward() 需要 (input, weights) 參數
            quantum_output = self.qnn.forward(
                normalized_features.reshape(1, -1),
                self.quantum_weights.reshape(1, -1)
            )
            
            # 處理輸出
            prediction = quantum_output[0][0] if len(quantum_output.shape) > 1 else quantum_output[0]
            
            # 確保在 [0, 1] 範圍
            prediction = float(np.clip(prediction, 0, 1))
            
            metadata = {
                'quantum_execution': True,
                'num_qubits': self.num_qubits,
                'num_parameters': len(self.quantum_weights),
                'backend_type': self.backend_type,
                'timestamp': datetime.now().isoformat()
            }
            
            return prediction, metadata
            
        except Exception as e:
            logger.error(f"Quantum forward pass failed: {e}")
            # Fallback to random prediction
            return 0.5, {'error': str(e), 'quantum_execution': False}
    
    def train(
        self,
        X_train: np.ndarray,
        y_train: np.ndarray,
        optimizer: str = 'COBYLA',
        max_iter: int = 100
    ) -> Dict:
        """
        訓練量子神經網絡
        
        Args:
            X_train: 訓練特徵 (n_samples, num_features)
            y_train: 訓練標籤 (n_samples,)
            optimizer: 優化器名稱
            max_iter: 最大迭代次數
        
        Returns:
            training_results: 訓練結果
        """
        try:
            logger.info(f"Training quantum model with {len(X_train)} samples...")
            
            # 選擇優化器
            if optimizer == 'COBYLA':
                opt = COBYLA(maxiter=max_iter)
            elif optimizer == 'SPSA':
                opt = SPSA(maxiter=max_iter)
            else:
                opt = L_BFGS_B(maxfun=max_iter)
            
            # 創建 VQC (Variational Quantum Classifier)
            vqc = VQC(
                num_qubits=self.num_qubits,
                feature_map=self.feature_map,
                ansatz=self.ansatz,
                optimizer=opt,
                sampler=self.sampler
            )
            
            # 訓練
            start_time = datetime.now()
            vqc.fit(X_train, y_train)
            training_time = (datetime.now() - start_time).total_seconds()
            
            # 更新權重
            self.quantum_weights = vqc.weights
            
            # 評估訓練集
            train_predictions = vqc.predict(X_train)
            train_accuracy = np.mean(train_predictions == y_train)
            
            results = {
                'training_time_seconds': training_time,
                'train_accuracy': float(train_accuracy),
                'num_parameters_trained': len(self.quantum_weights),
                'optimizer': optimizer,
                'max_iterations': max_iter,
                'quantum_execution': True
            }
            
            logger.info(
                f"Training completed: accuracy={train_accuracy:.2%}, "
                f"time={training_time:.1f}s"
            )
            
            return results
            
        except Exception as e:
            logger.error(f"Training failed: {e}")
            return {'error': str(e), 'quantum_execution': False}
    
    def get_circuit_info(self) -> Dict:
        """獲取電路資訊"""
        return {
            'num_qubits': self.num_qubits,
            'num_parameters': self.ansatz.num_parameters,
            'circuit_depth': self.quantum_circuit.depth(),
            'circuit_width': self.quantum_circuit.width(),
            'num_gates': sum(self.quantum_circuit.count_ops().values()),
            'gate_types': dict(self.quantum_circuit.count_ops()),
            'feature_map_type': 'ZZFeatureMap',
            'ansatz_type': 'RealAmplitudes'
        }


class QuantumThreatClassifier:
    """
    完整的量子威脅分類器
    結合古典預處理和真實量子計算
    """
    
    def __init__(self, use_real_quantum: bool = False):
        """
        初始化量子威脅分類器
        
        Args:
            use_real_quantum: True = 使用真實量子硬體/雲端, False = 本地模擬
        """
        self.use_real_quantum = use_real_quantum
        
        # 古典預處理層 (20 -> 4 特徵壓縮)
        self.classical_weights = np.random.randn(20, 4) * 0.1
        self.classical_bias = np.zeros(4)
        
        # 量子層
        backend_type = "hardware" if use_real_quantum else "simulator"
        self.quantum_layer = RealQuantumNeuralNetwork(
            num_features=4,
            num_qubits=4,
            feature_map_reps=2,
            ansatz_reps=3,
            backend_type=backend_type
        )
        
        # 古典後處理層 (量子輸出 -> 預測)
        self.output_weights = np.random.randn(1) * 0.1
        self.output_bias = 0.0
        
        logger.info(
            f"Quantum Threat Classifier initialized "
            f"(real_quantum={use_real_quantum})"
        )
    
    def _classical_preprocessing(self, features: np.ndarray) -> np.ndarray:
        """古典預處理：降維"""
        z = np.dot(features, self.classical_weights) + self.classical_bias
        return np.tanh(z)  # Activation
    
    def _classical_postprocessing(self, quantum_output: float) -> float:
        """古典後處理"""
        output = quantum_output * self.output_weights[0] + self.output_bias
        return float(1 / (1 + np.exp(-output)))  # Sigmoid
    
    async def predict(self, features: np.ndarray) -> Dict:
        """
        預測威脅
        
        Args:
            features: 20維特徵向量
        
        Returns:
            prediction_result: 包含預測機率、元數據等
        """
        try:
            start_time = datetime.now()
            
            # 步驟 1: 古典預處理
            compressed_features = self._classical_preprocessing(features)
            
            # 步驟 2: 量子處理 (真實量子計算!)
            quantum_output, quantum_metadata = self.quantum_layer.forward(compressed_features)
            
            # 步驟 3: 古典後處理
            final_prediction = self._classical_postprocessing(quantum_output)
            
            execution_time = (datetime.now() - start_time).total_seconds()
            
            result = {
                'attack_probability': final_prediction,
                'quantum_contribution': quantum_output,
                'execution_time_seconds': execution_time,
                'quantum_metadata': quantum_metadata,
                'model_type': 'Quantum-Classical-Hybrid',
                'qiskit_version': 'integrated',
                'timestamp': datetime.now().isoformat()
            }
            
            logger.info(
                f"Quantum prediction: prob={final_prediction:.2%}, "
                f"time={execution_time:.3f}s"
            )
            
            return result
            
        except Exception as e:
            logger.error(f"Quantum prediction failed: {e}")
            return {
                'error': str(e),
                'fallback': True,
                'attack_probability': 0.5
            }
    
    def train_quantum_layer(
        self,
        X_train: np.ndarray,
        y_train: np.ndarray
    ) -> Dict:
        """
        訓練量子層
        
        Args:
            X_train: 訓練特徵 (n_samples, 20)
            y_train: 訓練標籤 (n_samples,)
        """
        # 預處理所有訓練數據
        X_compressed = np.array([
            self._classical_preprocessing(x) for x in X_train
        ])
        
        # 訓練量子層
        results = self.quantum_layer.train(
            X_compressed,
            y_train,
            optimizer='COBYLA',
            max_iter=50
        )
        
        return results


# ========== 測試和基準測試 ==========

def generate_test_dataset(n_samples: int = 20) -> Tuple[np.ndarray, np.ndarray]:
    """
    生成測試數據集
    
    Returns:
        X: 特徵 (n_samples, 20)
        y: 標籤 (n_samples,) - 0 = 正常, 1 = 攻擊
    """
    np.random.seed(42)
    
    X = []
    y = []
    
    for i in range(n_samples):
        if i < n_samples // 2:
            # 正常流量
            features = np.random.randn(20) * 0.3 + 0.3
            label = 0
        else:
            # 攻擊流量
            features = np.random.randn(20) * 0.4 + 0.7
            label = 1
        
        X.append(features)
        y.append(label)
    
    return np.array(X), np.array(y)


async def benchmark_quantum_vs_classical():
    """基準測試：量子 vs 古典"""
    logger.info("=== Quantum vs Classical Benchmark ===")
    
    # 生成測試數據
    X_test, y_test = generate_test_dataset(n_samples=10)
    
    logger.info(f"Test dataset: {len(X_test)} samples")
    
    # 1. 古典基線（使用 NumPy 模擬）
    logger.info("\n--- Classical Baseline (NumPy Simulation) ---")
    from quantum_ml_hybrid import QuantumNeuralNetwork as ClassicalQNN
    
    classical_model = ClassicalQNN(input_dim=20, quantum_dim=4)
    
    classical_results = []
    classical_start = datetime.now()
    
    for features in X_test:
        pred, _ = classical_model.forward(features)
        classical_results.append(pred)
    
    classical_time = (datetime.now() - classical_start).total_seconds()
    classical_avg_time = classical_time / len(X_test)
    
    print(f"Classical predictions: {classical_results[:5]}")
    print(f"Classical total time: {classical_time:.3f}s")
    print(f"Classical avg time per prediction: {classical_avg_time*1000:.1f}ms")
    
    # 2. 量子模型（使用 Qiskit）
    if QISKIT_AVAILABLE:
        logger.info("\n--- Quantum Model (Qiskit Real Quantum) ---")
        
        quantum_model = QuantumThreatClassifier(use_real_quantum=False)
        
        quantum_results = []
        quantum_start = datetime.now()
        
        for features in X_test:
            result = await quantum_model.predict(features)
            quantum_results.append(result['attack_probability'])
        
        quantum_time = (datetime.now() - quantum_start).total_seconds()
        quantum_avg_time = quantum_time / len(X_test)
        
        print(f"Quantum predictions: {quantum_results[:5]}")
        print(f"Quantum total time: {quantum_time:.3f}s")
        print(f"Quantum avg time per prediction: {quantum_avg_time*1000:.1f}ms")
        
        # 比較
        logger.info("\n--- Comparison ---")
        print(f"Speedup: {quantum_time/classical_time:.2f}x")
        print(f"Quantum advantage: {'Yes' if quantum_time < classical_time else 'No (expected for small dataset)'}")
    
    logger.info("\n=== Benchmark Complete ===")


async def test_circuit_visualization():
    """測試電路可視化"""
    if not QISKIT_AVAILABLE:
        print("Qiskit not available")
        return
    
    logger.info("=== Circuit Visualization Test ===")
    
    qnn = RealQuantumNeuralNetwork(
        num_features=4,
        num_qubits=4,
        feature_map_reps=1,
        ansatz_reps=2
    )
    
    # 獲取電路資訊
    circuit_info = qnn.get_circuit_info()
    
    print("\n電路資訊:")
    for key, value in circuit_info.items():
        print(f"  {key}: {value}")
    
    # 繪製電路（需要 matplotlib）
    try:
        import matplotlib.pyplot as plt
        fig = qnn.quantum_circuit.draw('mpl')
        plt.savefig('quantum_circuit_diagram.png')
        print("\n電路圖已保存: quantum_circuit_diagram.png")
    except ImportError:
        print("\n無法繪製電路圖（需要安裝 matplotlib）")
    
    logger.info("=== Test Complete ===")


async def main():
    """主函數"""
    logger.info("=== Pandora Real Quantum Classifier PoC ===\n")
    
    if not QISKIT_AVAILABLE:
        logger.error("Qiskit is not installed. Please run: pip install qiskit qiskit-aer qiskit-machine-learning")
        return
    
    # 測試 1: 電路可視化
    await test_circuit_visualization()
    
    print("\n" + "="*60 + "\n")
    
    # 測試 2: 基準測試
    await benchmark_quantum_vs_classical()
    
    print("\n" + "="*60)
    print("\n✅ PoC 完成！")
    print("\n下一步:")
    print("  1. 設置 IBM Quantum Token")
    print("  2. 將 use_real_quantum=True 測試雲端執行")
    print("  3. 整合到 main.py FastAPI 服務")


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())

