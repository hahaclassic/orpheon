package search_ctrl

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hahaclassic/orpheon/backend/internal/controller/http/utils"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/search"
)

type SearchController struct {
	searchService  search.SearchService
	authMiddleware gin.HandlerFunc
}

func NewSearchController(searchService search.SearchService, authMiddleware gin.HandlerFunc) *SearchController {
	return &SearchController{
		searchService:  searchService,
		authMiddleware: authMiddleware,
	}
}

func (c *SearchController) RegisterRoutes(router *gin.RouterGroup) {
	search := router.Group("/search")
	search.Use(c.authMiddleware) // optional auth middleware
	{
		search.GET("", c.Search)
	}
}

func (c *SearchController) parseSearchRequest(ctx *gin.Context) (*entity.SearchRequest, error) {
	query := ctx.Query("query")
	if query == "" {
		return nil, fmt.Errorf("Query parameter is required")
	}

	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid limit parameter")
	}
	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid offset parameter")
	}

	genre := ctx.Query("genre")
	country := ctx.Query("country")

	searchRequest := &entity.SearchRequest{
		Query:  query,
		Limit:  limit,
		Offset: offset,
		Filters: entity.Filters{
			Genre:   genre,
			Country: country,
		},
	}
	return searchRequest, nil
}

func (c *SearchController) Search(ctx *gin.Context) {
	claims := utils.GetClaims(ctx)

	searchRequest, err := c.parseSearchRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentType := ctx.Query("type")

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

	case "playlist":
		result, err := c.searchService.SearchPlaylists(ctx.Request.Context(), claims, searchRequest)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform search"})
			return
		}
		ctx.JSON(http.StatusOK, result)

	default:
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}
}
