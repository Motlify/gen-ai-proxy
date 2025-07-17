package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"gen-ai-proxy/src"
	"gen-ai-proxy/src/database"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"` // changed from int64 to string
	Username string `json:"username"`
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with a username and password.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/register [post]
// @security []
func (s *Service) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
	}

	// Check if user already exists
	_, err = s.db.GetUserByUsername(c.Request().Context(), req.Username)
	if err == nil {
		log.Printf("User with username %s already exists", req.Username)
		return c.JSON(http.StatusConflict, ErrorResponse{Error: "User with this username already exists"})
	} else if errors.Is(err, sql.ErrNoRows) {
		// User does not exist, proceed to create
	} else {
		log.Printf("Error checking for existing user: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	user, err := s.db.CreateUser(c.Request().Context(), database.CreateUserParams{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		log.Printf("Error creating user in DB: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create user"})
	}

	var uuidStr string
	if user.ID.Valid {
		uuidStr = uuid.UUID(user.ID.Bytes).String()
	} else {
		log.Printf("User ID is not valid")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "invalid user ID"})
	}

	return c.JSON(http.StatusCreated, UserResponse{
		ID:       uuidStr,
		Username: user.Username,
	})
}

// Login godoc
// @Summary Log in a user
// @Description Log in a user with a username and password to get a JWT token.
// @Tags Users
// @Accept json,x-www-form-urlencoded
// @Produce json
// @Param user body LoginRequest true "User login details"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/login [post]
// @security []
func (s *Service) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	user, err := s.db.GetUserByUsername(c.Request().Context(), req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		}
		log.Printf("Error getting user by username during login: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
	}

	token, err := internal.CreateToken(user.ID, []byte(s.cfg.JWTSecret))
	if err != nil {
		log.Printf("Error generating token during login: %v", err)
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
	}

	log.Printf("Login successful for user %s, token generated.", req.Username)
	return c.JSON(http.StatusOK, LoginResponse{AccessToken: token})
}
