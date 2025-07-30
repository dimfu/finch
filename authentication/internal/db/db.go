package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var Pool *pgxpool.Pool

type config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DBNAME   string
}

func getConfig() (*config, error) {
	var cfg config
	envType := reflect.TypeOf(cfg)
	envValue := reflect.ValueOf(&cfg).Elem()

	for i := 0; i < envType.NumField(); i++ {
		field := envType.Field(i)
		envVar := field.Name

		value := os.Getenv(envVar)
		if value == "" {
			return nil, errors.New("environment variable " + envVar + " is required")
		}

		envValue.FieldByName(envVar).SetString(value)
	}

	return &cfg, nil
}

func Connect() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_DBNAME,
	)

	ctx := context.Background()
	Pool, err = pgxpool.New(ctx, url)

	if err != nil {
		return err
	}

	if err := Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
