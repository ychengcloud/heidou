package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"{{ . }}/internal/gen/models"
)

// 定义错误
var (
	ErrInvalidToken = errors.New("invalid token")
)

const defaultKey = "YOUCHENG"

var DefaultOptions = Options{
	TokenHeader:   "Authorization",
	TokenType:     "Bearer",
	Expired:       7200,
	signingMethod: jwt.SigningMethodHS512,
	SigningKey:    defaultKey,
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(defaultKey), nil
	},
	Issuer: "ycheng.pro",
}

// TokenInfo 令牌信息
type TokenInfo interface {
	// 获取访问令牌
	GetAccessToken() string
	// 获取令牌类型
	GetTokenType() string
	// 获取令牌到期时间戳
	GetExpiresAt() int64
	// JSON编码
	EncodeToJSON() ([]byte, error)
}

// tokenInfo 令牌信息
type tokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 令牌到期时间
}

func (t *tokenInfo) GetAccessToken() string {
	return t.AccessToken
}

func (t *tokenInfo) GetTokenType() string {
	return t.TokenType
}

func (t *tokenInfo) GetExpiresAt() int64 {
	return t.ExpiresAt
}

func (t *tokenInfo) EncodeToJSON() ([]byte, error) {
	return json.Marshal(t)
}

type Options struct {
	TokenHeader   string
	SigningMethod string
	Expired       int
	TokenType     string
	Issuer        string
	ClaimsKey     string
	SigningKey    string

	Logger *zap.Logger

	keyfunc       jwt.Keyfunc
	signingMethod jwt.SigningMethod
}

type Enforcer interface {
	Enforce(vals ...interface{}) (bool, error)
}

// JWTAuth jwt认证
type JWTAuth struct {
	*Options

	signingKey interface{}
	Enforcer   Enforcer
}

// New 创建认证实例
func New(enforcer Enforcer, opts *Options) (*JWTAuth, error) {
	auth := &JWTAuth{
		Options:  opts,
		Enforcer: enforcer,
	}
	auth.keyfunc = func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			auth.Logger.Info("initauth err")
			return nil, ErrInvalidToken
		}
		auth.Logger.Info("initauth :", zap.String("signingKey", opts.SigningKey))

		return []byte(auth.SigningKey), nil
	}
	auth.signingKey = []byte(auth.SigningKey)

	switch opts.SigningMethod {
	case "HS256":
		auth.signingMethod = jwt.SigningMethodHS256
	case "HS384":
		auth.signingMethod = jwt.SigningMethodHS384
	case "HS512":
		auth.signingMethod = jwt.SigningMethodHS512
	}

	return auth, nil
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(claims jwt.MapClaims) (TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(a.Expired) * time.Second).Unix()

	claims["iss"] = a.Issuer
	claims["iat"] = now.Unix()
	claims["exp"] = expiresAt
	claims["nbf"] = now.Unix()

	token := jwt.NewWithClaims(a.signingMethod, claims)

	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return nil, err
	}

	tokenInfo := &tokenInfo{
		ExpiresAt:   expiresAt,
		TokenType:   a.TokenType,
		AccessToken: tokenString,
	}

	return tokenInfo, nil
}

// 解析令牌
func (a *JWTAuth) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, a.keyfunc)
	if !token.Valid {
		return nil, err
	}

	if token != nil {
		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			return *claims, nil
		}
	}

	return nil, err
}

// ParseToken 解析ID
func (a *JWTAuth) ParseToken(tokenString string) (jwt.MapClaims, error) {
	return a.parseToken(tokenString)
}

func (a *JWTAuth) Authentication(c *gin.Context) (claims jwt.MapClaims, err error) {
	return a.authentication(c)
}

func (a *JWTAuth) authentication(c *gin.Context) (claims jwt.MapClaims, err error) {
	var token string
	auth := c.GetHeader(a.TokenHeader)
	prefix := a.TokenType
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = strings.Trim(auth[len(prefix):], " ")
	}

	if token == "" {
		err = models.ErrUnauthorized
		return nil, err
	}

	claims, err = a.ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			err = models.ErrTokenExpried
		default:
			err = models.ErrUnauthorized
		}
	}
	log.Println("jwt claims:", auth, claims, err)

	return claims, err
}

func (a *JWTAuth) Authorization(c *gin.Context) (bool, error) {
	claims, err := a.authentication(c)
	if err != nil {
		return false, err
	}

	return a.Enforcer.Enforce(claims["username"], c.Request.URL.Path, c.Request.Method)
}

var ProviderSet = wire.NewSet(New)
