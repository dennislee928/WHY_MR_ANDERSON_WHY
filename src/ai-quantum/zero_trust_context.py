#!/usr/bin/env python3
"""
Pandora Zero Trust Context Engine
Aggregates multi-dimensional trust signals for quantum-ML analysis
"""

import numpy as np
from dataclasses import dataclass
from typing import Dict, List, Optional
from datetime import datetime, timedelta
import logging
from math import radians, sin, cos, sqrt, asin

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)


@dataclass
class TrustContext:
    """Zero Trust Context Snapshot"""
    user_id: str
    device_id: str
    location_coordinates: tuple  # (lat, lon)
    timestamp: str
    
    # Identity Factors (0-1 normalized)
    authentication_strength: float  # MFA, biometric, etc.
    credential_age_hours: float
    role_privileges: List[str]
    
    # Device Factors
    device_posture_score: float  # OS patches, AV status
    device_trust_level: float  # Known device, compliant
    network_type: str  # corporate, home, public
    
    # Behavioral Factors
    access_pattern_anomaly: float  # Deviation from baseline
    resource_access_frequency: float
    time_of_day_anomaly: float
    geographic_velocity: float  # km/h between logins
    
    # Environmental Factors
    threat_intelligence_score: float  # IP reputation
    compliance_status: float  # Policy adherence
    data_sensitivity: float  # Resource classification


class ZeroTrustContextAggregator:
    """Aggregates trust signals from multiple sources"""
    
    def __init__(self):
        """初始化 Zero Trust Context Aggregator"""
        self.context_history = []
        self.baseline_profiles = {}
        logger.info("Zero Trust Context Aggregator initialized")
    
    async def collect_context(
        self, 
        user_id: str, 
        device_id: str,
        session_data: Dict
    ) -> TrustContext:
        """Collect comprehensive trust context"""
        
        # Identity signals
        auth_strength = self._calculate_auth_strength(session_data)
        
        # Device signals
        device_score = await self._assess_device_posture(device_id)
        
        # Behavioral signals
        behavior_anomaly = self._detect_behavioral_anomaly(user_id, session_data)
        
        # Environmental signals
        threat_score = await self._query_threat_intelligence(
            session_data.get('source_ip', '0.0.0.0')
        )
        
        context = TrustContext(
            user_id=user_id,
            device_id=device_id,
            location_coordinates=session_data.get('location', (0, 0)),
            timestamp=datetime.now().isoformat(),
            authentication_strength=auth_strength,
            credential_age_hours=session_data.get('credential_age', 0),
            role_privileges=session_data.get('roles', []),
            device_posture_score=device_score,
            device_trust_level=session_data.get('device_trust', 0.5),
            network_type=session_data.get('network_type', 'unknown'),
            access_pattern_anomaly=behavior_anomaly,
            resource_access_frequency=session_data.get('access_freq', 0),
            time_of_day_anomaly=self._time_anomaly(session_data),
            geographic_velocity=self._calculate_velocity(user_id, session_data),
            threat_intelligence_score=threat_score,
            compliance_status=session_data.get('compliance', 1.0),
            data_sensitivity=session_data.get('data_sensitivity', 0.5)
        )
        
        self.context_history.append(context)
        return context
    
    def _calculate_auth_strength(self, session_data: Dict) -> float:
        """Calculate authentication strength score"""
        score = 0.0
        
        # Base authentication
        if session_data.get('password_auth'):
            score += 0.3
        
        # Multi-factor
        if session_data.get('mfa_enabled'):
            score += 0.3
        
        # Biometric
        if session_data.get('biometric_auth'):
            score += 0.4
        
        return min(score, 1.0)
    
    async def _assess_device_posture(self, device_id: str) -> float:
        """Assess device security posture"""
        # Simplified - in production, query MDM/EDR systems
        checks = {
            'os_patched': 0.25,
            'antivirus_active': 0.25,
            'disk_encrypted': 0.25,
            'firewall_enabled': 0.25
        }
        
        # Simulate device checks
        score = sum(np.random.random() > 0.3 for _ in checks) * 0.25
        return score
    
    def _detect_behavioral_anomaly(self, user_id: str, session_data: Dict) -> float:
        """Detect anomalous behavior patterns"""
        if user_id not in self.baseline_profiles:
            return 0.0  # No baseline yet
        
        baseline = self.baseline_profiles[user_id]
        
        # Compare current vs baseline
        anomaly_score = 0.0
        
        # Access time anomaly
        current_hour = datetime.now().hour
        typical_hours = baseline.get('typical_hours', [9, 10, 11, 14, 15, 16])
        if current_hour not in typical_hours:
            anomaly_score += 0.3
        
        # Access frequency anomaly
        current_freq = session_data.get('access_freq', 0)
        typical_freq = baseline.get('avg_access_freq', 10)
        if abs(current_freq - typical_freq) > typical_freq * 2:
            anomaly_score += 0.4
        
        # Resource pattern anomaly
        current_resources = set(session_data.get('accessed_resources', []))
        typical_resources = set(baseline.get('typical_resources', []))
        overlap = len(current_resources & typical_resources) / max(len(typical_resources), 1)
        if overlap < 0.5:
            anomaly_score += 0.3
        
        return min(anomaly_score, 1.0)
    
    def _time_anomaly(self, session_data: Dict) -> float:
        """Calculate time-of-day anomaly"""
        current_hour = datetime.now().hour
        
        # Business hours: 9-17
        if 9 <= current_hour <= 17:
            return 0.0
        elif 6 <= current_hour <= 21:
            return 0.3
        else:
            return 0.8  # Late night access is suspicious
    
    def _calculate_velocity(self, user_id: str, session_data: Dict) -> float:
        """Calculate geographic velocity (impossible travel detection)"""
        if user_id not in self.baseline_profiles:
            return 0.0
        
        last_location = self.baseline_profiles[user_id].get('last_location')
        current_location = session_data.get('location', (0, 0))
        last_time = self.baseline_profiles[user_id].get('last_login_time')
        
        if not last_location or not last_time:
            return 0.0
        
        # Calculate distance (simplified)
        distance_km = self._haversine_distance(last_location, current_location)
        
        # Calculate time difference
        time_diff_hours = (datetime.now() - last_time).total_seconds() / 3600
        
        if time_diff_hours > 0:
            velocity = distance_km / time_diff_hours
            # Flag if > 800 km/h (impossible for normal travel)
            return min(velocity / 800.0, 1.0)
        
        return 0.0
    
    def _haversine_distance(self, coord1: tuple, coord2: tuple) -> float:
        """Calculate distance between two coordinates"""
        lat1, lon1 = map(radians, coord1)
        lat2, lon2 = map(radians, coord2)
        
        dlat = lat2 - lat1
        dlon = lon2 - lon1
        
        a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
        c = 2 * asin(sqrt(a))
        
        return 6371 * c  # Earth radius in km
    
    async def _query_threat_intelligence(self, ip_address: str) -> float:
        """Query threat intelligence for IP reputation"""
        # Simplified - in production, query threat intel APIs
        # Return score 0.0 (safe) to 1.0 (malicious)
        
        # Simulate threat check
        known_bad_ips = ['10.0.0.666', '192.168.1.evil']
        
        if ip_address in known_bad_ips:
            return 0.9
        
        # Random threat score for demo
        return np.random.random() * 0.3
    
    def update_baseline(self, user_id: str, context: TrustContext):
        """Update user behavioral baseline"""
        if user_id not in self.baseline_profiles:
            self.baseline_profiles[user_id] = {
                'typical_hours': [],
                'typical_resources': [],
                'avg_access_freq': 0,
                'last_location': context.location_coordinates,
                'last_login_time': datetime.now()
            }
        
        profile = self.baseline_profiles[user_id]
        
        # Update typical hours
        current_hour = datetime.now().hour
        if current_hour not in profile['typical_hours']:
            profile['typical_hours'].append(current_hour)
        
        # Update location
        profile['last_location'] = context.location_coordinates
        profile['last_login_time'] = datetime.now()
        
        logger.info(f"Updated baseline for user {user_id}")


class TemporalPatternAnalyzer:
    """Analyze temporal patterns for anomaly detection"""
    
    def __init__(self, window_hours: int = 24):
        """初始化 Temporal Pattern Analyzer"""
        self.window_hours = window_hours
        self.temporal_patterns = {}
        logger.info("Temporal Pattern Analyzer initialized")
    
    def analyze_access_pattern(
        self,
        user_id: str,
        access_times: List[datetime]
    ) -> Dict[str, float]:
        """Analyze temporal access patterns"""
        if len(access_times) < 5:
            return {'temporal_anomaly': 0.0}
        
        # Convert to hours
        hours = [t.hour + t.minute/60.0 for t in access_times]
        
        # Calculate statistics
        mean_hour = np.mean(hours)
        std_hour = np.std(hours)
        
        # Day of week analysis
        weekdays = [t.weekday() for t in access_times]
        weekend_ratio = sum(1 for d in weekdays if d >= 5) / len(weekdays)
        
        # Time clustering (are accesses clustered or scattered?)
        sorted_hours = sorted(hours)
        gaps = np.diff(sorted_hours)
        max_gap = np.max(gaps) if len(gaps) > 0 else 0
        
        # Anomaly score
        anomaly_score = 0.0
        
        # High weekend access is suspicious
        if weekend_ratio > 0.5:
            anomaly_score += 0.3
        
        # Large gaps in access time
        if max_gap > 12:
            anomaly_score += 0.3
        
        # Off-hours access (late night)
        late_night_count = sum(1 for h in hours if h < 6 or h > 22)
        if late_night_count / len(hours) > 0.3:
            anomaly_score += 0.4
        
        return {
            'temporal_anomaly': min(anomaly_score, 1.0),
            'mean_access_hour': mean_hour,
            'weekend_ratio': weekend_ratio,
            'max_gap_hours': max_gap
        }
    
    def detect_burst_activity(
        self,
        access_times: List[datetime],
        threshold_per_hour: int = 50
    ) -> bool:
        """Detect burst/flood activity"""
        if len(access_times) < 2:
            return False
        
        # Count accesses in sliding 1-hour windows
        for i, base_time in enumerate(access_times[:-1]):
            window_end = base_time + timedelta(hours=1)
            count = sum(1 for t in access_times[i:] if t <= window_end)
            
            if count > threshold_per_hour:
                logger.warning(f"Burst detected: {count} accesses in 1 hour")
                return True
        
        return False


# 示例使用
async def main():
    """主函數"""
    logger.info("=== Pandora Zero Trust Context Engine 測試 ===")
    
    aggregator = ZeroTrustContextAggregator()
    
    # 測試正常會話
    session_data = {
        'password_auth': True,
        'mfa_enabled': True,
        'credential_age': 2,
        'network_type': 'corporate',
        'location': (40.7128, -74.0060),  # New York
        'device_trust': 0.9,
        'access_freq': 20,
        'roles': ['user', 'developer'],
        'source_ip': '192.168.1.100'
    }
    
    context = await aggregator.collect_context('user123', 'device456', session_data)
    
    print(f"\n信任上下文:")
    print(f"  用戶: {context.user_id}")
    print(f"  認證強度: {context.authentication_strength:.2f}")
    print(f"  設備姿態分數: {context.device_posture_score:.2f}")
    print(f"  行為異常: {context.access_pattern_anomaly:.2f}")
    print(f"  威脅情報分數: {context.threat_intelligence_score:.2f}")
    
    # 更新基線
    aggregator.update_baseline('user123', context)
    
    logger.info("\n=== 測試完成 ===")


if __name__ == "__main__":
    import asyncio
    asyncio.run(main())

