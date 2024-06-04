package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/w1ns3c/passkeeper/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/w1ns3c/go-examples/crypto"
	"golang.org/x/crypto/bcrypt"

	"github.com/w1ns3c/passkeeper/internal/entities"
	"github.com/w1ns3c/passkeeper/internal/storage"
)

var (
	ErrGetUser       = fmt.Errorf("can't get user by ID")
	ErrGenHash       = fmt.Errorf("can't generate hash of password")
	ErrWrongPassword = fmt.Errorf("old password is wrong")
	ErrRepassNotSame = fmt.Errorf("new pass and repeat not the same")

	ErrChallengeGen  = fmt.Errorf("can't generate challenge")
	ErrChallengeLife = fmt.Errorf("challenge too old")

	ErrWrongAuth    = fmt.Errorf("wrong user/password")
	ErrInvalidToken = fmt.Errorf("token sign is not valid")
)

type UserUsecaseInf interface {
	RegisterUser(ctx context.Context, login string, password string, rePass string) (token string, err error)

	ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error)
	GetTokenSalt() string

	ChallengeGenerate(ctx context.Context, login string) (challenge string, err error)

	LoginUser(ctx context.Context, login string, password string) (token string, err error)
}

type UserUsecase struct {
	storage       storage.UserStorage
	salt          string
	tokenLifeTime time.Duration
	userSecretLen int
	log           *zerolog.Logger
}

func (u *UserUsecase) ChallengeGenerate(ctx context.Context, login string) (challenge string, err error) {
	exist, err := u.storage.CheckUserExist(ctx, login)
	if !exist {
		u.log.Error().Err(err).Msg("requested auth for non existed user")

		return "", ErrChallengeGen
	}

	if err != nil {
		u.log.Error().Err(err).Msg("can't check user exists")

		return "", ErrChallengeGen
	}

	challenge, err = crypto.GenRandStr(config.ChallengeLen)
	if err != nil {
		u.log.Error().Err(err).Msg("can't generate rand string")

		return "", ErrChallengeGen
	}

	err = u.storage.SaveChallenge(ctx, challenge)
	if err != nil {
		u.log.Error().Err(err).Msg("can't save current user challenge")

		return "", ErrChallengeGen
	}

	return challenge, nil
}

func (u *UserUsecase) ChangePassword(ctx context.Context, userID, oldPass, newPass, reNewPass string) (err error) {
	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGetUser.Error())

		return ErrGetUser
	}

	hOld, err := GenerateHash(oldPass, u.salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGenHash.Error())

		return ErrGenHash
	}

	if hOld != user.Hash {
		u.log.Error().
			Err(ErrWrongPassword).Send()

		return ErrWrongPassword
	}

	hNew1, err := GenerateHash(newPass, u.salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGenHash.Error())

		return ErrGenHash
	}

	hNew2, err := GenerateHash(reNewPass, u.salt)
	if err != nil {
		u.log.Error().Err(err).
			Msg(ErrGenHash.Error())

		return ErrGenHash
	}

	if hNew1 != hNew2 {
		u.log.Error().
			Err(ErrRepassNotSame).Send()

		return ErrRepassNotSame
	}

	user.Hash = hNew1
	return u.storage.SaveUser(ctx, user)

}

func (u *UserUsecase) LoginUser(ctx context.Context, login string, challengeRes string) (token string, err error) {

	user, err := u.storage.GetUserByLogin(ctx, login)
	if err != nil {
		return "", ErrWrongAuth
	}

	same := u.CompareChallenges(challengeRes, user)

	if !same {
		return "", fmt.Errorf("wrong auth for user: %s", login)
	}

	return GenerateToken(user.ID, u.salt, u.tokenLifeTime)
}

func (u *UserUsecase) CompareChallenges(actualChallenge string, user *entities.User) bool {
	if time.Since(user.ChallengeTime).Minutes() > config.ChallengeLifeTime {
		u.log.Error().Err(ErrChallengeLife).Msg("can't generate challenge")

		return false
	}

	decChallenge, err := crypto.DecryptAES([]byte(actualChallenge), []byte(user.Hash))
	if err != nil {
		u.log.Error().Err(err).Msg("can't decrypt challenge")

		return false
	}

	return string(decChallenge) == actualChallenge
}

func (u *UserUsecase) RegisterUser(ctx context.Context, login string, password string, rePass string) (token string, err error) {
	if password != rePass {
		return "", entities.ErrPassNotTheSame
	}

	// checking login free
	exist, err := u.storage.CheckUserExist(ctx, login)
	if !errors.Is(err, entities.ErrUserNotFound) || exist {
		return "", fmt.Errorf("user is already exist:%v", err)
	}

	hash, err := GenerateHash(password, u.salt)
	if err != nil {
		return "", fmt.Errorf("can't generate hash of password: %v", err)
	}

	m := md5.Sum([]byte(hash))
	id := GenerateID(hex.EncodeToString(m[:]), u.salt)

	secret, err := GenerateSecret(u.userSecretLen)
	if err != nil {
		return "", fmt.Errorf("can't generate secret for user: %v", err)
	}

	user := &entities.User{
		ID:    id,
		Login: login,
		//Credential: password,
		Hash:   hash,
		Secret: secret,
	}

	err = u.storage.SaveUser(ctx, user)
	if err != nil {
		return "", err
	}

	token, err = GenerateToken(user.ID, u.salt, u.tokenLifeTime)
	if err != nil {
		return "", fmt.Errorf("can't generate user token: %v", err)
	}

	return token, nil
}

func GenerateHash(password, salt string) (hash string, err error) {
	password = fmt.Sprintf("%s-%s.%s.%s", string(salt), string(password), string(password), string(salt))
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func GenerateID(secret, salt string) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s.%s.%s", salt, secret, salt)))

	return hex.EncodeToString(hash[:])
}

func (u *UserUsecase) GetTokenSalt() string {
	return u.salt
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func GenerateSecret(secretLen int) (secret string, err error) {
	sl, err := crypto.GenRandSlice(secretLen)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(sl), nil
}

func GenerateToken(userid string, secret string, lifetime time.Duration) (token string, err error) {
	tokenLife := time.Now().Add(lifetime)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenLife),
		},
		UserID: userid,
	})
	token, err = jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func CheckToken(tokenstr, secret string) (userID string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenstr, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	// return user ID in readable format
	return claims.UserID, nil
}
