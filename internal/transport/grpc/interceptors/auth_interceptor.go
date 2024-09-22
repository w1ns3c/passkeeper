package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/usecase/srv/usersUC"
)

var (
	ErrWrongAuth = "not authorized"
)

type AuthInterceptor struct {
	service usersUC.UserUsecaseInf
	grpc.UnaryServerInterceptor
}

func (in *AuthInterceptor) AuthFunc() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

		if info.FullMethod == "/UserSvc/LoginUser" ||
			info.FullMethod == "/UserSvc/RegisterUser" {

			return handler(ctx, req)
		}

		token, err := hashes.ExtractUserInfo(ctx)
		if err != nil {

		}

		uncheckedID, err := hashes.ExtractUserID(token)
		if err != nil {

		}

		salt := in.service.GetTokenSalt(ctx, uncheckedID)
		userID, err := hashes.CheckToken(token, salt)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, ErrWrongAuth)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Internal, "can't extract value from context")
		}
		md.Set(config.TokenHeader, userID)

		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
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
