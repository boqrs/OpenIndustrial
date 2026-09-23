package iot

import (
	"fmt"
	"net/http"
	"strconv"

	srv "github.com/boqrs/OpenIndustrial/cloud/internal/services/iot"
	"github.com/boqrs/zeus/ginx"
	"github.com/gin-gonic/gin"
)

// Handler handles IoT HTTP endpoints.
type Handler struct {
	service srv.Service
}

// NewHandler creates a new IoT handler.
func NewHandler(
	service srv.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

// RouterRegister registers IoT routes.
//
// These endpoints are internal integration endpoints for the MQTT
// infrastructure. They are not end-user APIs.
func (h *Handler) RouterRegister(
	router ginx.ZeroGinRouter,
) {
	internalGroup := router.Group(
		"/api/v1/internal/iot",
	)

	internalGroup.Handle(
		http.MethodPost,
		"/mqtt/authenticate",
		h.authenticateMQTT,
	)

	internalGroup.Handle(
		http.MethodPost,
		"/mqtt/authorize",
		h.authorizeMQTT,
	)

	internalGroup.Handle(
		http.MethodGet,
		"/devices/:resource_id/topics",
		h.getDeviceTopics,
	)

	internalGroup.Handle(
		http.MethodPost,
		"/devices/:resource_id/online",
		h.deviceOnline,
	)

	internalGroup.Handle(
		http.MethodPost,
		"/devices/:resource_id/offline",
		h.deviceOffline,
	)

	internalGroup.Handle(
		http.MethodPost,
		"/devices/:resource_id/heartbeat",
		h.deviceHeartbeat,
	)
}

func (h *Handler) authenticateMQTT(
	ctx *gin.Context,
) ginx.Render {
	var req srv.AuthenticateMQTTRequest

	if err := ctx.ShouldBindJSON(
		&req,
	); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	resp, err := h.service.AuthenticateMQTT(
		ctx.Request.Context(),
		req,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) authorizeMQTT(
	ctx *gin.Context,
) ginx.Render {
	var req srv.AuthorizeMQTTRequest

	if err := ctx.ShouldBindJSON(
		&req,
	); err != nil {
		return ginx.Error(
			fmt.Errorf("invalid param"),
		)
	}

	if err := h.service.AuthorizeMQTT(
		ctx.Request.Context(),
		req,
	); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(nil)
}

func (h *Handler) getDeviceTopics(
	ctx *gin.Context,
) ginx.Render {
	resourceID, err := parseResourceID(
		ctx.Param("resource_id"),
	)
	if err != nil {
		return ginx.Error(err)
	}

	resp, err := h.service.GetDeviceTopics(
		ctx.Request.Context(),
		resourceID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) deviceOnline(
	ctx *gin.Context,
) ginx.Render {
	resourceID, err := parseResourceID(
		ctx.Param("resource_id"),
	)
	if err != nil {
		return ginx.Error(err)
	}

	resp, err := h.service.DeviceOnline(
		ctx.Request.Context(),
		resourceID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) deviceOffline(
	ctx *gin.Context,
) ginx.Render {
	resourceID, err := parseResourceID(
		ctx.Param("resource_id"),
	)
	if err != nil {
		return ginx.Error(err)
	}

	resp, err := h.service.DeviceOffline(
		ctx.Request.Context(),
		resourceID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func (h *Handler) deviceHeartbeat(
	ctx *gin.Context,
) ginx.Render {
	resourceID, err := parseResourceID(
		ctx.Param("resource_id"),
	)
	if err != nil {
		return ginx.Error(err)
	}

	resp, err := h.service.DeviceHeartbeat(
		ctx.Request.Context(),
		resourceID,
	)
	if err != nil {
		return ginx.Error(err)
	}

	return ginx.Success(resp)
}

func parseResourceID(
	value string,
) (uint, error) {
	id, err := strconv.ParseUint(
		value,
		10,
		32,
	)
	if err != nil || id == 0 {
		return 0, fmt.Errorf(
			"invalid resource id",
		)
	}

	return uint(id), nil
}
