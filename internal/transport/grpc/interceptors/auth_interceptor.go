package interceptors

import (
	"context"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/hashes"
	"passkeeper/internal/entities/myerrors"
	"passkeeper/internal/usecase/srv/usersUC"
)

// AuthInterceptor using for check user auth
type AuthInterceptor struct {
	service usersUC.UserUsecaseInf
	grpc.UnaryServerInterceptor
}

// AuthFunc check user token in each request
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
			return nil, status.Errorf(codes.Unauthenticated, myerrors.ErrWrongAuth)
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

// NewAuthInterceptor is a constructor for AuthInterceptor
func NewAuthInterceptor(service usersUC.UserUsecaseInf) *AuthInterceptor {
	return &AuthInterceptor{service: service}
}
