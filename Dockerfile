# 1단계: 빌드 환경
FROM golang:1.26-alpine AS builder
WORKDIR /app

# 모듈 다운로드 (레이어 캐시 활용)
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 소스 코드 복사 및 빌드
# -w -s: 디버그 심볼 제거로 바이너리 크기 최소화
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /collector ./cmd/collector

# 2단계: 실행 환경
FROM alpine:3.21
WORKDIR /app

# TLS 통신을 위한 CA 인증서 설치
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S collector && \
    adduser -S -G collector collector

# 실행 파일 복사
COPY --from=builder /collector .

# 비루트 유저로 실행
USER collector

ENTRYPOINT ["./collector"]
