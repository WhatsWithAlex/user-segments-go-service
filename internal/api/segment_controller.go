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
