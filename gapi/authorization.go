package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/cristianemek/go-simplebank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)

	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)

	if len(values) == 0 {
		return nil, fmt.Errorf("missing auth header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != "Bearer" {
		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
	}

	accesToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accesToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return payload, nil
}
