package util

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// retrieve JWT key from .env file
var privateKey = []byte(viper.GetString("jwt.private_key"))

// generate JWT token
func GenerateJWT(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.RoleID,
		"iat":      time.Now().Unix(),
		"eat":      time.Now().Add(time.Second * time.Duration(viper.GetInt("jwt.token_ttl"))).Unix(),
	})

	return token.SignedString(privateKey)
}

// validate JWT token
func ValidateJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	_, ok := token.Claims.(jwt.MapClaims)
	if ok {
		return nil
	}

	return errors.New("invalid token provided")
}

// validate Admin role
func ValidateAdminRoleJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userRole := uint(claims["role"].(float64))
	if userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid admin token provided")
}

// validate Viewer role
func ValidateViewerRoleJWT(c echo.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userRole := uint(claims["role"].(float64))
	if userRole == model.ViewerRole || userRole == model.AdminRole {
		return nil
	}

	return errors.New("invalid viewer or admin token provided")
}

// GetUserIDOnToken return user id
func GetUserIDOnToken(c echo.Context) (string, error) {
	token, err := getToken(c)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", errors.New("invalid id claims")
	}

	return id, nil
}

// GetOwnerFromToken returns the owner details from the JWT token
func GetOwnerFromToken(c echo.Context) (model.Owner, error) {
	token, err := getToken(c)
	if err != nil {
		return model.Owner{}, err
	}

	if !token.Valid {
		return model.Owner{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return model.Owner{}, errors.New("invalid token claims")
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return model.Owner{}, errors.New("invalid id claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return model.Owner{}, errors.New("invalid username claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return model.Owner{}, errors.New("invalid email claims")
	}

	role, ok := claims["role"].(float64)
	if !ok {
		return model.Owner{}, errors.New("invalid role claims")
	}

	return model.Owner{
		ID:       int64(id),
		Username: username,
		Email:    email,
		RoleID:   uint(role),
	}, nil
}

// GetOwnerFromCtx returns the owner details from the context
func GetOwnerFromCtx(ctx context.Context) *model.Owner {
	owner, ok := ctx.Value("user").(model.Owner)
	if ok {
		return &owner
	}

	return nil
}

// check token validity
func getToken(context echo.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})

	return token, err
}

// extract token from request Authorization header
func getTokenFromRequest(c echo.Context) string {
	bearerToken := c.Request().Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
