package security

import (
	"errors"
	"fmt"
	"github.com/NidzamuddinMuzakki/go-libs/env"
	"github.com/NidzamuddinMuzakki/go-libs/log"
	"strings"

	"github.com/golang-jwt/jwt"
)

var (
	APPLICATION_NAME = env.String("MainSetup.ServiceName", "")
	// access_token_expiry = time.Duration(env.Int("Jwt.AccessToken", 60)) * time.Second
	JWT_SIGNING_METHOD = jwt.SigningMethodHS512
	JWT_SIGNATURE_KEY  = []byte(env.String("Jwt.SignatureKey", ""))
)

type Claims struct {
	jwt.StandardClaims
	JWTPayload
}

type JWTPayload struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ResponseToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

type JwtToken struct {
	Logger log.ILogging
}

type IJwtToken interface {
	ExtractToken(traceID string, header string) (token string, err error)
	ParseToken(traceID string, token string) (res *jwt.Token, err error)
	// GenerateToken(traceID string, data JWTPayload, refreshToken string) (res ResponseToken, err error)
	CheckRefreshToken(traceID string, refreshToken string) (res jwt.MapClaims, isOk bool, err error)
}

func NewJwtUtils(skipCaller int) IJwtToken {
	logger := log.NewLogging(skipCaller)
	return &JwtToken{
		Logger: logger,
	}
}

// func (lib *JwtToken) GenerateToken(traceID string, data JWTPayload, refreshToken string) (res ResponseToken, err error) {
// 	claims := Claims{
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(access_token_expiry).Unix(),
// 			Issuer:    APPLICATION_NAME,
// 		},
// 		JWTPayload: data,
// 	}

// 	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)

// 	accessToken, err := token.SignedString(JWT_SIGNATURE_KEY)
// 	if err != nil {
// 		lib.Logger.Error("GenerateToken()", "Error: %v", err)
// 		return res, err
// 	}

// 	if refreshToken == "" {
// 		reToken := jwt.New(JWT_SIGNING_METHOD)
// 		rtClaims := reToken.Claims.(jwt.MapClaims)
// 		rtClaims["id"] = data.ID
// 		rtClaims["username"] = data.Username

// 		refreshToken, err = reToken.SignedString(JWT_SIGNATURE_KEY)
// 		if err != nil {
// 			lib.Logger.Error(traceID, err.Error(), data)
// 			return res, err
// 		}
// 	}

// 	return ResponseToken{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 		Exp:          claims.ExpiresAt,
// 	}, nil
// }

func (lib *JwtToken) ExtractToken(traceID string, header string) (string, error) {

	if header == "" {
		lib.Logger.Error(traceID, "Checking Auth", "Error: Token is empty")
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, "Bearer ")
	if len(jwtToken) != 2 {
		lib.Logger.Error(traceID, "incorrect formatted header", "Error: Token is invalid")
		return "", errors.New("ncorrect formatted header")
	}

	return jwtToken[1], nil
}

func (lib *JwtToken) ParseToken(traceID string, token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if method, isOk := token.Method.(*jwt.SigningMethodHMAC); !isOk || method != JWT_SIGNING_METHOD {
			lib.Logger.Error(traceID, "bad signed method received", "Error: Token is invalid")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		lib.Logger.Error(traceID, "bad token received", token)
	}
	return parsedToken, nil
}

func (lib *JwtToken) CheckRefreshToken(traceID string, refreshToken string) (res jwt.MapClaims, isOk bool, err error) {
	token, err := lib.ParseToken(traceID, refreshToken)
	if err != nil {
		lib.Logger.Error(traceID, "CheckRefreshToken()", fmt.Sprintf("Error: %v", err))
		return res, isOk, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		lib.Logger.Error(traceID, "CheckRefreshToken()", fmt.Sprintf("Error: %v", err))
		return res, isOk, err
	}
	return claims, ok, err
}
