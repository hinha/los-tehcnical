package response

import "github.com/labstack/echo/v4"

// Structure Response Data
type Response struct {
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Code    int         `json:"code" example:"200"`
}

func DefaultResponse(c echo.Context, message string, data, errors interface{}, code int) error {
	var response Response
	response.Data = data
	response.Code = code
	response.Errors = errors
	response.Message = message

	return c.JSON(code, response)
}
