package api

import (
	"database/sql"
	"net/http"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/gin-gonic/gin"
)

type userSegmentController struct {
	userSegmentService domain.UserSegmentService
	env                *env.Env
}

func (usc *userSegmentController) getUserSegments(c *gin.Context) {
	var request domain.GetUserSegmentsRequest

	err := c.BindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	segments, err := usc.userSegmentService.GetActiveUserSegments(c, request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.GetUserSegmentsResponse{
		Segments: segments,
	})
}

func (usc *userSegmentController) updateUserSegments(c *gin.Context) {
	var request domain.UpdateUserSegmentsRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	var removeAt sql.NullTime
	if (request.RemoveAt != nil) && removeAt.Scan(*request.RemoveAt) != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse("remove_at field"))
		return
	}

	err = usc.userSegmentService.UpdateUserSegments(
		c,
		request.UserID,
		request.AddSegments,
		request.RemoveSegments,
		removeAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse("user's segments updated successfully"))
}
