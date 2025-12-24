# GORM Quick Reference Card

## Common Operations

### Create
```go
user := &domain.User{Name: "Tanaka"}
err := repo.Create(ctx, user)
// user.UserID now populated
```

### Read One
```go
user, err := repo.GetByID(ctx, 1)
if errors.Is(err, domain.ErrNotFound) {
    // Not found
}
```

### Read Many
```go
users, err := repo.GetByName(ctx, "Tanaka")
if errors.Is(err, domain.ErrNotFound) {
    // Empty result
}
```

### Update
```go
user.Name = "Updated"
err := repo.Update(ctx, user)
```

### Delete
```go
err := repo.Delete(ctx, 1)
```

## Direct GORM Queries

### Find Single Record
```go
var user domain.User
db.WithContext(ctx).Where("user_id = ?", 1).First(&user)
```

### Find Multiple Records
```go
var users []domain.User
db.WithContext(ctx).Where("name LIKE ?", "Tanaka%").Find(&users)
```

### Count
```go
var count int64
db.WithContext(ctx).Model(&domain.User{}).Where("name = ?", "Tanaka").Count(&count)
```

### Order & Limit
```go
db.WithContext(ctx).Order("date DESC").Limit(10).Find(&practices)
```

### Raw SQL
```go
db.WithContext(ctx).Raw("SELECT * FROM jpcorrect.user WHERE user_id = ?", 1).Scan(&user)
```

## Transactions

### Auto Transaction
```go
err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&record1).Error; err != nil {
        return err // auto rollback
    }
    if err := tx.Create(&record2).Error; err != nil {
        return err // auto rollback
    }
    return nil // auto commit
})
```

### Manual Transaction
```go
tx := db.WithContext(ctx).Begin()
defer tx.Rollback() // rollback if not committed

if err := tx.Create(&record).Error; err != nil {
    return err
}

return tx.Commit().Error
```

## Error Handling

```go
err := db.WithContext(ctx).First(&user, id).Error
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return domain.ErrNotFound
    }
    return err
}
```

## Batch Operations

### Batch Insert
```go
users := []*domain.User{{Name: "A"}, {Name: "B"}}
db.WithContext(ctx).Create(&users)
```

### Batch Insert with Size
```go
db.WithContext(ctx).CreateInBatches(users, 100)
```

## Advanced Queries

### Conditional Query Building
```go
query := db.WithContext(ctx).Model(&domain.Mistake{})
if userID != nil {
    query = query.Where("user_id = ?", *userID)
}
if status != nil {
    query = query.Where("mistake_status = ?", *status)
}
query.Find(&mistakes)
```

### Pagination
```go
offset := (page - 1) * pageSize
db.WithContext(ctx).
    Offset(offset).
    Limit(pageSize).
    Find(&records)
```

### Joins (Manual)
```go
db.WithContext(ctx).
    Table("jpcorrect.practice p").
    Joins("INNER JOIN jpcorrect.user u ON p.user_id = u.user_id").
    Where("u.name = ?", "Tanaka").
    Find(&practices)
```

## Updates

### Update Single Field
```go
db.WithContext(ctx).Model(&user).Update("name", "NewName")
```

### Update Multiple Fields (Map)
```go
db.WithContext(ctx).Model(&user).Updates(map[string]interface{}{
    "name": "NewName",
    "duration": 0, // Zero values work with map
})
```

### Update Multiple Fields (Struct - ignores zero values)
```go
db.WithContext(ctx).Model(&user).Updates(user)
```

### Update with Select (Include zero values)
```go
db.WithContext(ctx).Model(&practice).Select("duration").Updates(practice)
```

## Model Tags Reference

```go
type User struct {
    UserID int    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
    Name   string `json:"name" gorm:"column:name;not null"`
    Email  string `json:"email" gorm:"column:email;unique;index"`
    Active bool   `json:"active" gorm:"column:active;default:true"`
}

func (User) TableName() string {
    return "jpcorrect.user"
}
```

### Common GORM Tags
- `column:name` - Column name
- `primaryKey` - Primary key
- `autoIncrement` - Auto increment
- `not null` - NOT NULL constraint
- `unique` - UNIQUE constraint
- `index` - Create index
- `default:value` - Default value
- `type:text` - Custom SQL type
- `size:255` - Column size
- `foreignKey:UserID` - Foreign key

## Debugging

### Enable SQL Logging
```go
import "gorm.io/gorm/logger"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### Dry Run (Get SQL)
```go
stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1)
fmt.Println(stmt.Statement.SQL.String())
fmt.Println(stmt.Statement.Vars)
```

## Common Gotchas

### Zero Values in Updates
```go
// BAD: Won't update duration to 0
db.Save(&practice) // Ignores zero values

// GOOD: Updates zero values
db.Model(&practice).Updates(map[string]interface{}{"duration": 0})
```

### Date Handling
```go
// Parse string to time.Time
date, _ := time.Parse("2006-01-02", "2024-01-15")
practice.Date = date
```

### Enum Values
```go
mistake.MistakeStatus = "ai_detected" // Must match DB enum
mistake.MistakeType = "E1"            // Must match DB enum
```

## Connection Pool

```go
sqlDB, err := db.DB()
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

## Performance Tips

1. **Select only needed columns:**
   ```go
   db.Select("user_id", "name").Find(&users)
   ```

2. **Use batch operations:**
   ```go
   db.CreateInBatches(records, 100)
   ```

3. **Avoid N+1 with Preload:**
   ```go
   db.Preload("Mistakes").Find(&practices)
   ```

4. **Use indexes on WHERE clauses**

5. **Limit result sets:**
   ```go
   db.Limit(1000).Find(&records)
   ```

## Repository Pattern Template

```go
type postgresXxxRepository struct {
    db *gorm.DB
}

func NewPostgresXxx(conn *repository.Connection) domain.XxxRepository {
    return &postgresXxxRepository{db: conn.DB}
}

func (r *postgresXxxRepository) GetByID(ctx context.Context, id int) (*domain.Xxx, error) {
    var record domain.Xxx
    err := r.db.WithContext(ctx).Where("xxx_id = ?", id).First(&record).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return &record, nil
}

func (r *postgresXxxRepository) Create(ctx context.Context, record *domain.Xxx) error {
    return r.db.WithContext(ctx).Create(record).Error
}

func (r *postgresXxxRepository) Update(ctx context.Context, record *domain.Xxx) error {
    result := r.db.WithContext(ctx).Model(record).Where("xxx_id = ?", record.XxxID).Updates(map[string]interface{}{
        "field1": record.Field1,
        "field2": record.Field2,
    })
    return result.Error
}

func (r *postgresXxxRepository) Delete(ctx context.Context, id int) error {
    result := r.db.WithContext(ctx).Where("xxx_id = ?", id).Delete(&domain.Xxx{})
    return result.Error
}
```

---

**For more details, see:**
- [GORM_MIGRATION.md](./GORM_MIGRATION.md) - Migration guide
- [GORM_EXAMPLES.md](./GORM_EXAMPLES.md) - Detailed examples
- [GORM Docs](https://gorm.io/docs/) - Official documentation