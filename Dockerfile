# ----------------------------
# Frontend
# ----------------------------
FROM node:18-alpine AS frontend-build
WORKDIR /src/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build


# ----------------------------
# Backend build
# ----------------------------
FROM golang:1.24-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates

COPY backend/go.mod backend/go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

COPY backend/ ./

COPY --from=frontend-build /src/frontend/dist ./frontend/dist

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /app/payroll ./main.go


# ----------------------------
# Runtime
# ----------------------------
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    fontconfig \
    libfreetype6 \
    libjpeg62-turbo \
    libpng16-16 \
    libssl3 \
    libx11-6 \
    libxcb1 \
    libxext6 \
    libxrender1 \
    wget \
    xfonts-75dpi \
    xfonts-base \
    && wget -q https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-3/wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && dpkg -i wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && rm wkhtmltox_0.12.6.1-3.bookworm_amd64.deb \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/payroll .
COPY --from=builder /src/frontend/dist ./frontend/dist
COPY backend/templates ./templates

RUN mkdir -p /data

EXPOSE 8080

ENV PORT=8080
ENV FRONTEND_DIST=/app/frontend/dist
ENV DB_PATH=/data/payroll.db

CMD ["/app/payroll"]