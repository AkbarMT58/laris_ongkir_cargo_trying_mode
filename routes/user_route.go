package routes

import (
	"github.com/AdonisVillanueva/golang-echo-mongo-api/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {

	e.GET("/getongkir", controllers.GetAllOngkir)

}
