# 知猿 (Zioran) 运维 Makefile
# 服务器: 120.26.192.163 | 域名: zioran.com

SSH = sshpass -p 'Admin66668888' ssh -o StrictHostKeyChecking=no root@120.26.192.163
SCP = sshpass -p 'Admin66668888' scp -o StrictHostKeyChecking=no
REMOTE_DIR = /opt/zioran

# ==================== 本地开发 ====================

.PHONY: dev dev-backend dev-frontend dev-admin build test

## 启动本地后端
dev-backend:
	cd backend && go run cmd/server/main.go

## 启动本地前端
dev-frontend:
	cd frontend && pnpm dev

## 启动本地管理端
dev-admin:
	cd admin && pnpm dev

## 本地全部启动（需要3个终端，此处仅构建验证）
dev: build
	@echo "请分别在3个终端执行: make dev-backend / make dev-frontend / make dev-admin"

## 构建全部
build: build-backend build-frontend build-admin

build-backend:
	cd backend && go build -o ../bin/zioran-api ./cmd/server/
	@echo "✅ Backend built"

build-frontend:
	cd frontend && pnpm build
	@echo "✅ Frontend built"

build-admin:
	cd admin && pnpm build
	@echo "✅ Admin built"

## 运行测试
test:
	cd backend && go test ./...
	@echo "✅ All tests passed"

# ==================== 服务器部署 ====================

.PHONY: deploy deploy-backend deploy-frontend deploy-admin deploy-full

## 一键全量部署
deploy-full: deploy-code deploy-migrate deploy-backend deploy-frontend deploy-admin deploy-restart
	@echo "✅ 全量部署完成"

## 推送代码到服务器
deploy-code:
	$(SSH) "cd $(REMOTE_DIR) && git pull origin main"

## 执行数据库迁移
deploy-migrate:
	$(SSH) "cd $(REMOTE_DIR) && for f in backend/migrations/*.sql; do mysql -uroot -proot123456 zioran < \$$f 2>/dev/null; done"
	@echo "✅ Migration done"

## 部署后端
deploy-backend:
	$(SSH) "export PATH=\$$PATH:/usr/local/go/bin && cd $(REMOTE_DIR)/backend && go build -o $(REMOTE_DIR)/bin/zioran-api ./cmd/server/ && cp config.yaml $(REMOTE_DIR)/bin/config.yaml"
	@echo "✅ Backend deployed"

## 部署前端
deploy-frontend:
	$(SSH) "cd $(REMOTE_DIR)/frontend && pnpm install --frozen-lockfile 2>/dev/null && NEXT_PUBLIC_API_URL=https://api.zioran.com/api/v1 pnpm build"
	@echo "✅ Frontend deployed"

## 部署管理端
deploy-admin:
	$(SSH) "cd $(REMOTE_DIR)/admin && pnpm install --frozen-lockfile 2>/dev/null && pnpm build"
	@echo "✅ Admin deployed"

## 重启所有服务
deploy-restart:
	$(SSH) "pm2 restart zioran-api && pm2 restart zioran-frontend && pm2 save"
	@echo "✅ Services restarted"

# ==================== 服务器管理 ====================

.PHONY: ssh status logs logs-api logs-frontend start stop restart

## SSH连接服务器
ssh:
	$(SSH)

## 查看服务状态
status:
	$(SSH) "pm2 list"

## 查看全部日志
logs:
	$(SSH) "pm2 logs --lines 30 --nostream"

## 查看后端日志
logs-api:
	$(SSH) "pm2 logs zioran-api --lines 50 --nostream"

## 查看前端日志
logs-frontend:
	$(SSH) "pm2 logs zioran-frontend --lines 30 --nostream"

## 启动所有服务
start:
	$(SSH) "pm2 start zioran-api zioran-frontend"

## 停止所有服务
stop:
	$(SSH) "pm2 stop zioran-api zioran-frontend"

## 重启所有服务
restart:
	$(SSH) "pm2 restart zioran-api zioran-frontend && pm2 save"
	@echo "✅ Restarted"

# ==================== 数据库 ====================

.PHONY: db-shell db-backup db-restore seed

## 进入MySQL命令行
db-shell:
	$(SSH) "mysql -uroot -proot123456 zioran"

## 数据库备份
db-backup:
	$(SSH) "mysqldump -uroot -proot123456 zioran > /opt/zioran/backup/zioran_$$(date +%Y%m%d_%H%M%S).sql && echo '✅ Backup done'"

## 初始化种子数据
seed:
	$(SSH) "mysql -uroot -proot123456 zioran < $(REMOTE_DIR)/backend/migrations/004_seed_data.sql"
	@echo "✅ Seed data loaded"

# ==================== 快捷 ====================

.PHONY: push quick-deploy health

## Git提交并推送
push:
	git add -A && git commit -m "$(msg)" && git push origin main

## 快速部署（拉代码+重启，不重新构建前端）
quick-deploy:
	$(SSH) "cd $(REMOTE_DIR) && git pull origin main && export PATH=\$$PATH:/usr/local/go/bin && cd backend && go build -o $(REMOTE_DIR)/bin/zioran-api ./cmd/server/ && cp config.yaml $(REMOTE_DIR)/bin/config.yaml && pm2 restart zioran-api && pm2 save"
	@echo "✅ Quick deploy done"

## 健康检查
health:
	@echo "API:" && curl -s http://120.26.192.163:8080/api/v1/categories | python3 -c "import sys,json;d=json.loads(sys.stdin.read());print('  ✅ OK' if d['code']==0 else '  ❌ FAIL')"
	@echo "Frontend:" && curl -s -o /dev/null -w "  HTTP %{http_code}\n" http://120.26.192.163:3000/
	@echo "Nginx:" && curl -s -o /dev/null -w "  HTTP %{http_code}\n" http://120.26.192.163/api/v1/categories

## 显示帮助
help:
	@echo ""
	@echo "知猿 (Zioran) 运维命令"
	@echo "=========================="
	@echo ""
	@echo "本地开发:"
	@echo "  make dev-backend    启动后端 (localhost:8080)"
	@echo "  make dev-frontend   启动前端 (localhost:3000)"
	@echo "  make dev-admin      启动管理端 (localhost:5173)"
	@echo "  make build          构建全部"
	@echo "  make test           运行测试"
	@echo ""
	@echo "部署:"
	@echo "  make deploy-full    一键全量部署"
	@echo "  make quick-deploy   快速部署(仅后端)"
	@echo "  make deploy-restart 重启服务"
	@echo ""
	@echo "运维:"
	@echo "  make ssh            SSH连接服务器"
	@echo "  make status         查看服务状态"
	@echo "  make logs           查看日志"
	@echo "  make restart        重启服务"
	@echo "  make health         健康检查"
	@echo ""
	@echo "数据库:"
	@echo "  make db-shell       MySQL命令行"
	@echo "  make db-backup      数据库备份"
	@echo "  make seed           初始化种子数据"
	@echo ""
	@echo "Git:"
	@echo "  make push msg=\"提交信息\"  提交并推送"
	@echo ""
