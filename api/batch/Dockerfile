# -------- 建構階段 --------
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# 複製 go.mod/go.sum 並安裝依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製需要的程式碼
COPY api/batch ./api/batch
COPY common ./common

# 編譯 main.go（🔥 這是你的 batch 入口點）
RUN go build -o batch ./api/batch/cmd/main.go

# -------- 執行階段 --------
FROM gcr.io/distroless/static-debian12

WORKDIR /

# 複製執行檔
COPY --from=builder /app/batch .

EXPOSE 8080

# 啟動 batch server（可由 Cloud Run 或 Cloud Scheduler 呼叫）
ENTRYPOINT ["/batch"]