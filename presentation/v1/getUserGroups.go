package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetUserGroups struct {
	groupRepo repository.GroupRepository
}

func NewGetUserGroups(groupRepo repository.GroupRepository) *GetUserGroups {
	return &GetUserGroups{
		groupRepo: groupRepo,
	}
}

// @Summary get user groups
// @Description get user groups
// @Tags groups
// @Accept json
// @Produce json
// @Param user_id path string true "user_id"
// @Success 200 {array} usecase.UserGroup
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /users/{user_id}/groups [get]
func (g *GetUserGroups) Handler(c echo.Context) error {
	userIDParam := c.Param("user_id")

	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザーIDは必須です"))
	}

	userID := entity.UserID(userIDParam)

	result, err := usecase.NewFetchUserGroupsUseCase(g.groupRepo, userID).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}
