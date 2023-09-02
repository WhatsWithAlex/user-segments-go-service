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

// getOperations godoc
//
//	@Summary		Get operations file link
//	@Description	get csv-file link for operations made in specified period with given user's segments
//	@Tags			operations
//	@Produce		json
//	@Param			user_id	query		integer	true	"user's identificator"
//	@Param			year	query		integer	true	"year of the starting date"		minimum(1970)	maximum(9999)
//	@Param			month	query		integer	true	"month of the starting date"	minimum(1)		maximum(12)
//	@Success		200		{object}	domain.GetOperationsResponse
//	@Failure		400		{object}	domain.CommonResponse
//	@Failure		404		{object}	domain.CommonResponse
//	@Failure		500		{object}	domain.CommonResponse
//	@Router			/operations/ [get]
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
