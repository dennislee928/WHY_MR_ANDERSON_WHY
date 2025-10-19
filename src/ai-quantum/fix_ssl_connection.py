#!/usr/bin/env python3
"""
修復 IBM Quantum SSL 連接問題
"""
import os
import sys
import ssl
import certifi
import urllib3

print("="*60)
print("SSL 連接診斷與修復")
print("="*60)

# 方案 1: 禁用 SSL 驗證（僅用於測試）
def test_without_ssl_verify():
    """測試不驗證 SSL（不推薦用於生產環境）"""
    print("\n[方案 1] 測試不驗證 SSL...")
    
    try:
        # 禁用 SSL 警告
        urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)
        
        from qiskit_ibm_runtime import QiskitRuntimeService
        
        token = os.getenv('IBM_QUANTUM_TOKEN')
        if not token:
            print("❌ Token 未設定")
            return False
        
        # 設定環境變數禁用 SSL 驗證
        os.environ['CURL_CA_BUNDLE'] = ''
        os.environ['REQUESTS_CA_BUNDLE'] = ''
        
        print("正在連接（跳過 SSL 驗證）...")
        service = QiskitRuntimeService(
            channel='ibm_quantum',
            token=token,
            verify=False  # 跳過 SSL 驗證
        )
        
        print("✅ 連接成功！")
        backends = service.backends()
        print(f"可用後端: {len(backends)} 個")
        return True
        
    except Exception as e:
        print(f"❌ 失敗: {type(e).__name__}")
        print(f"   {str(e)[:200]}")
        return False

# 方案 2: 更新 CA 證書
def test_with_certifi():
    """使用 certifi 提供的證書庫"""
    print("\n[方案 2] 使用 certifi 證書...")
    
    try:
        cert_path = certifi.where()
        print(f"證書路徑: {cert_path}")
        
        os.environ['REQUESTS_CA_BUNDLE'] = cert_path
        os.environ['CURL_CA_BUNDLE'] = cert_path
        
        from qiskit_ibm_runtime import QiskitRuntimeService
        
        token = os.getenv('IBM_QUANTUM_TOKEN')
        if not token:
            print("❌ Token 未設定")
            return False
        
        print("正在連接（使用 certifi 證書）...")
        service = QiskitRuntimeService(
            channel='ibm_quantum',
            token=token
        )
        
        print("✅ 連接成功！")
        backends = service.backends()
        print(f"可用後端: {len(backends)} 個")
        return True
        
    except Exception as e:
        print(f"❌ 失敗: {type(e).__name__}")
        print(f"   {str(e)[:200]}")
        return False

# 方案 3: 使用系統 CA 證書
def test_with_system_ca():
    """使用系統 CA 證書"""
    print("\n[方案 3] 使用系統 CA 證書...")
    
    try:
        # 可能的系統證書路徑
        ca_paths = [
            '/etc/ssl/certs/ca-certificates.crt',
            '/etc/pki/tls/certs/ca-bundle.crt',
            '/etc/ssl/ca-bundle.pem',
            '/usr/local/share/certs/ca-root-nss.crt',
        ]
        
        ca_bundle = None
        for path in ca_paths:
            if os.path.exists(path):
                ca_bundle = path
                print(f"找到系統證書: {path}")
                break
        
        if ca_bundle:
            os.environ['REQUESTS_CA_BUNDLE'] = ca_bundle
            os.environ['CURL_CA_BUNDLE'] = ca_bundle
        
        from qiskit_ibm_runtime import QiskitRuntimeService
        
        token = os.getenv('IBM_QUANTUM_TOKEN')
        if not token:
            print("❌ Token 未設定")
            return False
        
        print("正在連接（使用系統證書）...")
        service = QiskitRuntimeService(
            channel='ibm_quantum',
            token=token
        )
        
        print("✅ 連接成功！")
        backends = service.backends()
        print(f"可用後端: {len(backends)} 個")
        return True
        
    except Exception as e:
        print(f"❌ 失敗: {type(e).__name__}")
        print(f"   {str(e)[:200]}")
        return False

# 方案 4: 檢查網路連接
def test_direct_connection():
    """測試直接 HTTPS 連接"""
    print("\n[方案 4] 測試直接 HTTPS 連接...")
    
    try:
        import requests
        
        print("測試連接到 IBM Quantum API...")
        response = requests.get(
            'https://auth.quantum-computing.ibm.com/api/version',
            timeout=10
        )
        
        print(f"✅ HTTP 狀態碼: {response.status_code}")
        print(f"回應內容: {response.text[:100]}")
        return True
        
    except Exception as e:
        print(f"❌ 失敗: {type(e).__name__}")
        print(f"   {str(e)[:200]}")
        return False

# 方案 5: 使用代理（如果需要）
def test_with_proxy():
    """測試使用代理連接"""
    print("\n[方案 5] 檢查代理設定...")
    
    http_proxy = os.getenv('HTTP_PROXY') or os.getenv('http_proxy')
    https_proxy = os.getenv('HTTPS_PROXY') or os.getenv('https_proxy')
    
    if http_proxy or https_proxy:
        print(f"HTTP 代理: {http_proxy or 'None'}")
        print(f"HTTPS 代理: {https_proxy or 'None'}")
    else:
        print("ℹ️  未設定代理")
    
    return False

def main():
    """主函數"""
    print("\n開始診斷 SSL 連接問題...\n")
    
    # 顯示環境資訊
    print("環境資訊:")
    print(f"  Python 版本: {sys.version.split()[0]}")
    print(f"  SSL 版本: {ssl.OPENSSL_VERSION}")
    print(f"  certifi 路徑: {certifi.where()}")
    
    token = os.getenv('IBM_QUANTUM_TOKEN')
    print(f"  IBM Token: {'已設定' if token else '未設定'}")
    print("")
    
    # 依序測試各種方案
    results = []
    
    # 測試 1: 直接連接
    results.append(("直接 HTTPS", test_direct_connection()))
    
    # 測試 2: certifi
    results.append(("certifi 證書", test_with_certifi()))
    
    # 測試 3: 系統證書
    results.append(("系統證書", test_with_system_ca()))
    
    # 測試 4: 不驗證 SSL
    results.append(("跳過 SSL 驗證", test_without_ssl_verify()))
    
    # 測試 5: 代理
    test_with_proxy()
    
    # 總結
    print("\n" + "="*60)
    print("診斷總結")
    print("="*60)
    
    for method, success in results:
        status = "✅ 成功" if success else "❌ 失敗"
        print(f"  {method:<20} : {status}")
    
    # 如果所有方案都失敗
    if not any(success for _, success in results):
        print("\n⚠️  所有連接方案都失敗")
        print("\n可能的原因:")
        print("  1. 企業防火牆/代理阻擋 HTTPS")
        print("  2. IBM Quantum 服務暫時不可用")
        print("  3. 網路環境限制（如 VPN、內網）")
        print("  4. SSL/TLS 協議版本不相容")
        print("\n建議:")
        print("  ✅ 使用本地 Aer 模擬器（已驗證可用）")
        print("  ✅ 檢查網路連接和防火牆設定")
        print("  ✅ 聯繫網路管理員確認出站 HTTPS 規則")
    else:
        successful = [method for method, success in results if success]
        print(f"\n✅ 成功的方案: {', '.join(successful)}")

if __name__ == "__main__":
    main()

