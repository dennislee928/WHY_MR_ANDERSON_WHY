部署策略藍圖
微服務 (Microservice)	建議部署平台	部署方式	理由 (Why here?)
postgres (資料庫)	Railway.app	Managed Database	Railway 提供非常易於使用且免費的託管式 PostgreSQL。將其獨立出來，可以為您的核心數據提供一個穩定的基礎。
redis (快取)	Render	Managed Redis	Render 同樣提供免費的託管式 Redis，將它與資料庫分開部署，可以分散風險並利用不同平台的資源。
pandora-agent (核心後端)	Koyeb	Web Service (Docker)	這是您最重要的服務。Koyeb 的免費方案提供 2 個永不休眠 (Always-on) 的 Nano 容器，非常適合需要 24/7 全時運行的核心後端 API 服務。
axiom-ui (UI 伺服器)	Patr.io	Container Deployment	Patr.io 同樣提供慷慨的免費容器託管方案。將 UI 服務部署在這裡，可以體驗另一個平台的部署流程，並將使用者介面的流量與後端服務分離。
nginx (反向代理)	Render	Web Service	Render 的 Web Service 非常適合部署像 Nginx 這樣的反向代理。您可以利用其路由規則，將來自外部的請求導向部署在 Patr.io 上的 axiom-ui 和 Fly.io 上的 grafana。
prometheus (監控指標)	Fly.io	App (Docker)	Prometheus 需要永久儲存空間 (Persistent Volume) 來存放時間序列數據。Fly.io 是少數在免費方案中提供數 GB 永久儲存空間的平台，是部署 Prometheus 的不二之選。
loki (日誌聚合)	Fly.io	App (Docker)	與 Prometheus 相同，Loki 也需要永久儲存空間。將它們放在同一個平台可以簡化網路設定。
grafana (視覺化)	Fly.io	App (Docker)	Grafana 需要頻繁地與 Prometheus 和 Loki 通訊，將它們部署在同一個平台內部可以大幅降低網路延遲。同時，Fly.io 的不休眠特性也能確保您的儀表板隨時可用。
alertmanager (告警)	Fly.io	App (Docker)	告警系統必須全時運行。將整個監控與告警系統 (Prometheus, Loki, Grafana, Alertmanager) 集中部署在 Fly.io 是最合理且高效的選擇。
promtail (日誌收集)	Koyeb	Web Service (Sidecar)	Promtail 應與它要監控的應用程式部署在一起。因此，將它與核心後端 pandora-agent 一同部署在 Koyeb 上，並設定它將日誌透過公網 URL 推送到 Fly.io 上的 Loki。
node-exporter (系統指標)	不部署	N/A	在 PaaS 環境下，您無法也無需監控底層主機，這個服務可以直接捨棄。