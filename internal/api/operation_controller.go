package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/gin-gonic/gin"
)

type operationController struct {
	operationService domain.OperationService
	env              *env.Env
}

func (oc *operationController) getOperations(c *gin.Context) {
	var request domain.GetOperationsRequest

	err := c.BindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	var fromTime, toTime sql.NullTime

	err = fromTime.Scan(time.Date(request.Year, time.Month(request.Month), 0, 0, 0, 0, 0, time.Local))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse("year field"))
		return
	}
	err = toTime.Scan(time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse("month field"))
		return
	}

	fileName, err := oc.operationService.GetUserOperations(c, request.UserID, fromTime, toTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse(err.Error()))
		return
	}
	fileUrl := "http://" + c.Request.Host + "/static/" + fileName

	c.JSON(http.StatusOK, domain.GetOperationsResponse{
		FileURL: fileUrl,
	})
}
