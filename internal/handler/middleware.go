package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/pkg/authProto"
	"strings"
)

const (
	authorizationHeader = "Authorization"

	managerCtx    = "managerID"
	restaurantCtx = "restaurantID"
)

func corsMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func (h *Handler) ManagerIdentity(c *gin.Context) {
	user, err := h.parseAuthHeader(c)
	if err != nil {
		switch err {
		case domain.ErrInternalServer:
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		default:
			newResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
	}

	if user.Role != "Restaurant manager" {
		newResponse(c, http.StatusForbidden, "the given user has a different role")
	}

	restaurantID, err := h.services.GetRestaurantID(int(user.UserId))
	if err != nil {
		switch err {
		case domain.ErrManagerNotAssigned:
			newResponse(c, http.StatusUnauthorized, err.Error())
			return
		default:
			newResponse(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	c.Set(restaurantCtx, restaurantID)
	c.Set(managerCtx, user)
}

func (h *Handler) SuperAdminIdentity(c *gin.Context) {
	user, err := h.parseAuthHeader(c)
	if err != nil {
		switch err {
		case domain.ErrInternalServer:
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		default:
			newResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
	}

	if user.Role != "Superadmin" {
		newResponse(c, http.StatusForbidden, "the given user has a different role")
	}

	c.Set(managerCtx, user)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (*authProto.UserRole, error) {
	header := c.GetHeader(authorizationHeader)

	if header == "" {
		return nil, errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return nil, errors.New("token is empty")
	}

	user, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
			log.Error().Err(err).Msg("error occurred while verification JWT in auth service")
			return nil, errors.New("internal server error")
	}

	return user, nil
}
