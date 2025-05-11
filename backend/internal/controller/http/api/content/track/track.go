package track

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/track"
)

type TrackController struct {
	trackService   track.TrackMetaService
	authMiddleware gin.HandlerFunc
}

func New(trackService track.TrackMetaService, authMiddleware gin.HandlerFunc) *TrackController {
	return &TrackController{
		trackService:   trackService,
		authMiddleware: authMiddleware,
	}
}

func (c *TrackController) RegisterRoutes(router *gin.RouterGroup) {
	tracks := router.Group("/tracks")
	{
		tracks.GET("/:id", c.GetTrack)

		protected := tracks.Group("")
		protected.Use(c.authMiddleware)
		{
			protected.POST("", c.CreateTrack)
			protected.PUT("/:id", c.UpdateTrack)
			protected.DELETE("/:id", c.DeleteTrack)
		}
	}
}

// GetTrack godoc
// @Summary Get track by ID
// @Description Get track details by its ID
// @Tags tracks
// @Produce json
// @Param id path string true "Track ID"
// @Success 200 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [get]
func (c *TrackController) GetTrack(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	track, err := c.trackService.GetTrackMeta(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	ctx.JSON(http.StatusOK, track)
}

// CreateTrack godoc
// @Summary Create a new track
// @Description Create a new track with the provided details
// @Tags tracks
// @Accept json
// @Produce json
// @Param track body entity.Track true "Track object"
// @Success 201 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks [post]
func (c *TrackController) CreateTrack(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var track entity.TrackMeta
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdTrack, err := c.trackService.CreateTrackMeta(ctx.Request.Context(), claims, &track)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to create track"})
		return
	}

	ctx.JSON(http.StatusCreated, createdTrack)
}

// UpdateTrack godoc
// @Summary Update a track
// @Description Update an existing track with new details
// @Tags tracks
// @Accept json
// @Produce json
// @Param id path string true "Track ID"
// @Param track body entity.Track true "Updated track object"
// @Success 200 {object} entity.Track
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [put]
func (c *TrackController) UpdateTrack(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	var track entity.TrackMeta
	if err := ctx.ShouldBindJSON(&track); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track.ID = id
	err = c.trackService.UpdateTrackMeta(ctx.Request.Context(), claims, &track)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to update track"})
		return
	}

	ctx.Status(http.StatusOK)
}

// DeleteTrack godoc
// @Summary Delete a track
// @Description Delete a track by its ID
// @Tags tracks
// @Produce json
// @Param id path string true "Track ID"
// @Success 204 "No Content"
// @Failure 400 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/tracks/{id} [delete]
func (c *TrackController) DeleteTrack(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)
	if claims == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	err = c.trackService.DeleteTrackMeta(ctx.Request.Context(), claims, id)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Failed to delete track"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
