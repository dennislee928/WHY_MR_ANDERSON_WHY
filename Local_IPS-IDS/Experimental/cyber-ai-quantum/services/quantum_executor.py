#!/usr/bin/env python3
"""
Pandora Quantum Executor Service
處理所有與 IBM Quantum 後端的交互
"""

import os
import logging
import asyncio
from typing import Dict, List, Optional, Tuple
from datetime import datetime
from enum import Enum
import json
import numpy as np

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


class JobStatus(str, Enum):
    """作業狀態"""
    QUEUED = "QUEUED"
    RUNNING = "RUNNING"
    DONE = "DONE"
    ERROR = "ERROR"
    CANCELLED = "CANCELLED"


class QuantumBackendType(str, Enum):
    """量子後端類型"""
    LOCAL_SIMULATOR = "local_simulator"
    CLOUD_SIMULATOR = "cloud_simulator"
    REAL_HARDWARE = "real_hardware"


class QuantumExecutorService:
    """
    量子執行器服務
    
    負責:
    1. 管理 IBM Quantum 連接
    2. 提交量子作業
    3. 追蹤作業狀態
    4. 獲取結果
    """
    
    def __init__(self):
        """初始化 Quantum Executor Service"""
        self.jobs = {}  # job_id -> job_info
        self.results_cache = {}  # job_id -> results
        self.ibm_service = None
        self.current_backend = None
        
        # 配置
        self.use_real_quantum = os.getenv('QUANTUM_REAL_HARDWARE_ENABLED', 'false').lower() == 'true'
        self.backend_name = os.getenv('QUANTUM_BACKEND', 'simulator_statevector')
        self.job_timeout = int(os.getenv('QUANTUM_JOB_TIMEOUT_SECONDS', '300'))
        self.job_retry_count = int(os.getenv('QUANTUM_JOB_RETRY_COUNT', '3'))
        
        logger.info("Quantum Executor Service initialized")
        
        # 嘗試初始化後端
        self._initialize_backend()
    
    def _initialize_backend(self):
        """初始化量子後端"""
        try:
            # 首先檢查 Qiskit 是否可用
            try:
                from qiskit_aer.primitives import Sampler as AerSampler
                from qiskit_aer import AerSimulator
                self.qiskit_available = True
            except ImportError:
                self.qiskit_available = False
                logger.warning("Qiskit Aer not available, using fallback mode")
                return
            
            if self.use_real_quantum:
                # IBM Quantum 雲端/硬體
                logger.info("Initializing IBM Quantum backend...")
                self._initialize_ibm_quantum()
            else:
                # 本地模擬器
                logger.info(f"Using local simulator: {self.backend_name}")
                self.current_backend = AerSimulator(method='statevector')
                self.sampler = AerSampler()
                self.backend_type = QuantumBackendType.LOCAL_SIMULATOR
                
        except Exception as e:
            logger.error(f"Backend initialization failed: {e}")
            self.qiskit_available = False
    
    def _initialize_ibm_quantum(self):
        """初始化 IBM Quantum 服務"""
        try:
            from qiskit_ibm_runtime import QiskitRuntimeService, Sampler as RuntimeSampler
            
            # 載入 Token
            token = os.getenv('IBM_QUANTUM_TOKEN')
            if not token or token == 'your_ibm_quantum_api_token_here':
                logger.warning("IBM Quantum Token not configured, falling back to simulator")
                from qiskit_aer.primitives import Sampler as AerSampler
                self.sampler = AerSampler()
                self.backend_type = QuantumBackendType.LOCAL_SIMULATOR
                return
            
            # 初始化 IBM Quantum Service
            self.ibm_service = QiskitRuntimeService(
                channel='ibm_quantum',
                token=token
            )
            
            # 選擇後端
            if 'simulator' in self.backend_name:
                # 雲端模擬器
                backend = self.ibm_service.get_backend(self.backend_name)
                self.backend_type = QuantumBackendType.CLOUD_SIMULATOR
            else:
                # 真實量子硬體
                backend = self.ibm_service.least_busy(operational=True, simulator=False)
                self.backend_type = QuantumBackendType.REAL_HARDWARE
            
            self.current_backend = backend
            self.sampler = RuntimeSampler(backend=backend)
            
            logger.info(f"IBM Quantum backend initialized: {backend.name}")
            
        except Exception as e:
            logger.error(f"IBM Quantum initialization failed: {e}")
            # Fallback to local simulator
            from qiskit_aer.primitives import Sampler as AerSampler
            self.sampler = AerSampler()
            self.backend_type = QuantumBackendType.LOCAL_SIMULATOR
    
    async def submit_quantum_job(
        self,
        circuit,  # QuantumCircuit
        parameters: Dict,
        job_metadata: Optional[Dict] = None
    ) -> str:
        """
        提交量子作業
        
        Args:
            circuit: Qiskit QuantumCircuit
            parameters: 電路參數
            job_metadata: 作業元數據
        
        Returns:
            job_id: 作業ID
        """
        try:
            if not self.qiskit_available:
                raise RuntimeError("Qiskit not available")
            
            # 生成作業ID
            job_id = f"qjob_{datetime.now().strftime('%Y%m%d%H%M%S')}_{np.random.randint(1000,9999)}"
            
            # 作業資訊
            job_info = {
                'job_id': job_id,
                'status': JobStatus.QUEUED,
                'submitted_at': datetime.now().isoformat(),
                'backend': self.current_backend.name if self.current_backend else 'local',
                'backend_type': self.backend_type.value,
                'metadata': job_metadata or {},
                'parameters': parameters
            }
            
            self.jobs[job_id] = job_info
            
            # 異步執行作業
            asyncio.create_task(self._execute_job(job_id, circuit, parameters))
            
            logger.info(f"Quantum job submitted: {job_id}")
            
            return job_id
            
        except Exception as e:
            logger.error(f"Job submission failed: {e}")
            raise
    
    async def _execute_job(self, job_id: str, circuit, parameters: Dict):
        """執行量子作業（異步）"""
        try:
            # 更新狀態
            self.jobs[job_id]['status'] = JobStatus.RUNNING
            self.jobs[job_id]['started_at'] = datetime.now().isoformat()
            
            logger.info(f"Executing quantum job: {job_id}")
            
            # 執行電路（在線程池中運行以避免阻塞）
            loop = asyncio.get_event_loop()
            result = await loop.run_in_executor(
                None,
                self._run_circuit_sync,
                circuit,
                parameters
            )
            
            # 儲存結果
            self.results_cache[job_id] = result
            self.jobs[job_id]['status'] = JobStatus.DONE
            self.jobs[job_id]['completed_at'] = datetime.now().isoformat()
            
            # 計算執行時間
            started = datetime.fromisoformat(self.jobs[job_id]['started_at'])
            execution_time = (datetime.now() - started).total_seconds()
            self.jobs[job_id]['execution_time_seconds'] = execution_time
            
            logger.info(f"Job completed: {job_id} (time={execution_time:.2f}s)")
            
        except Exception as e:
            logger.error(f"Job execution failed: {job_id} - {e}")
            self.jobs[job_id]['status'] = JobStatus.ERROR
            self.jobs[job_id]['error'] = str(e)
    
    def _run_circuit_sync(self, circuit, parameters: Dict):
        """同步執行電路（在執行器中調用）"""
        # 這裡執行真實的量子計算
        # 簡化版本：直接返回模擬結果
        
        if not self.qiskit_available:
            return {'probabilities': [0.5, 0.5], 'counts': {}, 'method': 'fallback'}
        
        try:
            from qiskit import transpile
            
            # 轉譯電路for後端
            if self.current_backend:
                transpiled = transpile(circuit, self.current_backend, optimization_level=3)
            else:
                transpiled = circuit
            
            # 使用 Sampler 執行
            # 注意：實際實現會更複雜，需要處理參數綁定等
            job = self.sampler.run([transpiled])
            result = job.result()
            
            # 提取機率
            quasi_dists = result.quasi_dists
            
            return {
                'quasi_distributions': quasi_dists,
                'metadata': result.metadata,
                'method': 'qiskit_execution'
            }
            
        except Exception as e:
            logger.error(f"Circuit execution error: {e}")
            return {'error': str(e), 'method': 'failed'}
    
    async def get_job_status(self, job_id: str) -> Dict:
        """獲取作業狀態"""
        if job_id not in self.jobs:
            return {'error': 'Job not found'}
        
        job_info = self.jobs[job_id]
        
        return {
            'job_id': job_id,
            'status': job_info['status'],
            'submitted_at': job_info.get('submitted_at'),
            'started_at': job_info.get('started_at'),
            'completed_at': job_info.get('completed_at'),
            'execution_time_seconds': job_info.get('execution_time_seconds'),
            'backend': job_info.get('backend'),
            'backend_type': job_info.get('backend_type')
        }
    
    async def get_job_result(self, job_id: str) -> Dict:
        """獲取作業結果"""
        if job_id not in self.jobs:
            return {'error': 'Job not found'}
        
        job_status = self.jobs[job_id]['status']
        
        if job_status != JobStatus.DONE:
            return {
                'job_id': job_id,
                'status': job_status,
                'message': f'Job is {job_status}, not ready yet'
            }
        
        if job_id in self.results_cache:
            return {
                'job_id': job_id,
                'status': JobStatus.DONE,
                'result': self.results_cache[job_id]
            }
        
        return {'error': 'Result not found'}
    
    async def cancel_job(self, job_id: str) -> Dict:
        """取消作業"""
        if job_id not in self.jobs:
            return {'error': 'Job not found'}
        
        if self.jobs[job_id]['status'] in [JobStatus.DONE, JobStatus.ERROR]:
            return {'error': f"Cannot cancel job in {self.jobs[job_id]['status']} status"}
        
        self.jobs[job_id]['status'] = JobStatus.CANCELLED
        self.jobs[job_id]['cancelled_at'] = datetime.now().isoformat()
        
        logger.info(f"Job cancelled: {job_id}")
        
        return {
            'job_id': job_id,
            'status': JobStatus.CANCELLED
        }
    
    def get_all_jobs(self, limit: int = 50) -> List[Dict]:
        """獲取所有作業列表"""
        jobs_list = list(self.jobs.values())
        
        # 按提交時間排序（最新的在前）
        jobs_list.sort(key=lambda x: x.get('submitted_at', ''), reverse=True)
        
        return jobs_list[:limit]
    
    def get_statistics(self) -> Dict:
        """獲取執行器統計"""
        total_jobs = len(self.jobs)
        
        status_counts = {
            JobStatus.QUEUED: 0,
            JobStatus.RUNNING: 0,
            JobStatus.DONE: 0,
            JobStatus.ERROR: 0,
            JobStatus.CANCELLED: 0
        }
        
        execution_times = []
        
        for job in self.jobs.values():
            status_counts[job['status']] += 1
            
            if 'execution_time_seconds' in job:
                execution_times.append(job['execution_time_seconds'])
        
        avg_execution_time = np.mean(execution_times) if execution_times else 0.0
        
        return {
            'total_jobs': total_jobs,
            'status_distribution': {k.value: v for k, v in status_counts.items()},
            'average_execution_time_seconds': float(avg_execution_time),
            'backend_type': self.backend_type.value if self.backend_type else 'unknown',
            'qiskit_available': self.qiskit_available,
            'quantum_cache_size': len(self.results_cache)
        }


# 全局執行器實例
_quantum_executor = None


def get_quantum_executor() -> QuantumExecutorService:
    """獲取全局量子執行器實例"""
    global _quantum_executor
    
    if _quantum_executor is None:
        _quantum_executor = QuantumExecutorService()
    
    return _quantum_executor


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora Quantum Executor Service 測試 ===")
    
    executor = get_quantum_executor()
    
    # 顯示統計
    stats = executor.get_statistics()
    print(f"\n執行器統計:")
    print(json.dumps(stats, indent=2))
    
    logger.info("\n=== 測試完成 ===")


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())

