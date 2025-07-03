# === Member 前台服務 ===
run-member-post:
	go run ./api/member/post/cmd

run-member-category:
	go run ./api/member/category/cmd

run-member-apigw:
	go run ./api/member/apigw/cmd

run-member-all:
	@echo "🚀 啟動 member 所有服務中..."
	@make run-member-post & \
	make run-member-category & \
	make run-member-apigw

# === Admin 後台服務 ===
run-admin-post:
	go run ./api/admin/post/cmd

run-admin-category:
	go run ./api/admin/category/cmd

run-admin-apigw:
	go run ./api/admin/apigw/cmd

run-admin-all:
	@echo "🚀 啟動 admin 所有服務中..."
	@make run-admin-post & \
	make run-admin-category & \
	make run-admin-apigw

# === 通用指令 ===
run-all:
	@echo "🚀 同時啟動 member + admin 所有服務中..."
	@make run-member-all & \
	make run-admin-all

stop:
	@echo "🛑 請手動停止所有背景服務"
	@echo "👉 包含 member-xxx 與 admin-xxx"
