param([string]$Version = "dev")
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Push-Location (Join-Path $root "web/admin")
try {
  pnpm install --frozen-lockfile
  pnpm build
} finally { Pop-Location }
Copy-Item -Path (Join-Path $root "web/admin/dist/*") -Destination (Join-Path $root "internal/assets/admin") -Recurse -Force
New-Item -ItemType Directory -Force (Join-Path $root "bin") | Out-Null
go build -trimpath -ldflags "-s -w -X main.version=$Version" -o (Join-Path $root "bin/netupdown.exe") ./cmd/netupdown
