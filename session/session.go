package session

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"eikva.ru/eikva/database"
	"eikva.ru/eikva/models"
	"github.com/golang-jwt/jwt/v5"
)

type EikvaClaims struct {
	UserUUID  string `json:"user_id"`
	UserLogin string `json:"user_login"`
	TokenID   string `json:"token_id"`
	jwt.RegisteredClaims
}

type EikvaSessionTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var parser = jwt.NewParser(jwt.WithValidMethods([]string{"HS512"}))
var jwtSecret = os.Getenv("JWT_SECRET")

func CreateToken(user *models.User, tokenID string, ttl time.Duration) string {
	now := time.Now()
	tokenCalims := EikvaClaims{
		UserUUID:  user.UUID,
		UserLogin: user.Login,
		TokenID:   tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Eikva",
			Subject:   user.UUID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, tokenCalims)
	result, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		panic("token signature fail")
	}

	return result
}

func CreateSessionTokens(user *models.User) EikvaSessionTokens {
	return EikvaSessionTokens{
		AccessToken:  CreateToken(user, user.AccessTokenID.String, time.Hour),
		RefreshToken: CreateToken(user, user.RefreshTokenID.String, time.Duration(time.Hour*168)),
	}
}

type ErrNotMatchingId struct {
	Message string
}

func (e ErrNotMatchingId) Error() string {
	return e.Message
}

func GetTokenClaims(token string) (*EikvaClaims, *jwt.Token, error) {
	claims := &EikvaClaims{}
	parsed, err := parser.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	return claims, parsed, err
}

type TokenType int

const (
	TokenTypeAccess TokenType = iota
	TokenTypeRefresh
)

func ValidateSessionTokenAndGetUser(token string, tokenType TokenType) (*models.User, error) {
	claims, parsed, err := GetTokenClaims(token)

	if err != nil {
		return nil, err
	}

	if parsed.Valid {
		user, err := database.GetExistingUserByUUID(claims.UserUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, errors.New("Сессия не найдена")
			}
			return nil, err
		}
		var tokenID sql.NullString
		switch tokenType {
		case TokenTypeAccess:
			tokenID = user.AccessTokenID
			break;
		case TokenTypeRefresh:
			tokenID = user.RefreshTokenID
			break;
		default:
			return nil, errors.New("Не верный тип токена")
		}

		if !tokenID.Valid || tokenID.String == "" {
			return nil, errors.New("placeholder error: no AccessTokenID in DB")
		}

		if claims.TokenID == "" || claims.TokenID != tokenID.String {
			return nil, ErrNotMatchingId{
				Message: "Не совпадает ID токена",
			}
		}

		return user, nil
	}

	return nil, errors.New("placeholder !parsed.Valid")
}
