package gapi

import (
	"fmt"

	db "github.com/cristianemek/go-simplebank/db/sqlc"
	"github.com/cristianemek/go-simplebank/pb"
	"github.com/cristianemek/go-simplebank/token"
	"github.com/cristianemek/go-simplebank/util"
	"github.com/cristianemek/go-simplebank/worker"
)

// Server va a servir todas las peticiones gRPC para nuestro proyecto
type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

// Funcion para crear una nueva instancia del servidor y configurar las rutas
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
