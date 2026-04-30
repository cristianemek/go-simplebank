package gapi

import (
	"context"
	"errors"

	db "github.com/cristianemek/go-simplebank/db/sqlc"
	"github.com/cristianemek/go-simplebank/pb"
	"github.com/cristianemek/go-simplebank/util"
	"github.com/cristianemek/go-simplebank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "error getting user  %s", err)
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "invalidad password: %s", err)
	}

	accesToken, accesPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccesTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating acces token: %s", err)
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating refresh token: %s", err)
	}

	mtdt := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx,
		db.CreateSessionParams{
			ID:           refreshPayload.ID,
			Username:     user.Username,
			RefreshToken: refreshToken,
			UserAgent:    mtdt.UserAgent,
			ClientIp:     mtdt.ClientIp,
			IsBlocked:    false,
			ExpiresAt:    refreshPayload.ExpiresAt.Time,
		})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating session: %s", err)
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accesToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accesPayload.ExpiresAt.Time),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
	}

	return rsp, nil

}

func validateLoginUserRequest(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
