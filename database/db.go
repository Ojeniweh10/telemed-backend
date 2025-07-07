package database

import (
	"context"
	"fmt"
	"telemed/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewConnection() *pgxpool.Pool {
	host := config.Db().Host
	user := config.Db().User
	password := config.Db().Password
	dbname := config.Db().Name

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, dbname)
	ctx := context.Background()

	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		panic(err)
	}

	poolConfig.MaxConns = 5
	dbPool, err := pgxpool.ConnectConfig(ctx, poolConfig)

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")
	return dbPool
}

func Insert(db *pgxpool.Pool, table string, data map[string]any) error {
	var columns, placeholders string
	var values []any
	i := 1
	for k, v := range data {
		if i > 1 {
			columns += ","
			placeholders += ","
		}
		columns += k
		placeholders += fmt.Sprintf("$%d", i)
		values = append(values, v)
		i++
	}
	sql := fmt.Sprintf("INSERT INTO %v (%s) VALUES (%s)", table, columns, placeholders)
	_, err := db.Exec(context.Background(), sql, values...)
	if err != nil {
		return err
	}

	return nil
}
