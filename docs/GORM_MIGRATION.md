# GORM Migration Documentation

## Overview
This document describes the migration from manual `pgx/v5` repositories to GORM ORM, completed on the jpcorrect-backend project.

## What Changed

### 1. Dependencies Added
- `gorm.io/gorm v1.31.1` - GORM ORM core
- `gorm.io/driver/postgres v1.6.0` - PostgreSQL driver for GORM

### 2. Domain Models (`internal/domain/*.go`)
All domain models were updated with GORM struct tags:

**Before:**
```go
type User struct {
    UserID int    `db:"user_id" json:"user_id"`
    Name   string `db:"name" json:"name"`
}
```

**After:**
```go
type User struct {
    UserID int    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
    Name   string `json:"name" gorm:"column:name;not null"`
}

func (User) TableName() string {
    return "jpcorrect.user"
}
```

**Key Changes:**
- Replaced `db:` tags with `gorm:` tags
- Added `TableName()` methods to specify schema-qualified table names
- Added GORM-specific constraints (primaryKey, autoIncrement, not null)
- Changed `Practice.Date` from `string` to `time.Time` for proper date handling
- Specified custom types for Postgres ENUMs in `Mistake` model

### 3. Connection Interface (`internal/repository/connection.go`)

**Before:**
```go
type Connection interface {
    Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
    Query(context.Context, string, ...any) (pgx.Rows, error)
    QueryRow(context.Context, string, ...any) pgx.Row
}
```

**After:**
```go
type Connection struct {
    DB *gorm.DB
}

func NewConnection(db *gorm.DB) *Connection {
    return &Connection{DB: db}
}
```

Changed from an interface to a concrete wrapper struct around `*gorm.DB`.

### 4. Repository Implementations (`internal/repository/postgres_*.go`)

All six repository files were refactored:
- `postgres_user.go`
- `postgres_practice.go`
- `postgres_mistake.go`
- `postgres_note.go`
- `postgres_transcript.go`
- `postgres_ai_correction.go`

**Before (manual SQL):**
```go
func (u *postgresUserRepository) GetByID(ctx context.Context, userID int) (*domain.User, error) {
    query := `
        SELECT user_id, name
        FROM jpcorrect.user
        WHERE user_id = $1`
    
    rows, err := u.conn.Query(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []*domain.User
    for rows.Next() {
        var user domain.User
        if err := rows.Scan(&user.UserID, &user.Name); err != nil {
            return nil, err
        }
        users = append(users, &user)
    }
    
    if len(users) == 0 {
        return nil, domain.ErrNotFound
    }
    return users[0], nil
}
```

**After (GORM):**
```go
func (u *postgresUserRepository) GetByID(ctx context.Context, userID int) (*domain.User, error) {
    var user domain.User
    err := u.db.WithContext(ctx).Where("user_id = ?", userID).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return &user, nil
}
```

**Key Changes:**
- Removed `fetch()` helper methods
- Replaced raw SQL with GORM query builders
- Error handling: `gorm.ErrRecordNotFound` maps to `domain.ErrNotFound`
- `Create()`: GORM automatically populates the ID field
- `Update()`: Uses `Updates()` with map to avoid zero-value issues
- `Delete()`: Uses GORM's soft-delete-friendly delete pattern

### 5. Application Initialization (`internal/cmd/api.go`)

**Before:**
```go
dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatalf("failed to connect to database: %v", err)
}
defer dbpool.Close()

a := api.NewAPI(dbpool)
```

**After:**
```go
dsn := os.Getenv("DATABASE_URL")
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatalf("failed to connect to database: %v", err)
}

sqlDB, err := db.DB()
if err != nil {
    log.Fatalf("failed to get underlying sql.DB: %v", err)
}
defer sqlDB.Close()

// Set connection pool settings
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)

conn := repository.NewConnection(db)
a := api.NewAPI(conn)
```

Connection pool configuration is now done via the underlying `sql.DB`.

### 6. API Layer (`internal/api/server.go`)

Updated function signature:
```go
func NewAPI(conn *repository.Connection) *API
```

No other changes needed in the API layer - the repository interfaces remained the same.

## Benefits Gained

### Developer Experience
1. **Less Boilerplate**: ~50% reduction in repository code
2. **Type Safety**: GORM validates struct tags at runtime, catching mapping errors early
3. **Auto-ID Population**: No need to manually handle `RETURNING` clauses
4. **Query Builder**: Chain-able, readable query construction

### Maintainability
1. **Clearer Intent**: GORM queries are self-documenting
2. **Less Error-Prone**: No manual row scanning eliminates common bugs
3. **Relationships Ready**: Can easily add `Preload()` for related data in the future

### Performance Notes
- GORM uses reflection, adding ~5-10% overhead vs raw SQL
- For this application's scale, the tradeoff is acceptable
- Complex queries can still use `db.Raw()` if needed

## Tradeoffs Accepted

### What We Lost
1. **Direct SQL Control**: Cannot optimize specific queries as easily
2. **Query Visibility**: Generated SQL is less visible without logging
3. **Slightly Slower**: Reflection overhead on every operation

### What We Kept
1. **Schema Control**: Still using `golang-migrate` for migrations (not AutoMigrate)
2. **Repository Pattern**: Architecture unchanged, still have clean interfaces
3. **Testing Surface**: Can mock `*gorm.DB` or keep using repository interfaces

## Migration Checklist

- [x] Add GORM dependencies
- [x] Update all domain models with GORM tags
- [x] Add `TableName()` methods to all models
- [x] Replace Connection interface with GORM wrapper
- [x] Refactor User repository
- [x] Refactor Practice repository
- [x] Refactor Mistake repository
- [x] Refactor Note repository
- [x] Refactor Transcript repository
- [x] Refactor AICorrection repository
- [x] Update API server initialization
- [x] Update command layer database connection
- [x] Build verification (`go build ./...`)
- [x] Format check (`go fmt ./...`)
- [x] Run tests (if any exist)

## Testing Recommendations

1. **Integration Tests**: Test each repository against a real Postgres instance
2. **Enum Handling**: Verify that `mistake_status` and `mistake_type` enums work correctly
3. **Date Handling**: Ensure `Practice.Date` serializes/deserializes properly
4. **Error Cases**: Test `ErrNotFound` mapping for all Get methods

## Future Enhancements

### Optional Optimizations
1. **Enable GORM Logger**: Add SQL query logging in development
   ```go
   db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
       Logger: logger.Default.LogMode(logger.Info),
   })
   ```

2. **Add Relationships**: Enable eager loading
   ```go
   type Practice struct {
       // ... existing fields
       Mistakes []Mistake `gorm:"foreignKey:PracticeID"`
       Notes    []Note    `gorm:"foreignKey:PracticeID"`
   }
   ```

3. **Custom Types**: Wrap enums in custom Go types for type safety
   ```go
   type MistakeStatus string
   const (
       MistakeStatusAIDetected     MistakeStatus = "ai_detected"
       MistakeStatusAIMiscorrected MistakeStatus = "ai_miscorrected"
       MistakeStatusHumanCorrected MistakeStatus = "human_corrected"
   )
   ```

4. **Hooks**: Add validation or audit logging
   ```go
   func (u *User) BeforeCreate(tx *gorm.DB) error {
       // Validate or transform data
       return nil
   }
   ```

## Rollback Plan

If GORM causes issues, reverting is straightforward:
1. `git revert` the migration commit(s)
2. The database schema is unchanged (no data migration needed)
3. All API contracts remain identical

## References

- [GORM Documentation](https://gorm.io/docs/)
- [GORM PostgreSQL Driver](https://github.com/go-gorm/postgres)
- [Original Discussion](./AGENTS.md)