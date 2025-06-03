package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/usecase"
	"fmt"
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"chikokulympic-api/middleware"
)
type GetEventsResponse struct {
	Events []entity.Event `json:"events" example:"[{\"event_id\":\"event123\",\"event_name\":\"Event Name\",\"event_description\":\"Event Description\",\"event_date\":\"2023-10-01T00:00:00Z\"}]"`
}

type GetEvents struct {
	eventRepo repository.EventRepository
	groupRepo repository.GroupRepository
}

func NewGetEvents(eventRepo repository.EventRepository, groupRepo repository.GroupRepository) *GetEvents {
	return &GetEvents{
		eventRepo: eventRepo,
		groupRepo: groupRepo,
	}
}

func (g *GetEvents) Handler(c echo.Context) error {
	groupIDsParam := c.QueryParam("group_ids")
	if groupIDsParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループIDは必須です"))
	}

	groupIDStrings := strings.Split(groupIDsParam, ",")
	var groupIDs []entity.GroupID

	for _, idStr := range groupIDStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		// グループIDをUUIDに変換
		groupUUID, err := uuid.Parse(idStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(fmt.Sprintf("無効なグループID: %s", idStr)))
		}

		groupIDs = append(groupIDs, entity.GroupID(groupUUID.String()))
	}

	if len(groupIDs) == 0 {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("有効なグループIDが指定されていません"))
	}

	events, err := usecase.NewFetchEventInfoUsecase(g.groupRepo, g.eventRepo, groupIDs).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(fmt.Sprintf("イベント情報の取得中にエラーが発生しました: %v", err)))
	}

	response := GetEventsResponse{
		Events: events,
	}

	return c.JSON(http.StatusOK, response)
}
