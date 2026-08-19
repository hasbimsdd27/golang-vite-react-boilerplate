# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-build
RUN corepack enable && corepack prepare pnpm@9.15.0 --activate
WORKDIR /app/frontend
COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY frontend/ ./
RUN pnpm build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 3: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-build /app/backend/server .
COPY --from=frontend-build /app/frontend/dist ./frontend/dist
ENV STATIC_DIR=./frontend/dist
EXPOSE 8080
CMD ["./server"]
