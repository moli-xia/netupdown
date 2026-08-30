VERSION ?= dev

.PHONY: admin build test
admin:
	cd web/admin && pnpm install --frozen-lockfile && pnpm build
	rm -rf internal/assets/admin
	cp -r web/admin/dist internal/assets/admin

build:
	go build -trimpath -ldflags "-s -w -X main.version=$(VERSION)" -o bin/netupdown ./cmd/netupdown

test:
	go test ./...
