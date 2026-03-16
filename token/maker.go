package token

import "time"

//interfaz para manejar tokens
type Maker interface {
	// CreateToken para crear un nuevo token para un username especifico y duracion
	CreateToken(username string, duration time.Duration) (string, *Payload, error)

	// Comprobar que el token es valido
	VerifyToken(token string) (*Payload, error)
}
