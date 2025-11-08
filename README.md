# go-ad-admin

Lightweight Admin UI for Samba AD DC and Kea DHCP, written in Go.  
No Node/npm toolchain required (server-side rendered templates).

## Quick start

```bash
go mod tidy
go test ./... -cover -covermode=atomic
go run ./cmd/go-ad-admin
# -> http://localhost:8080
```

## Configuration (env)

| VAR               | Default | Description |
|-------------------|---------|-------------|
| GO_AD_LISTEN      | :8080 | Bind address |
| GO_AD_ENV         | dev | dev/prod |
| GO_AD_SESSION_KEY | (random at start) | 32+ bytes recommended |
| GO_AD_LDAP_URL    | ldap://127.0.0.1:389 | LDAP/LDAPS URL |
| GO_AD_LDAP_BASEDN | dc=example,dc=com | Base DN |
| GO_AD_PRIVACY     | low | low/high (pseudonymize listings) |

## Layout

- `cmd/go-ad-admin` – main entry
- `internal/config` – env config & validation
- `internal/web` – HTTP handlers (SSR templates)
- `internal/ldap` – interface & mocks (to be implemented)
- `internal/audit` – append-only JSONL audit log
- `internal/kea` – Kea HTTP client (to be implemented)
- `web/templates` – Go `html/template` files
- `docs` – Requirements & Use Cases
- `features` – BDD Gherkin features
