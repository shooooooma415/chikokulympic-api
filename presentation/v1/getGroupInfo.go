package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetGroupInfo struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

type GroupInfoResponse struct {
	GroupName      entity.GroupName     `json:"group_name" example:"テストグループ"`
	Password       entity.GroupPassword `json:"password" example:"password123"`
	GroupMembers   []usecase.Member     `json:"group_members"`
	GroupManagerID entity.UserID        `json:"manager_id" example:"user456"`
}

func NewGetGroupInfo(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *GetGroupInfo {
	return &GetGroupInfo{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// @Summary get group info
// @Description get chosen group info
// @Tags groups
// @Accept json
// @Produce json
// @Param group_id path string true "group_id"
// @Success 200 {object} GroupInfoResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /groups/{group_id} [get]
func (g *GetGroupInfo) Handler(c echo.Context) error {
	groupIDParam := c.Param("group_id")
	if groupIDParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループIDは必須です"))
	}

	groupID := entity.GroupID(groupIDParam)
	fetchGroupInfoUseCase := usecase.NewFetchGroupInfoUsecase(g.groupRepo, g.userRepo)
	result, err := fetchGroupInfoUseCase.Execute(groupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	if result == nil {
		return c.JSON(http.StatusNotFound, middleware.NewErrorResponse("グループが見つかりません"))
	}

	response := &GroupInfoResponse{
		GroupName:      result.GroupName,
		Password:       result.Password,
		GroupMembers:   result.Members,
		GroupManagerID: result.GroupManagerID,
	}

	return c.JSON(http.StatusOK, response)
}
