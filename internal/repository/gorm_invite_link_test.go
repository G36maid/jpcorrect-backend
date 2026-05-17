package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"jpcorrect-backend/internal/domain"
)

func TestGormInviteLinkRepository_GetByToken(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormInviteLinkRepository(db)

	t.Run("Success", func(t *testing.T) {
		token := "abc123token456def789gh012ij345"
		linkID := uuid.New()
		guildID := uuid.New()
		userID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE token = $1 ORDER BY "invite_link"."id" LIMIT $2`)).
			WithArgs(token, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "guild_id", "token", "expires_at", "created_by_user_id", "created_at"}).
				AddRow(linkID, guildID, token, time.Now(), userID, time.Now()))

		link, err := repo.GetByToken(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, link)
		assert.Equal(t, token, link.Token)
		assert.Equal(t, guildID, link.GuildID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		token := "nonexistent_token_1234567890ab"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE token = $1 ORDER BY "invite_link"."id" LIMIT $2`)).
			WithArgs(token, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		link, err := repo.GetByToken(context.Background(), token)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, link)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		token := "err_token_1234567890abcdef"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE token = $1 ORDER BY "invite_link"."id" LIMIT $2`)).
			WithArgs(token, 1).
			WillReturnError(fmt.Errorf("db error"))

		link, err := repo.GetByToken(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, link)
	})
}

func TestGormInviteLinkRepository_GetByGuildID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormInviteLinkRepository(db)
	guildID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE guild_id = $1`)).
			WithArgs(guildID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "guild_id", "token", "expires_at", "created_by_user_id", "created_at"}).
				AddRow(uuid.New(), guildID, "token111111111111111111111", time.Now(), uuid.New(), time.Now()).
				AddRow(uuid.New(), guildID, "token222222222222222222222", time.Now(), uuid.New(), time.Now()))

		links, err := repo.GetByGuildID(context.Background(), guildID)

		assert.NoError(t, err)
		assert.Len(t, links, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE guild_id = $1`)).
			WithArgs(guildID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "guild_id", "token", "expires_at", "created_by_user_id", "created_at"}))

		links, err := repo.GetByGuildID(context.Background(), guildID)

		assert.NoError(t, err)
		assert.Empty(t, links)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "invite_link" WHERE guild_id = $1`)).
			WithArgs(guildID).
			WillReturnError(fmt.Errorf("db error"))

		links, err := repo.GetByGuildID(context.Background(), guildID)

		assert.Error(t, err)
		assert.Nil(t, links)
	})
}

func TestGormInviteLinkRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormInviteLinkRepository(db)

	t.Run("Success", func(t *testing.T) {
		link := &domain.InviteLink{
			GuildID:         uuid.New(),
			Token:           "newtoken123456789012345678",
			ExpiresAt:       time.Now().Add(24 * time.Hour),
			CreatedByUserID: uuid.New(),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "invite_link"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), link)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, link.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DuplicateEntry", func(t *testing.T) {
		link := &domain.InviteLink{
			GuildID:         uuid.New(),
			Token:           "duplicatetoken123456789012",
			ExpiresAt:       time.Now().Add(24 * time.Hour),
			CreatedByUserID: uuid.New(),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "invite_link"`)).
			WillReturnError(&pgconn.PgError{Code: "23505"})
		mock.ExpectRollback()

		err := repo.Create(context.Background(), link)

		assert.ErrorIs(t, err, domain.ErrDuplicateEntry)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		link := &domain.InviteLink{
			GuildID:         uuid.New(),
			Token:           "err_token_123456789012345",
			ExpiresAt:       time.Now().Add(24 * time.Hour),
			CreatedByUserID: uuid.New(),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "invite_link"`)).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), link)

		assert.Error(t, err)
	})
}

func TestGormInviteLinkRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormInviteLinkRepository(db)

	t.Run("Success", func(t *testing.T) {
		link := &domain.InviteLink{
			ID:        uuid.New(),
			GuildID:   uuid.New(),
			Token:     "updated_token_123456789012",
			ExpiresAt: time.Now().Add(48 * time.Hour),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invite_link"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Update(context.Background(), link)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		link := &domain.InviteLink{
			ID:        uuid.New(),
			GuildID:   uuid.New(),
			Token:     "dberr_token_12345678901234",
			ExpiresAt: time.Now().Add(48 * time.Hour),
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "invite_link"`)).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Update(context.Background(), link)

		assert.Error(t, err)
	})
}
