# Nginx Dockerfile for Render deployment
FROM nginx:1.29.2-alpine3.22-slim

# 安裝必要工具
RUN apk add --no-cache curl

# 複製 Nginx 配置
COPY configs/nginx/nginx-paas.conf /etc/nginx/nginx.conf
COPY configs/nginx/default-paas.conf /etc/nginx/conf.d/default.conf

# 複製靜態檔案
COPY web/static /usr/share/nginx/html

# 建立健康檢查端點
RUN echo "OK" > /usr/share/nginx/html/health

# 暴露端口
EXPOSE 80

# 創建非 root 用戶
RUN addgroup -g 101 -S nginx && \
    adduser -S -D -H -u 101 -h /var/cache/nginx -s /sbin/nologin -G nginx -g nginx nginx

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/health || exit 1

# 切換到非 root 用戶
USER nginx

# 啟動 Nginx
CMD ["nginx", "-g", "daemon off;"]

