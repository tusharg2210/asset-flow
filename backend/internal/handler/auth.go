package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/config"
	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type AuthHandler struct {
	users  *repository.UserRepository
	jwtCfg config.JWTConfig
}

func NewAuthHandler(users *repository.UserRepository, jwtCfg config.JWTConfig) *AuthHandler {
	return &AuthHandler{users: users, jwtCfg: jwtCfg}
}

type signupRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var req signupRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return util.Fail(c, http.StatusInternalServerError, "could not process password")
	}

	u := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     model.RoleEmployee,
		Status:   model.StatusActive,
	}

	if err := h.users.Create(c.Request().Context(), u); err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			return util.Fail(c, http.StatusConflict, "email already registered")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	u.Password = ""
	return util.SuccessMsg(c, http.StatusCreated, "User registered successfully", echo.Map{"user": u})
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	u, err := h.users.GetByEmail(c.Request().Context(), req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return util.Fail(c, http.StatusUnauthorized, "invalid email or password")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	if u.Status != model.StatusActive {
		return util.Fail(c, http.StatusForbidden, "account is inactive")
	}
	if !util.CheckPassword(req.Password, u.Password) {
		return util.Fail(c, http.StatusUnauthorized, "invalid email or password")
	}

	accessToken, err := util.GenerateAccessToken(h.jwtCfg.JWTSecret, u.ID, u.Role, h.jwtCfg.AccessTokenTTL)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	refreshToken, err := util.GenerateRefreshToken(h.jwtCfg.JWTSecret, u.ID, u.Role, h.jwtCfg.RefreshTokenTTL)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	h.setRefreshCookie(c, refreshToken)

	u.Password = ""
	return util.Success(c, http.StatusOK, echo.Map{
		"token": accessToken,
		"user":  echo.Map{"id": u.ID, "name": u.Name, "role": u.Role},
	})
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie(h.jwtCfg.RefreshCookieKey)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "missing refresh token")
	}

	claims, err := util.ParseRefreshToken(h.jwtCfg.JWTSecret, cookie.Value)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "invalid or expired refresh token")
	}

	u, err := h.users.GetByID(c.Request().Context(), claims.UserID)
	if err != nil || u.Status != model.StatusActive {
		return util.Fail(c, http.StatusUnauthorized, "account no longer active")
	}

	accessToken, err := util.GenerateAccessToken(h.jwtCfg.JWTSecret, u.ID, u.Role, h.jwtCfg.AccessTokenTTL)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	newRefresh, err := util.GenerateRefreshToken(h.jwtCfg.JWTSecret, u.ID, u.Role, h.jwtCfg.RefreshTokenTTL)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	h.setRefreshCookie(c, newRefresh)

	return util.Success(c, http.StatusOK, echo.Map{"token": accessToken})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     h.jwtCfg.RefreshCookieKey,
		Value:    "",
		HttpOnly: true,
		Path:     "/api/auth",
		MaxAge:   -1,
	})
	return util.SuccessMsg(c, http.StatusOK, "logged out", nil)
}

func (h *AuthHandler) setRefreshCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     h.jwtCfg.RefreshCookieKey,
		Value:    token,
		HttpOnly: true,
		Secure:   false, // set false only in local dev over http
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth",
		Expires:  time.Now().Add(h.jwtCfg.RefreshTokenTTL),
	})
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID, err := util.UserID(c)
	if err != nil {
		return util.Fail(c, http.StatusUnauthorized, "unauthenticated")
	}

	u, err := h.users.GetByID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return util.Fail(c, http.StatusNotFound, "user not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}

	u.Password = ""
	return util.Success(c, http.StatusOK, u)
}
