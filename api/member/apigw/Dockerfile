# -------- 第一階段：建置階段 --------
# 使用官方 Go 映像，Alpine 版體積較小
FROM golang:1.24-alpine AS builder

# 設定環境變數，讓 Go modules 可在 container 中使用
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# ✅ 專案模組根目錄
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# ✅ 拷貝所有需要編譯的程式碼（包含 middleware）
COPY api ./api
COPY common ./common

# ✅ 用資料夾編譯而不是單一 .go 檔
RUN go build -o apigw ./api/member/apigw/cmd

# -------- 第二階段：運行階段 --------
# 使用更精簡的容器映像：distroless 只有 libc，不含 shell
FROM gcr.io/distroless/static-debian12

# 設定容器啟動時的工作目錄（可選）
WORKDIR /

# 從上一階段拷貝已編譯完成的 binary
COPY --from=builder /app/apigw .

EXPOSE 8080

# 設定容器啟動後執行的指令
ENTRYPOINT ["/apigw"]