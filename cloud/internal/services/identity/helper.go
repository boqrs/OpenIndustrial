package identity

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/boqrs/OpenIndustrial/cloud/internal/persistence/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// =====================================================
// Password
// =====================================================

const passwordCost = bcrypt.DefaultCost

// HashPassword
//
// bcrypt生成密码hash
func HashPassword(
	password string,
) (
	string,
	error,
) {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		passwordCost,
	)

	if err != nil {
		return "",
			err
	}

	return string(hash), nil
}

// VerifyPassword
//
// 校验密码
func VerifyPassword(
	hash string,
	password string,
) bool {

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}

// =====================================================
// JWT
// =====================================================

type IdentityClaims struct {
	UserID uuid.UUID `json:"user_id"`

	TenantID uint `json:"tenant_id"`

	RoleID uint `json:"role_id"`

	TokenType string `json:"token_type"`

	jwt.RegisteredClaims
}

const (
	TokenTypeAccess = "access"

	TokenTypeRefresh = "refresh"
)

// GenerateAccessToken
//
// 生成访问token
func GenerateAccessToken(
	secret string,
	user *model.User,
	expire time.Duration,
) (
	string,
	error,
) {

	claims := IdentityClaims{

		UserID: user.UUID,

		TenantID: user.TenantID,

		RoleID: user.RoleID,

		TokenType: TokenTypeAccess,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(
				time.Now().
					Add(expire),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ID: uuid.New().
				String(),
		},
	}

	return signJWT(
		secret,
		claims,
	)
}

// GenerateRefreshToken
//
// refresh token
func GenerateRefreshToken(
	secret string,
	user *model.User,
	expire time.Duration,
) (
	string,
	error,
) {

	claims := IdentityClaims{

		UserID: user.UUID,

		TenantID: user.TenantID,

		RoleID: user.RoleID,

		TokenType: TokenTypeRefresh,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(
				time.Now().
					Add(expire),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),

			ID: uuid.New().
				String(),
		},
	}

	return signJWT(
		secret,
		claims,
	)

}

func signJWT(
	secret string,
	claims IdentityClaims,
) (
	string,
	error,
) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(
		[]byte(secret),
	)
}

// ParseRefreshToken
//
// 验证refresh token
func ParseRefreshToken(
	secret string,
	tokenString string,
) (
	*IdentityClaims,
	error,
) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&IdentityClaims{},
		func(token *jwt.Token) (
			interface{},
			error,
		) {

			if token.Method != jwt.SigningMethodHS256 {
				return nil,
					errors.New(
						"invalid signing method",
					)
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok :=
		token.Claims.(*IdentityClaims)

	if !ok ||
		!token.Valid {

		return nil,
			errors.New(
				"invalid token",
			)
	}

	if claims.TokenType != TokenTypeRefresh {

		return nil,
			errors.New(
				"not refresh token",
			)
	}

	if claims.UserID == uuid.Nil ||
		claims.TenantID == 0 ||
		claims.RoleID == 0 {

		return nil,
			errors.New(
				"invalid identity claims",
			)
	}

	return claims, nil
}

// =====================================================
// Invitation Token
// =====================================================

// GenerateInvitationToken
//
// 返回:
//
// token 明文
// hash存数据库
func GenerateInvitationToken() (
	string,
	string,
	error,
) {

	buf := make([]byte, 32)

	_, err :=
		rand.Read(buf)

	if err != nil {

		return "",
			"",
			err
	}

	token := hex.EncodeToString(
		buf,
	)

	hash, err := HashToken(
		token,
	)

	if err != nil {
		return "",
			"",
			err
	}

	return token, hash, nil
}

// HashToken
//
// Invitation Token 使用 SHA-256 做确定性 hash。
//
// 与密码不同，Invitation Token 本身已经是高熵随机值，
// 因此不需要 bcrypt 的随机 salt。
// Repository 可以直接通过 hash 查询 invitation。
func HashToken(
	token string,
) (
	string,
	error,
) {

	if token == "" {
		return "",
			errors.New(
				"token is empty",
			)
	}

	digest := sha256.Sum256(
		[]byte(token),
	)

	return hex.EncodeToString(
		digest[:],
	), nil
}

// VerifyToken
//
// 校验 Invitation Token。
func VerifyToken(
	hash string,
	token string,
) bool {

	computed, err := HashToken(
		token,
	)

	if err != nil {
		return false
	}

	return computed == hash
}
