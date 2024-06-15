package application

import (
	"context"
	"fmt"

	"github.com/AngelVlc/todos_backend/src/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/errors"
	"github.com/AngelVlc/todos_backend/src/internal/api/shared/domain/events"
)

type MoveListItemService struct {
	repo     domain.ListsRepository
	eventBus events.EventBus
}

func NewMoveListItemService(repo domain.ListsRepository, eventBus events.EventBus) *MoveListItemService {
	return &MoveListItemService{repo, eventBus}
}

func (s *MoveListItemService) MoveListItem(ctx context.Context, originListID int32, originListItemID int32, destinationListID int32, userID int32) error {
	foundOriginList, err := s.repo.FindList(ctx, domain.ListRecord{ID: originListID, UserID: userID})
	if err != nil {
		return err
	}

	foundDestinationList, err := s.repo.FindList(ctx, domain.ListRecord{ID: destinationListID, UserID: userID})
	if err != nil {
		return &appErrors.BadRequestError{Msg: "The destination list does not exist"}
	}

	indexToRemove := -1

	for i, item := range foundOriginList.Items {
		if item.ID == originListItemID {
			indexToRemove = i
			item.ListID = destinationListID
			item.Position = foundDestinationList.GetMaxItemPosition() + 1
			foundDestinationList.Items = append(foundDestinationList.Items, item)

			break
		}
	}

	if indexToRemove < 0 {
		return &appErrors.BadRequestError{Msg: fmt.Sprintf("An item with id %v doesn't exist in the original list", originListItemID)}
	}

	foundOriginList.Items = append(foundOriginList.Items[:indexToRemove], foundOriginList.Items[indexToRemove+1:]...)

	if err = s.repo.UpdateList(ctx, &foundOriginList); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating the original list", InternalError: err}
	}

	if err = s.repo.UpdateList(ctx, &foundDestinationList); err != nil {
		return &appErrors.UnexpectedError{Msg: "Error updating the destination list", InternalError: err}
	}

	go s.eventBus.Publish(events.ListUpdated, foundOriginList.ID)
	go s.eventBus.Publish(events.ListUpdated, foundDestinationList.ID)

	return nil
}
