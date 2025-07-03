# -------- 建構階段 --------
    FROM golang:1.24-alpine AS builder

    ENV CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64
    
    WORKDIR /app/member
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY api/member/post ./api/member/post
    COPY common ./common
    
    RUN go build -o post ./api/member/post/cmd/main.go
    
    # -------- 執行階段 --------
    FROM gcr.io/distroless/static-debian12
    
    WORKDIR /
    
    COPY --from=builder /app/member/post .

    EXPOSE 8080
    
    ENTRYPOINT ["/post"]