package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"jpcorrect-backend/internal/domain"
)

func TestGormGuildDefaultSlotRepository_GetByGuildID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormGuildDefaultSlotRepository(db)
	guildID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "guild_default_slot" WHERE guild_id = $1 ORDER BY "guild_default_slot"."guild_id" LIMIT $2`)).
			WithArgs(guildID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"guild_id", "day_of_week", "time_of_day", "timezone", "created_at", "updated_at"}).
				AddRow(guildID, 3, "19:00", "Asia/Taipei", time.Now(), time.Now()))

		slot, err := repo.GetByGuildID(context.Background(), guildID)

		assert.NoError(t, err)
		assert.NotNil(t, slot)
		assert.Equal(t, guildID, slot.GuildID)
		assert.Equal(t, 3, slot.DayOfWeek)
		assert.Equal(t, "19:00", slot.TimeOfDay)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "guild_default_slot" WHERE guild_id = $1 ORDER BY "guild_default_slot"."guild_id" LIMIT $2`)).
			WithArgs(guildID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		slot, err := repo.GetByGuildID(context.Background(), guildID)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, slot)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "guild_default_slot" WHERE guild_id = $1 ORDER BY "guild_default_slot"."guild_id" LIMIT $2`)).
			WithArgs(guildID, 1).
			WillReturnError(fmt.Errorf("db error"))

		slot, err := repo.GetByGuildID(context.Background(), guildID)

		assert.Error(t, err)
		assert.Nil(t, slot)
	})
}

func TestGormGuildDefaultSlotRepository_Upsert(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormGuildDefaultSlotRepository(db)

	t.Run("Success", func(t *testing.T) {
		slot := &domain.GuildDefaultSlot{
			GuildID:   uuid.New(),
			DayOfWeek: 3,
			TimeOfDay: "19:00",
			Timezone:  "Asia/Taipei",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "guild_default_slot"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Upsert(context.Background(), slot)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		slot := &domain.GuildDefaultSlot{
			GuildID:   uuid.New(),
			DayOfWeek: 5,
			TimeOfDay: "20:00",
			Timezone:  "Asia/Taipei",
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "guild_default_slot"`)).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Upsert(context.Background(), slot)

		assert.Error(t, err)
	})
}

func TestGormGuildDefaultSlotRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormGuildDefaultSlotRepository(db)
	guildID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "guild_default_slot" WHERE guild_id = $1`)).
			WithArgs(guildID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Delete(context.Background(), guildID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "guild_default_slot" WHERE guild_id = $1`)).
			WithArgs(guildID).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Delete(context.Background(), guildID)

		assert.Error(t, err)
	})
}
