package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type UserHandler struct {
	ctx context.Context
	svc *UserService
}

func NewUserHandler(ctx context.Context, conn *pgx.Conn) *UserHandler {
	return &UserHandler{
		ctx: ctx,
		svc: NewUserService(conn),
	}
}

func (h *UserHandler) BindRoutes(router *gin.Engine) {
	router.POST("/users/register", h.register)
}

func (h *UserHandler) register(c *gin.Context) {
	var pa RegisterParams
	if err := c.ShouldBindJSON(&pa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.Register(h.ctx, &pa)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}
