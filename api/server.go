package api

import (
	db "github.com/cristianemek/go-simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server va a servir todas las peticiones http para nuestro proyecto
type Server struct {
	store  *db.Store
	router *gin.Engine
}

// Funcion para crear una nueva instancia del servidor y configurar las rutas
func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	// Configurar las rutas
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

// Funcion para iniciar el servidor en una direccion especifica
func (server *Server) Start(address string) error {
	return server.router.Run(address) // campo de enrutador es privado, solo se puede acceder a el dentro del paquete api
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
