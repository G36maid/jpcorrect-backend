# AGENTS.md

## Commands
- **Build:** `go build ./...`
- **Lint:** `golangci-lint run`
- **Test All:** `go test ./...`
- **Test Single:** `go test -v ./path/to/pkg -run TestName`
- **Database:** Run `make migrate-up` for schema changes. GORM is used for data access (no sqlc).
- **Dev Server:** `make air` (hot reload).

## Code Style
- **Formatting:** Run `gofmt` before committing.
- **Naming:** PascalCase for exported, camelCase for unexported identifiers.
- **Imports:** Group: 1. Standard Lib, 2. Third-party, 3. Internal (e.g. `jpcorrect-backend/internal/...`).
- **Errors:** Handle errors explicitly (`if err != nil`). Do not ignore errors.
- **Architecture:** `cmd/` for entrypoints, `internal/` for logic. `api` layer handles HTTP, `domain` defines types/interfaces, `repository` implements storage using GORM.
- **Libraries:** Uses Gin (HTTP) and GORM (ORM/Postgres). Do not introduce new frameworks without permission.

## GORM Usage
- **Models:** All domain models in `internal/domain/` use GORM struct tags (e.g., `gorm:"column:user_id;primaryKey"`).
- **Table Names:** Each model has a `TableName()` method to specify schema-qualified names (e.g., `jpcorrect.user`).
- **Repositories:** Use GORM query builders (`db.Where()`, `db.First()`, `db.Create()`, etc.) instead of raw SQL.
- **Error Mapping:** `gorm.ErrRecordNotFound` → `domain.ErrNotFound`.
- **Context:** Always use `db.WithContext(ctx)` for proper context propagation.
- **Migrations:** Continue using `golang-migrate` for schema changes. Do NOT use GORM's AutoMigrate in production.
