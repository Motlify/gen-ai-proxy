package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"gen-ai-proxy/src/database"
	"gen-ai-proxy/src/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	userContextKey = "userID"
)


func APIKeyAuthMiddleware(db *database.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				authHeader := c.Request().Header.Get("Authorization")
				if authHeader == "" {
					return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authentication header is missing"})
				}
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid Authorization header format"})
				}
				apiKey = parts[1]
			}

			hash := sha256.Sum256([]byte(apiKey))
			apiKeyHash := hex.EncodeToString(hash[:])

			apiKeyRecord, err := db.GetAPIKeyByHash(context.Background(), apiKeyHash)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid API Key"})
			}

			// Update last_used_at timestamp
			err = db.UpdateAPIKeyLastUsed(context.Background(), apiKeyRecord.ID)
			if err != nil {
				// Log the error but don't block the request
				c.Logger().Errorf("Failed to update API key last_used_at: %v", err)
			}

			c.Set(userContextKey, apiKeyRecord.UserID)

			return next(c)
		}
	}
}


func JWTAuthMiddleware(jwtSecret []byte) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Authorization header is missing"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid Authorization header format"})
			}

			tokenString := parts[1]

			token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid or expired token"})
			}

			claims, ok := token.Claims.(*models.JwtCustomClaims)
			if !ok || !token.Valid {
				return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid token claims"})
			}

			c.Set(userContextKey, claims.UserID)

			return next(c)
		}
	}
}

func GetUserIDFromContext(c echo.Context) (pgtype.UUID, error) {
	userID, ok := c.Get(userContextKey).(pgtype.UUID)
	if !ok {
		return pgtype.UUID{}, echo.NewHTTPError(http.StatusInternalServerError, "User ID not found in context")
	}
	return userID, nil
}
