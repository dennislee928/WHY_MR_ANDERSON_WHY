.PHONY: help up down logs scan-nuclei scan-nmap clean backup

help:
	@echo "Security Stack Commands:"
	@echo "  make up           - Start all services"
	@echo "  make down         - Stop all services"
	@echo "  make logs         - View logs"
	@echo "  make scan-nuclei  - Run Nuclei scan"
	@echo "  make scan-nmap    - Run Nmap scan"
	@echo "  make backup       - Backup databases"
	@echo "  make clean        - Remove all volumes"

up:
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 10
	@docker-compose ps

down:
	docker-compose down

logs:
	docker-compose logs -f

scan-nuclei:
	docker-compose run --rm scanner-nuclei \
		nuclei -u $(TARGET) -o /results/nuclei-$(shell date +%Y%m%d-%H%M%S).json

scan-nmap:
	docker-compose run --rm nmap \
		nmap $(TARGET) -oX /results/nmap-$(shell date +%Y%m%d-%H%M%S).xml

backup:
	@mkdir -p backups
	docker-compose exec -T postgres pg_dump -U sectools security > backups/db-$(shell date +%Y%m%d-%H%M%S).sql
	@echo "Backup created in backups/"

clean:
	docker-compose down -v
	@echo "All volumes removed"

restart:
	docker-compose restart

ps:
	docker-compose ps

health:
	@docker-compose ps --format "table {{.Service}}\t{{.Status}}"
