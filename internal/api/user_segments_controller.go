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

// getUserSegments godoc
//
//	@Summary		Get user's segments
//	@Description	get active user's segments
//	@Tags			user_segments
//	@Param			user_id	query	int	true	"user's identificator"
//	@Produce		json
//	@Success		200	{object}	domain.GetUserSegmentsResponse
//	@Failure		400	{object}	domain.CommonResponse
//	@Failure		404	{object}	domain.CommonResponse
//	@Failure		500	{object}	domain.CommonResponse
//	@Router			/user_segments/ [get]
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

// updateUserSegments godoc
//
//	@Summary		Update user's segments
//	@Description	add and remove user to/from specified segments
//	@Tags			user_segments
//	@Accept			json
//	@Param			user_id			body	int			true	"user's identificator"
//	@Param			add_segments	body	[]string	false	"segments to add"
//	@Param			remove_segments	body	[]string	false	"segments to remove"
//	@Param			remove_at		body	string		false	"user will automatically removed from assigned segments at this date" Format(date-time)
//	@Produce		json
//	@Success		200	{object}	domain.CommonResponse
//	@Failure		400	{object}	domain.CommonResponse
//	@Failure		500	{object}	domain.CommonResponse
//	@Router			/user_segments/ [post]
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
