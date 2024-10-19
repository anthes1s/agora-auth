package authctx

import (
	"context"
	"os"
	"regexp"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// String need to consist of only alphanumeric symbols + underscores and dots, not shorter than 4 symbols and not longer than 64
	regexUsername = `^[a-zA-Z0-9_.]{4,64}$`
)

type AuthContext struct {
	Pool *pgxpool.Pool
	// This is bad, but it'll work ¯\_(ツ)_/¯
	ValidateUsername *regexp.Regexp
}

func New() (*AuthContext, error) {
	pl, err := acquireDBPool()
	if err != nil {
		return nil, err
	}

	reUsername, err := regexp.Compile(regexUsername)
	if err != nil {
		return nil, err
	}

	ctx := AuthContext{
		Pool:             pl,
		ValidateUsername: reUsername,
	}

	return &ctx, nil
}

func acquireDBPool() (*pgxpool.Pool, error) {
	dbcfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), dbcfg)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return pool, nil
}
