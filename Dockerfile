# Frontend Build
FROM node:22-alpine AS frontend
RUN corepack enable
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ .
RUN pnpm build

# Go Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o picourl ./cmd/server

# Runtime 
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/picourl /usr/local/bin/picourl
COPY --from=builder /app/frontend/dist /static
COPY --from=builder /app/internal/store/migrations /migrations
EXPOSE 8080
CMD [ "picourl" ]
