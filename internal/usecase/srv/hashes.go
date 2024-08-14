package srv

// GenerateHash func gen sha256 hash of (password with salt)
func GenerateHash(password, salt string) string {
	password = fmt.Sprintf("%s-%s.%s.%s", string(salt), string(password), string(password), string(salt))
	hash := sha256.Sum256([]byte(password))

	return fmt.Sprintf("%x", hash)
}

// GenerateCryptoHash func
func GenerateCryptoHash(password, salt string) (hash string, err error) {
	password = GenerateHash(password, salt)
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(h), nil
}

func ComparePassAndCryptoHash(password, hash string, salt string) bool {
	genHash := GenerateHash(password, salt)

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(genHash)); err != nil {
		return false
	}
	return true
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

func CheckToken(tokenStr, secret string) (userID string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims,
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
