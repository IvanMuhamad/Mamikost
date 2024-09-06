package server

import (
	"context"
	"log"
	"mamikost/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func InitDatabase(config *config.Config) *pgxpool.Conn {
	connectionString := viper.GetString("database.connection_string")
	maxIdleConnections := viper.GetDuration("database.max_idle_connections")
	maxOpenConnections := viper.GetInt32("database.max_open_connections")
	connectionMaxLifetime := viper.GetDuration("database.connection_max_lifetime")
	//driverName := config.GetString("database.driver_name")

	if connectionString == "" {
		log.Fatalf("Database connection string is missing")
	}

	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = maxOpenConnections
	dbConfig.MaxConnIdleTime = maxIdleConnections
	dbConfig.MaxConnLifetime = connectionMaxLifetime

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := connPool.Acquire(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = connection.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	defer log.Println("Database already connected")

	return connection
}

func Close(conn *pgxpool.Conn) {
	err := conn.Conn().Close(context.Background())
	log.Println("Database already closed..")
	if err != nil {
		log.Fatalf("Unable to close connection %v\n", err)
	}
}
