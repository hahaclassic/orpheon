package track_ctrl

import (
	"github.com/gin-gonic/gin"
)

type TrackRouter struct {
	trackMetaController *TrackMetaController
	segmentService      *TrackSegmentController
	audioService        *TrackAudioController
	authMiddleware      gin.HandlerFunc
}

func NewTrackRouter(trackMetaController *TrackMetaController,
	segmentService *TrackSegmentController,
	audioService *TrackAudioController,
	authMiddleware gin.HandlerFunc) *TrackRouter {
	return &TrackRouter{
		trackMetaController: trackMetaController,
		segmentService:      segmentService,
		audioService:        audioService,
		authMiddleware:      authMiddleware,
	}
}

func (r *TrackRouter) RegisterRoutes(router *gin.RouterGroup) {
	tracks := router.Group("/tracks")
	{
		tracks.GET("/:id", r.trackMetaController.GetTrack)
		tracks.GET("/:id/segments", r.segmentService.GetSegments)

		tracksProtected := tracks.Group("")
		tracksProtected.Use(r.authMiddleware)
		{
			tracksProtected.POST("", r.trackMetaController.CreateTrack)
			tracksProtected.PUT("/:id", r.trackMetaController.UpdateTrack)
			tracksProtected.DELETE("/:id", r.trackMetaController.DeleteTrack)
		}

		tracksAudio := tracks.Group("/:id/audio")
		{
			tracksAudio.GET("", r.audioService.GetAudioChunk)
			tracksAudioProtected := tracksAudio.Group("")
			tracksAudioProtected.Use(r.authMiddleware)
			{
				tracksAudioProtected.POST("", r.audioService.UploadAudioFile)
				tracksAudioProtected.DELETE("", r.audioService.DeleteAudioFile)
			}
		}
	}
}
