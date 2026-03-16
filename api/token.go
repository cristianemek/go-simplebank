package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccesToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccesTokenResponse struct {
	AccesToken          string    `json:"access_token"`
	AccesTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewAccesToken(ctx *gin.Context) {
	var req renewAccesToken
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if session.IsBlocked {
		err := fmt.Errorf("blocked session")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("missmatched session token")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("expired session")
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	accesToken, accesPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccesTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	rsp := renewAccesTokenResponse{
		AccesToken:          accesToken,
		AccesTokenExpiresAt: accesPayload.ExpiresAt.Time,
	}

	ctx.JSON(http.StatusOK, rsp)
}
