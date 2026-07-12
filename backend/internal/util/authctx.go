package util

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var ErrUnauthenticated = errors.New("unauthenticated")

const (
	ctxUserID = "user_id"
	ctxRole   = "user_role"
	ctxDeptID = "department_id"
)

func UserID(c echo.Context) (int64, error) {
	v, ok := c.Get(ctxUserID).(int64)
	if !ok {
		return 0, ErrUnauthenticated
	}
	return v, nil
}

func Role(c echo.Context) (string, error) {
	v, ok := c.Get(ctxRole).(string)
	if !ok {
		return "", ErrUnauthenticated
	}
	return v, nil
}

func DepartmentID(c echo.Context) (int64, bool) {
	v, ok := c.Get(ctxDeptID).(int64)
	return v, ok
}