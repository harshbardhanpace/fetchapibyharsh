package health

import (
	"net/http"
	apihelpers "space/apiHelpers"
	"space/business/health"
	"space/helpers"
	"space/loggerconfig"
	"space/models"

	"github.com/gin-gonic/gin"
)

func GetHealthStatus(ctx *gin.Context) {

	res := apihelpers.APIRes{}
	resp := health.GetHealthInfo()
	if len(resp) == 0 {
		pong := models.Pong{
			DT: helpers.GetCurrentTimeInIST(),
		}

		res.Status = true
		res.Message = "SUCCESS"
		res.Data = pong
		// utils.PrintlnLog(constants.INFO, "health check api success : ", resp)
		loggerconfig.Info("health (controller), check api success : ", resp)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// utils.PrintlnLog(constants.ERROR, "health check api error : ", resp)
	loggerconfig.Error("health (controller), check api error : ", resp)
	res.Status = false
	res.Message = "FAILURE"
	res.Data = resp
	ctx.JSON(http.StatusInternalServerError, res)
}
