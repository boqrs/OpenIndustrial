package asset

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// API provides the HTTP handlers for the asset domain.
type API struct {
	service *Service
}

// NewAPI creates a new asset API handler.
func NewAPI(service *Service) *API {
	return &API{
		service: service,
	}
}

// RegisterRoutes registers the asset API routes with a Gin router.
func (a *API) RegisterRoutes(router *gin.RouterGroup) {
	assetRoutes := router.Group("/assets")
	{
		assetRoutes.POST("", a.createAsset)
		assetRoutes.GET("", a.listAssets)
		// It's common to look up assets by SN
		assetRoutes.GET("/sn/:sn", a.getAssetBySN)
	}
}

func (a *API) createAsset(c *gin.Context) {
	orgIDStr := c.Param("org_id") // Assuming org_id is a URL parameter from a parent router
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	asset, err := a.service.CreateAsset(c.Request.Context(), orgID, productID, req.SN)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create asset"})
		return
	}

	c.JSON(http.StatusCreated, ToAssetResponse(asset))
}

func (a *API) listAssets(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	assets, err := a.service.ListAssetsForOrg(c.Request.Context(), orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list assets"})
		return
	}

	c.JSON(http.StatusOK, ToAssetListResponse(assets))
}

func (a *API) getAssetBySN(c *gin.Context) {
	orgIDStr := c.Param("org_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid organization id"})
		return
	}

	sn := c.Param("sn")

	asset, err := a.service.GetAssetBySN(c.Request.Context(), orgID, sn)
	if err != nil {
		if err == ErrAssetNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get asset"})
		return
	}

	c.JSON(http.StatusOK, ToAssetResponse(asset))
}