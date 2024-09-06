package services

import (
	"context"
	db "mamikost/db/sqlc"
	"mamikost/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=./mocks/mock_store.go -build_flags=-mod=mod -package=mocks . Store
type Store interface {
	db.Querier
	CreateOrderTx(ctx context.Context, args db.CreateOrderParams) (*db.OrderRentProperty, error)
	Signup(ctx context.Context, userReq models.CreateUserRequest) (*models.UserResponse, *models.Error)
	Signin(ctx context.Context, userReq models.CreateUserRequest) (*models.UserResponse, *models.Error)
	Signout(ctx context.Context, accessToken string) *models.Error
}

type StoreManager struct {
	*db.Queries // implements Querier
	dbConn      *pgxpool.Conn
}

func NewStoreManager(dbConn *pgxpool.Conn) Store {
	return &StoreManager{
		Queries: db.New(dbConn),
		dbConn:  dbConn,
	}
}
