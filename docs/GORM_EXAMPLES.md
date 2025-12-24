# GORM Usage Examples

This document provides practical examples for working with GORM in the jpcorrect-backend project.

## Table of Contents
1. [Basic CRUD Operations](#basic-crud-operations)
2. [Querying with Conditions](#querying-with-conditions)
3. [Error Handling](#error-handling)
4. [Advanced Queries](#advanced-queries)
5. [Transaction Management](#transaction-management)
6. [Common Patterns](#common-patterns)

---

## Basic CRUD Operations

### Create a New User
```go
user := &domain.User{
    Name: "Tanaka Taro",
}

err := userRepo.Create(ctx, user)
if err != nil {
    log.Printf("failed to create user: %v", err)
    return err
}

// user.UserID is now populated by GORM
fmt.Printf("Created user with ID: %d\n", user.UserID)
```

### Read a User by ID
```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    if errors.Is(err, domain.ErrNotFound) {
        log.Println("User not found")
        return nil, err
    }
    log.Printf("database error: %v", err)
    return nil, err
}

fmt.Printf("Found user: %s (ID: %d)\n", user.Name, user.UserID)
```

### Update a User
```go
user, err := userRepo.GetByID(ctx, 1)
if err != nil {
    return err
}

user.Name = "Suzuki Hanako"
err = userRepo.Update(ctx, user)
if err != nil {
    log.Printf("failed to update user: %v", err)
    return err
}
```

### Delete a User
```go
err := userRepo.Delete(ctx, 1)
if err != nil {
    log.Printf("failed to delete user: %v", err)
    return err
}
```

---

## Querying with Conditions

### Find Users by Name
```go
users, err := userRepo.GetByName(ctx, "Tanaka Taro")
if err != nil {
    if errors.Is(err, domain.ErrNotFound) {
        log.Println("No users found with that name")
        return nil
    }
    return err
}

for _, user := range users {
    fmt.Printf("User ID: %d, Name: %s\n", user.UserID, user.Name)
}
```

### Find All Practices for a User
```go
practices, err := practiceRepo.GetByUserID(ctx, userID)
if err != nil {
    if errors.Is(err, domain.ErrNotFound) {
        log.Println("No practices found for this user")
        return nil
    }
    return err
}

for _, practice := range practices {
    fmt.Printf("Practice on %s: %.2f hours\n", 
        practice.Date.Format("2006-01-02"), 
        practice.Duration)
}
```

### Find All Mistakes for a Practice
```go
mistakes, err := mistakeRepo.GetByPracticeID(ctx, practiceID)
if err != nil {
    if errors.Is(err, domain.ErrNotFound) {
        log.Println("No mistakes found for this practice")
        return nil
    }
    return err
}

for _, mistake := range mistakes {
    fmt.Printf("Mistake: %s (%s) at %.2f-%.2f seconds\n",
        mistake.MistakeType,
        mistake.MistakeStatus,
        mistake.StartTime,
        mistake.EndTime)
}
```

---

## Error Handling

### Standard Error Checking Pattern
```go
func (s *SomeService) GetUser(ctx context.Context, id int) (*domain.User, error) {
    user, err := s.userRepo.GetByID(ctx, id)
    if err != nil {
        // Check for "not found" specifically
        if errors.Is(err, domain.ErrNotFound) {
            return nil, fmt.Errorf("user %d does not exist", id)
        }
        // All other errors are database/system errors
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### Handling Empty Results
```go
// For single-item queries, ErrNotFound is returned
practice, err := practiceRepo.GetByID(ctx, 999)
if errors.Is(err, domain.ErrNotFound) {
    // Handle missing practice
}

// For list queries, ErrNotFound means no items found
mistakes, err := mistakeRepo.GetByUserID(ctx, 123)
if errors.Is(err, domain.ErrNotFound) {
    // No mistakes for this user - this is valid
    mistakes = []*domain.Mistake{} // Return empty slice
}
```

---

## Advanced Queries

### Custom Raw SQL Query
If you need complex SQL that GORM's builder doesn't support well:

```go
type MistakeStats struct {
    UserID       int
    MistakeCount int
    AvgDuration  float64
}

func (r *postgresUserRepository) GetMistakeStats(ctx context.Context) ([]MistakeStats, error) {
    var stats []MistakeStats
    err := r.db.WithContext(ctx).Raw(`
        SELECT 
            user_id,
            COUNT(*) as mistake_count,
            AVG(end_time - start_time) as avg_duration
        FROM jpcorrect.mistake
        GROUP BY user_id
        ORDER BY mistake_count DESC
    `).Scan(&stats).Error
    
    return stats, err
}
```

### Ordering and Limiting Results
```go
// Get the 10 most recent practices
var practices []*domain.Practice
err := db.WithContext(ctx).
    Where("user_id = ?", userID).
    Order("date DESC").
    Limit(10).
    Find(&practices).Error
```

### Counting Records
```go
var count int64
err := db.WithContext(ctx).
    Model(&domain.Mistake{}).
    Where("practice_id = ? AND mistake_status = ?", practiceID, "ai_detected").
    Count(&count).Error

fmt.Printf("Found %d AI-detected mistakes\n", count)
```

---

## Transaction Management

### Basic Transaction
```go
func (s *Service) CreatePracticeWithNote(ctx context.Context, userID int, note string) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Create practice
        practice := &domain.Practice{
            UserID:   userID,
            Date:     time.Now(),
            Duration: 0,
        }
        if err := tx.Create(practice).Error; err != nil {
            return err // rollback
        }

        // Create note
        noteObj := &domain.Note{
            PracticeID: practice.PracticeID,
            Content:    note,
        }
        if err := tx.Create(noteObj).Error; err != nil {
            return err // rollback
        }

        return nil // commit
    })
}
```

### Manual Transaction Control
```go
func (s *Service) ComplexOperation(ctx context.Context) error {
    tx := s.db.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    // Operation 1
    if err := tx.Create(&someRecord).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Operation 2
    if err := tx.Model(&otherRecord).Update("field", "value").Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

---

## Common Patterns

### Batch Insert
```go
mistakes := []*domain.Mistake{
    {PracticeID: 1, UserID: 1, StartTime: 0, EndTime: 5, MistakeStatus: "ai_detected", MistakeType: "E1"},
    {PracticeID: 1, UserID: 1, StartTime: 10, EndTime: 15, MistakeStatus: "ai_detected", MistakeType: "E2"},
    {PracticeID: 1, UserID: 1, StartTime: 20, EndTime: 25, MistakeStatus: "ai_detected", MistakeType: "E3"},
}

// Create all at once
err := db.WithContext(ctx).Create(&mistakes).Error
if err != nil {
    log.Printf("batch insert failed: %v", err)
}
```

### Upsert (Insert or Update)
```go
// If you want to update if exists, insert if not
transcript := &domain.Transcript{
    MistakeID: 123,
    Content:   "新しいコンテンツ",
    Furigana:  "あたらしいこんてんつ",
    Accent:    "LHHHHHHH",
}

// First, try to find existing
var existing domain.Transcript
err := db.WithContext(ctx).Where("mistake_id = ?", transcript.MistakeID).First(&existing).Error

if errors.Is(err, gorm.ErrRecordNotFound) {
    // Doesn't exist, create it
    err = db.WithContext(ctx).Create(transcript).Error
} else if err == nil {
    // Exists, update it
    existing.Content = transcript.Content
    existing.Furigana = transcript.Furigana
    existing.Accent = transcript.Accent
    err = db.WithContext(ctx).Save(&existing).Error
}
```

### Pagination
```go
func (r *postgresPracticeRepository) GetPaginated(ctx context.Context, userID int, page, pageSize int) ([]*domain.Practice, error) {
    var practices []*domain.Practice
    offset := (page - 1) * pageSize
    
    err := r.db.WithContext(ctx).
        Where("user_id = ?", userID).
        Order("date DESC").
        Limit(pageSize).
        Offset(offset).
        Find(&practices).Error
    
    if err != nil {
        return nil, err
    }
    if len(practices) == 0 {
        return nil, domain.ErrNotFound
    }
    return practices, nil
}
```

### Preload Related Data (Future Enhancement)
```go
// After adding relationships to domain models:
// type Practice struct {
//     ...
//     Mistakes []Mistake `gorm:"foreignKey:PracticeID"`
// }

var practice domain.Practice
err := db.WithContext(ctx).
    Preload("Mistakes").
    Where("practice_id = ?", practiceID).
    First(&practice).Error

// practice.Mistakes is now populated
for _, mistake := range practice.Mistakes {
    fmt.Printf("Mistake: %s\n", mistake.MistakeType)
}
```

### Conditional Queries
```go
func (r *postgresMistakeRepository) Search(ctx context.Context, filters MistakeFilters) ([]*domain.Mistake, error) {
    query := r.db.WithContext(ctx).Model(&domain.Mistake{})
    
    if filters.UserID != nil {
        query = query.Where("user_id = ?", *filters.UserID)
    }
    
    if filters.Status != nil {
        query = query.Where("mistake_status = ?", *filters.Status)
    }
    
    if filters.MistakeType != nil {
        query = query.Where("mistake_type = ?", *filters.MistakeType)
    }
    
    if filters.MinDuration != nil {
        query = query.Where("end_time - start_time >= ?", *filters.MinDuration)
    }
    
    var mistakes []*domain.Mistake
    err := query.Find(&mistakes).Error
    
    if err != nil {
        return nil, err
    }
    if len(mistakes) == 0 {
        return nil, domain.ErrNotFound
    }
    return mistakes, nil
}
```

---

## Debugging Tips

### Enable SQL Logging (Development Only)
In `internal/cmd/api.go`:

```go
import "gorm.io/gorm/logger"

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info), // Shows all SQL queries
})
```

### Dry Run (Get SQL Without Executing)
```go
stmt := db.WithContext(ctx).
    Where("user_id = ?", 123).
    Session(&gorm.Session{DryRun: true}).
    First(&domain.Practice{})

fmt.Println(stmt.Statement.SQL.String())  // Prints the SQL
fmt.Println(stmt.Statement.Vars)          // Prints the parameters
```

---

## Performance Considerations

### Use Select to Limit Columns
```go
// Only fetch specific columns
var users []domain.User
db.WithContext(ctx).
    Select("user_id", "name").
    Where("user_id IN ?", []int{1, 2, 3}).
    Find(&users)
```

### Avoid N+1 Queries
```go
// BAD: N+1 query problem
practices, _ := practiceRepo.GetByUserID(ctx, userID)
for _, practice := range practices {
    mistakes, _ := mistakeRepo.GetByPracticeID(ctx, practice.PracticeID) // N queries!
}

// GOOD: Fetch all at once
practices, _ := practiceRepo.GetByUserID(ctx, userID)
practiceIDs := make([]int, len(practices))
for i, p := range practices {
    practiceIDs[i] = p.PracticeID
}

var allMistakes []*domain.Mistake
db.WithContext(ctx).
    Where("practice_id IN ?", practiceIDs).
    Find(&allMistakes) // Single query
```

### Batch Operations
```go
// Use CreateInBatches for large inserts
db.WithContext(ctx).CreateInBatches(mistakes, 100) // Insert 100 at a time
```

---

## Common Gotchas

### Zero Values in Updates
```go
// This WON'T update duration to 0!
practice.Duration = 0
db.Save(&practice)

// Use Updates with a map instead
db.Model(&practice).Updates(map[string]interface{}{
    "duration": 0,
})

// Or use Select to specify fields
db.Model(&practice).Select("duration").Updates(practice)
```

### Date/Time Handling
```go
// Parse date strings properly
dateStr := "2024-01-15"
date, err := time.Parse("2006-01-02", dateStr)

practice := &domain.Practice{
    UserID:   1,
    Date:     date,
    Duration: 1.5,
}
```

### Enum Values
```go
// Use string literals matching your Postgres enum
mistake := &domain.Mistake{
    MistakeStatus: "ai_detected",     // Valid: ai_detected, ai_miscorrected, human_corrected
    MistakeType:   "E1",               // Valid: E1-E9
}
```

---

## References
- [GORM Official Documentation](https://gorm.io/docs/)
- [GORM Query Guide](https://gorm.io/docs/query.html)
- [GORM Advanced Query](https://gorm.io/docs/advanced_query.html)
- [GORM Performance Tips](https://gorm.io/docs/performance.html)