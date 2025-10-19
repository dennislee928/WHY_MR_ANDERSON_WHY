/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  
  // 輸出配置（用於靜態導出）
  output: 'export',
  trailingSlash: true,
  images: {
    unoptimized: true,
  },
  
  // 環境變數
  env: {
    NEXT_PUBLIC_APP_NAME: 'Pandora Box Console IDS-IPS',
    NEXT_PUBLIC_APP_VERSION: process.env.npm_package_version || '1.0.0',
  },
  
  // Webpack 配置
  webpack: (config, { isServer }) => {
    // 修復某些套件的問題
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
      };
    }
    
    return config;
  },
  
  // 實驗性功能
  experimental: {
    // 啟用服務器動作
    serverActions: true,
  },
  
  // 編譯時忽略的路徑
  eslint: {
    ignoreDuringBuilds: false,
  },
  
  typescript: {
    ignoreBuildErrors: false,
  },
};

module.exports = nextConfig;

