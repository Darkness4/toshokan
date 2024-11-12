package indexer

import (
	"context"
	"time"

	"github.com/Darkness4/toshokan/db"
	"github.com/Darkness4/toshokan/scan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Indexer struct {
	scanner *scan.Scanner

	db *pgxpool.Pool
	q  *db.Queries
}

func NewIndexer(scanner *scan.Scanner, db *pgxpool.Pool, q *db.Queries) *Indexer {
	return &Indexer{
		scanner: scanner,
		db:      db,
		q:       q,
	}
}

func (i *Indexer) Init() {
	i.scanner.Init()
}

func (i *Indexer) Index(ctx context.Context) {
	for res := range i.scanner.Scan() {
		i.indexMeta(ctx, res)
	}
}

func (i *Indexer) indexMeta(ctx context.Context, res scan.ScanResult) error {
	// Prepare TX
	tx, err := i.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := i.q.WithTx(tx)

	// Insert tags
	for _, category := range res.Meta.Categories {
		err := q.CreateTag(ctx, db.CreateTagParams{
			Namespace: pgtype.Text{
				String: category.Namespace,
				Valid:  true,
			},
			Value: pgtype.Text{
				String: category.Value,
				Valid:  true,
			},
		})
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
	}

	now := time.Now()

	// TODO: find thumbnail path
	// TODO: parse series

	// Insert meta
	if err := q.CreateArchive(
		ctx,
		db.CreateArchiveParams{
			Title: pgtype.Text{
				String: res.Meta.Title,
				Valid:  true,
			},
			DateAdded: pgtype.Timestamptz{
				Time:  now,
				Valid: true,
			},
			DateUpdated: pgtype.Timestamptz{
				Time:  now,
				Valid: true,
			},
			DateIssued: pgtype.Date{
				Time:  time.Unix(res.Meta.Issued, 0),
				Valid: true,
			},
			FilePath: pgtype.Text{
				String: res.Path,
				Valid:  true,
			},
		},
	); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	// Commit TX
	return tx.Commit(ctx)
}
