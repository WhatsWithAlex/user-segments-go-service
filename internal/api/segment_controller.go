package api

import (
	"net/http"

	"github.com/WhatsWithAlex/user-segments-go-service/internal/domain"
	"github.com/WhatsWithAlex/user-segments-go-service/internal/env"
	"github.com/gin-gonic/gin"
)

type segmentController struct {
	segmentService domain.SegmentService
	env            *env.Env
}

// createSegment godoc
//
//	@Summary		Create segment
//	@Description	create segment with given unique name (slug)
//	@Tags			segments
//	@Accept			json
//	@Param			slug		body	string	true	"segment name"
//	@Param			probability	body	number	false	"probability of auto assignment"
//	@Produce		json
//	@Success		200	{object}	domain.CommonResponse
//	@Failure		400	{object}	domain.CommonResponse
//	@Failure		500	{object}	domain.CommonResponse
//	@Router			/segments/ [post]
func (sc *segmentController) createSegment(c *gin.Context) {
	var request domain.CreateSegmentRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	err = sc.segmentService.CreateSegment(c, request.Slug, request.Probability)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse("segment created successfully"))
}

// deleteSegment godoc
//
//	@Summary		Delete segment
//	@Description	delete segment by slug
//	@Tags			segments
//	@Param			slug	query	string	true	"segment slug"
//	@Produce		json
//	@Success		200	{object}	domain.CommonResponse
//	@Failure		400	{object}	domain.CommonResponse
//	@Failure		404	{object}	domain.CommonResponse
//	@Failure		500	{object}	domain.CommonResponse
//	@Router			/segments/ [delete]
func (sc *segmentController) deleteSegment(c *gin.Context) {
	var request domain.DeleteSegmentRequest

	err := c.BindQuery(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse(err.Error()))
		return
	}

	err = sc.segmentService.DeleteSegment(c, request.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse("segment deleted successfully"))
}
