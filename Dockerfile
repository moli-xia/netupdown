FROM node:24-alpine AS admin
WORKDIR /src/web/admin
COPY web/admin/package.json web/admin/pnpm-lock.yaml web/admin/pnpm-workspace.yaml ./
RUN corepack enable && pnpm install --frozen-lockfile
COPY web/admin ./
RUN pnpm build

FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=admin /src/web/admin/dist internal/assets/admin
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o /out/netupdown ./cmd/netupdown

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1000 netupdown
USER netupdown
WORKDIR /app
COPY --from=build /out/netupdown .
VOLUME ["/app/data"]
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/healthz || exit 1
ENTRYPOINT ["./netupdown"]
CMD ["serve", "-c", "/app/config.yaml"]
