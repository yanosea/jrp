package history

import (
	"context"
)

// HistoryRepository is an interface that provides the repository for the history table in the jrp database.
type HistoryRepository interface {
	DeleteAll(ctx context.Context) (int, error)
	DeleteByIdIn(ctx context.Context, ids []int) (int, error)
	DeleteByIdInAndIsFavoritedIs(ctx context.Context, ids []int, isFavorited int) (int, error)
	DeleteByIsFavoritedIs(ctx context.Context, isFavorited int) (int, error)
	FindAll(ctx context.Context) ([]*History, error)
	FindByIsFavoritedIs(ctx context.Context, isFavorited int) ([]*History, error)
	FindByIsFavoritedIsAndPhraseContains(ctx context.Context, keywords []string, and bool, isFavorited int) ([]*History, error)
	FindByPhraseContains(ctx context.Context, keywords []string, and bool) ([]*History, error)
	FindTopNByIsFavoritedIsAndByOrderByIdAsc(ctx context.Context, number int, isFavorited int) ([]*History, error)
	FindTopNByIsFavoritedIsAndByPhraseContainsOrderByIdAsc(ctx context.Context, keywords []string, and bool, number int, isFavorited int) ([]*History, error)
	FindTopNByOrderByIdAsc(ctx context.Context, number int) ([]*History, error)
	FindTopNByPhraseContainsOrderByIdAsc(ctx context.Context, keywords []string, and bool, number int) ([]*History, error)
	SaveAll(ctx context.Context, jrps []*History) ([]*History, error)
	UpdateIsFavoritedByIdIn(ctx context.Context, isFavorited int, ids []int) (int, error)
	UpdateIsFavoritedByIsFavoritedIs(ctx context.Context, isFavorited int, isFavoritedIs int) (int, error)
}

