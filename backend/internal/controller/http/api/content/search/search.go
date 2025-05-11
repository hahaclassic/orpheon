package search

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
)

type SearchController struct {
	searchService search.SearchService
}

func New(searchService search.SearchService) *SearchController {
	return &SearchController{
		searchService: searchService,
	}
}

func (c *SearchController) RegisterRoutes(router *gin.RouterGroup) {
	search := router.Group("/search")
	{
		search.GET("", c.Search)
	}
}

// Search godoc
// @Summary Search for content
// @Description Search for tracks, albums, artists, and playlists
// @Tags search
// @Produce json
// @Param query query string true "Search query"
// @Param type query string false "Content type (track, album, artist, playlist)"
// @Param limit query int false "Maximum number of results"
// @Param offset query int false "Number of results to skip"
// @Success 200 {object} search.SearchResult
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search [get]
func (c *SearchController) Search(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
		return
	}

	contentType := ctx.Query("type")
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset parameter"})
		return
	}

	searchRequest := &entity.SearchRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}

	switch contentType {
	case "track":
		result, err := c.searchService.SearchTracks(ctx.Request.Context(), searchRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
			return
		}
		ctx.JSON(http.StatusOK, result)

	case "album":
		result, err := c.searchService.SearchAlbums(ctx.Request.Context(), searchRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
			return
		}
		ctx.JSON(http.StatusOK, result)

	case "artist":
		result, err := c.searchService.SearchArtists(ctx.Request.Context(), searchRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
			return
		}
		ctx.JSON(http.StatusOK, result)

	// case "playlist":
	// 	result, err := c.searchService.SearchPlaylists(ctx.Request.Context(), searchRequest)
	// 	if err != nil {
	// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusOK, result)

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
}
