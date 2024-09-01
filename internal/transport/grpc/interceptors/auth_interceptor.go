package interceptors

import (
	"context"
	"passkeeper/internal/entities/hashes"

	"passkeeper/internal/config"
	"passkeeper/internal/usecase/srv/usersUC"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrWrongAuth = "not authorized"
)

type AuthInterceptor struct {
	service usersUC.UserUsecaseInf
}

func (in *AuthInterceptor) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := hashes.ExtractUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := hashes.CheckToken(token, in.service.GetTokenSalt())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, ErrWrongAuth)
	}

	metadata.AppendToOutgoingContext(ctx, config.TokenHeader, userID)

	return ctx, nil
}

func NewAuthInterceptor(service usersUC.UserUsecaseInf) *AuthInterceptor {
	return &AuthInterceptor{service: service}
}

// Will use `grpc-ecosystem/go-grpc-middleware/blob/main/interceptors/auth`
//
//func (in *AuthInterceptor) Intercept(ctx context.Context, req interface{},
//	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//
//	token, err := handlers.ExtractUserInfo(ctx)
//	if err != nil {
//		in.log.Error().Err(err).
//			Msg(ErrWrongAuth)
//
//		return nil, status.Errorf(codes.Unauthenticated, ErrWrongAuth)
//	}
//
//	userID, err := usecase.CheckToken(token, in.service.GetTokenSalt())
//	if err != nil {
//		in.log.Error().Err(err).
//			Msg(ErrWrongAuth)
//
//		return nil, status.Errorf(codes.Unauthenticated, ErrWrongAuth)
//	}
//
//	metadata.AppendToOutgoingContext(ctx, config.TokenHeader, userID)
//
//	return handler(ctx, req)
//}
