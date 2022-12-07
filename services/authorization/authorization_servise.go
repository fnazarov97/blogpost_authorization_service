package authorization

import (
	"blockpost/genprotos/authorization"
	"blockpost/util"
	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login ...
func (s *authorizationService) Login(ctx context.Context, req *authorization.LoginRequest) (*authorization.TokenResponse, error) {
	log.Println("Login...")

	errAuth := errors.New("username or password wrong")

	user, err := s.stg.GetUserByUsername(req.Username)
	if err != nil {
		log.Println(err.Error())
		return nil, status.Errorf(codes.Unauthenticated, errAuth.Error())
	}

	match, err := util.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "util.ComparePassword: %s", err.Error())
	}

	if !match {
		return nil, status.Errorf(codes.Unauthenticated, errAuth.Error())
	}

	m := map[string]interface{}{
		"user_id":  user.Id,
		"username": user.Username,
	}

	tokenStr, err := util.GenerateJWT(m, 10*time.Minute, s.cfg.SecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "util.GenerateJWT: %s", err.Error())
	}

	return &authorization.TokenResponse{
		Token: tokenStr,
	}, nil
}

// HasAccess ...
func (s *authorizationService) HasAccess(ctx context.Context, req *authorization.TokenRequest) (*authorization.HasAccessResponse, error) {
	log.Println("HasAccess...")

	result, err := util.ParseClaims(req.Token, s.cfg.SecretKey)
	if err != nil {
		log.Println(status.Errorf(codes.Unauthenticated, "util.ParseClaims: %s", err.Error()))
		return &authorization.HasAccessResponse{
			User:      nil,
			HasAccess: false,
		}, nil
	}

	log.Println(result.Username)

	user, err := s.stg.GetUserByID(result.UserID)
	if err != nil {
		log.Println(status.Errorf(codes.Unauthenticated, "s.stg.GetUserByID: %s", err.Error()))
		return &authorization.HasAccessResponse{
			User:      nil,
			HasAccess: false,
		}, nil
	}

	return &authorization.HasAccessResponse{
		User:      user,
		HasAccess: true,
	}, nil
}
