# -------- 建構階段 --------
    FROM golang:1.24-alpine AS builder

    ENV CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64
    
    WORKDIR /app/admin
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY api/admin/category ./api/admin/category
    COPY common ./common
    
    RUN go build -o category ./api/admin/category/cmd/main.go
    
    # -------- 執行階段 --------
    FROM gcr.io/distroless/static-debian12
    
    WORKDIR /
    
    COPY --from=builder /app/admin/category .
    
    EXPOSE 8080
    
    ENTRYPOINT ["/category"]
    