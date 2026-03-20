package gapi

import (
	"fmt"

	db "github.com/cristianemek/go-simplebank/db/sqlc"
	"github.com/cristianemek/go-simplebank/pb"
	"github.com/cristianemek/go-simplebank/token"
	"github.com/cristianemek/go-simplebank/util"
)

// Server va a servir todas las peticiones gRPC para nuestro proyecto
type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

// Funcion para crear una nueva instancia del servidor y configurar las rutas
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
