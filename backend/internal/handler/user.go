package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"asset-flow/internal/model"
	"asset-flow/internal/repository"
	"asset-flow/internal/util"
)

type UserHandler struct {
	users *repository.UserRepository
}

func NewUserHandler(users *repository.UserRepository) *UserHandler {
	return &UserHandler{users}
}

func (h *UserHandler) List(c echo.Context) error {
	var departmentID *int64
	if v := c.QueryParam("department_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return util.Fail(c, http.StatusBadRequest, "invalid department_id")
		}
		departmentID = &id
	}

	users, err := h.users.List(c.Request().Context(), departmentID)
	if err != nil {
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	return util.Success(c, http.StatusOK, users)
}


type updateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=DEPARTMENT_HEAD ASSET_MANAGER"`
}


func (h *UserHandler) UpdateRole(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid user id")
	}

	var req updateRoleRequest
	if err := c.Bind(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return util.Fail(c, http.StatusBadRequest, err.Error())
	}

	u, err := h.users.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return util.Fail(c, http.StatusNotFound, "user not found")
		}
		return util.FailErr(c, http.StatusInternalServerError, err)
	}
	if u.Role != model.RoleEmployee {
		return util.Fail(c, http.StatusConflict, "only employees can be promoted")
	}

u, err = h.users.UpdateRole(c.Request().Context(), id, req.Role)
if err != nil {
	if errors.Is(err, repository.ErrUserNotPromotable) {
		return util.Fail(c, http.StatusConflict, "user not found or not eligible for promotion")
	}
	return util.FailErr(c, http.StatusInternalServerError, err)
}
return util.Success(c, http.StatusOK, echo.Map{"id": u.ID, "name": u.Name, "role": u.Role})}