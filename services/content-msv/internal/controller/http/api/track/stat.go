package track_ctrl

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/services/content-msv/internal/domain/entity"
	stats "github.com/hahaclassic/orpheon/services/content-msv/internal/domain/usecase/track"
	ctxclaims "github.com/hahaclassic/orpheon/services/content-msv/internal/providers/http/ctxclaims"
)

type StatController struct {
	statService stats.ListeningStatService
}

func NewStatController(statService stats.ListeningStatService) *StatController {
	return &StatController{statService: statService}
}

func (c *StatController) UpdateStat(ctx *gin.Context) {
	trackIDStr := ctx.Param("id")
	trackID, err := uuid.Parse(trackIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var event entity.ListeningEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := ctxclaims.GetClaims(ctx)
	if claims != nil {
		event.UserID = claims.UserID
	}

	event.TrackID = trackID

	if err := c.statService.UpdateStat(ctx.Request.Context(), &event); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
