package application

import (
	"context"

	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UpdateListService struct {
	repo domain.ListsRepository
}

func NewUpdateListService(repo domain.ListsRepository) *UpdateListService {
	return &UpdateListService{repo}
}

func (s *UpdateListService) UpdateList(ctx context.Context, listID int32, name domain.ListName, userID int32, IDsByPosition []int32) (*domain.List, error) {
	foundList, err := s.repo.FindListByID(ctx, listID, userID)
	if err != nil {
		return nil, err
	}

	if foundList.Name != name {
		err = name.CheckIfAlreadyExists(ctx, userID, s.repo)
		if err != nil {
			return nil, err
		}
	}

	foundList.Name = name

	err = s.repo.UpdateList(ctx, foundList)

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user list", InternalError: err}
	}

	foundItems, err := s.repo.GetAllListItems(ctx, listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting all list items", InternalError: err}
	}

	if len(foundItems) == 0 {
		return foundList, nil
	}

	for i := 0; i < len(IDsByPosition); i++ {
		for j := 0; j < len(foundItems); j++ {
			if foundItems[j].ID == int32(IDsByPosition[i]) {
				foundItems[j].Position = int32(i)
				break
			}
		}
	}

	err = s.repo.BulkUpdateListItems(ctx, foundItems)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error bulk updating", InternalError: err}
	}

	return foundList, nil
}
