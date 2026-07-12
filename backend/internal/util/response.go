package util

import "github.com/labstack/echo/v4"

type Envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type ErrEnvelope struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func Success(c echo.Context, status int, data any) error {
	return c.JSON(status, Envelope{Success: true, Data: data})
}

func SuccessMsg(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, Envelope{Success: true, Message: message, Data: data})
}

func Fail(c echo.Context, status int, message string) error {
	return c.JSON(status, ErrEnvelope{Success: false, Error: message})
}

func FailErr(c echo.Context, status int, err error) error {
	return c.JSON(status, ErrEnvelope{Success: false, Error: err.Error()})
}