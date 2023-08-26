package response

import (
	"little_mangamee/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccesResponse(c *gin.Context, bodyPayload interface{}) {

	returnData := gin.H{}
	returnData["result_code"] = "00"
	returnData["message"] = "success"
	if bodyPayload != nil {
		returnData["data"] = bodyPayload
	}
	c.JSON(http.StatusOK, returnData)
	return
}

func ErrorResponse(c *gin.Context, err error, message interface{}) {

	dataErr, ok := utils.ERR_MAP[err]
	if !ok {
		dataErr = utils.ERR_MAP[utils.ERR_INTERNAL_SERVER]
		c.JSON(dataErr.HttpCode, gin.H{
			"result_code": dataErr.ErrorCode,
			"message":     dataErr.Message,
		})
		return
	}

	if message == nil {
		c.JSON(dataErr.HttpCode, gin.H{
			"result_code": dataErr.ErrorCode,
			"message":     dataErr.Message,
		})
		return
	}

	c.JSON(dataErr.HttpCode, gin.H{
		"result_code": dataErr.ErrorCode,
		"message":     message,
	})

	return
}
