package responses

import (
	"github.com/labstack/echo/v4"
)

type UserResponse struct {
	Data    *echo.Map `json:"data"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
}
