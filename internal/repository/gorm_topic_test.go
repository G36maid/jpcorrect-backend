package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"jpcorrect-backend/internal/domain"
)

func TestGormTopicRepository_GetByID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormTopicRepository(db)
	topicID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE id = $1 AND "topic"."deleted_at" IS NULL ORDER BY "topic"."id" LIMIT $2`)).
			WithArgs(topicID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}).
				AddRow(topicID, domain.TopicKindAnnounce, "テスト", "medium", nil, nil))

		topic, err := repo.GetByID(context.Background(), topicID)

		assert.NoError(t, err)
		assert.NotNil(t, topic)
		assert.Equal(t, topicID, topic.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE id = $1 AND "topic"."deleted_at" IS NULL ORDER BY "topic"."id" LIMIT $2`)).
			WithArgs(topicID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		topic, err := repo.GetByID(context.Background(), topicID)

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, topic)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE id = $1 AND "topic"."deleted_at" IS NULL ORDER BY "topic"."id" LIMIT $2`)).
			WithArgs(topicID, 1).
			WillReturnError(fmt.Errorf("db error"))

		topic, err := repo.GetByID(context.Background(), topicID)

		assert.Error(t, err)
		assert.Nil(t, topic)
	})
}

func TestGormTopicRepository_SearchAnnounce(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormTopicRepository(db)

	t.Run("Success", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE (kind = $1 AND title_jp LIKE $2) AND "topic"."deleted_at" IS NULL`)).
			WithArgs(domain.TopicKindAnnounce, "%テスト%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}).
				AddRow(id, domain.TopicKindAnnounce, "テストのトピック", "medium", nil, nil))

		topics, err := repo.SearchAnnounce(context.Background(), "テスト")

		assert.NoError(t, err)
		assert.Len(t, topics, 1)
		assert.Equal(t, id, topics[0].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE (kind = $1 AND title_jp LIKE $2) AND "topic"."deleted_at" IS NULL`)).
			WithArgs(domain.TopicKindAnnounce, "%notfound%").
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}))

		topics, err := repo.SearchAnnounce(context.Background(), "notfound")

		assert.NoError(t, err)
		assert.Empty(t, topics)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE (kind = $1 AND title_jp LIKE $2) AND "topic"."deleted_at" IS NULL`)).
			WithArgs(domain.TopicKindAnnounce, "%error%").
			WillReturnError(fmt.Errorf("db error"))

		topics, err := repo.SearchAnnounce(context.Background(), "error")

		assert.Error(t, err)
		assert.Nil(t, topics)
	})
}

func TestGormTopicRepository_GetRandom(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormTopicRepository(db)

	t.Run("Success", func(t *testing.T) {
		id := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL ORDER BY RANDOM(),"topic"."id" LIMIT $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}).
				AddRow(id, domain.TopicKindRandom, "ランダム", "easy", nil, nil))

		topic, err := repo.GetRandom(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, topic)
		assert.Equal(t, id, topic.ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL ORDER BY RANDOM(),"topic"."id" LIMIT $1`)).
			WithArgs(1).
			WillReturnError(gorm.ErrRecordNotFound)

		topic, err := repo.GetRandom(context.Background())

		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Nil(t, topic)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL ORDER BY RANDOM(),"topic"."id" LIMIT $1`)).
			WithArgs(1).
			WillReturnError(fmt.Errorf("db error"))

		topic, err := repo.GetRandom(context.Background())

		assert.Error(t, err)
		assert.Nil(t, topic)
	})
}

func TestGormTopicRepository_List(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewGormTopicRepository(db)

	t.Run("Success", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}).
				AddRow(id1, domain.TopicKindAnnounce, "トピック1", "medium", nil, nil).
				AddRow(id2, domain.TopicKindRandom, "トピック2", "hard", nil, nil))

		topics, err := repo.List(context.Background())

		assert.NoError(t, err)
		assert.Len(t, topics, 2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("EmptyResult", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "kind", "title_jp", "difficulty", "hint_vocab", "hint_grammar"}))

		topics, err := repo.List(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, topics)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "topic" WHERE "topic"."deleted_at" IS NULL`)).
			WillReturnError(fmt.Errorf("db error"))

		topics, err := repo.List(context.Background())

		assert.Error(t, err)
		assert.Nil(t, topics)
	})
}
