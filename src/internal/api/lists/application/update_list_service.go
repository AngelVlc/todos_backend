package application

import (
	"github.com/AngelVlc/todos/internal/api/lists/domain"
	appErrors "github.com/AngelVlc/todos/internal/api/shared/domain/errors"
)

type UpdateListService struct {
	repo domain.ListsRepository
}

func NewUpdateListService(repo domain.ListsRepository) *UpdateListService {
	return &UpdateListService{repo}
}

func (s *UpdateListService) UpdateList(listID int32, name domain.ListName, userID int32, IDsByPosition []int32) (*domain.List, error) {
	foundList, err := s.repo.FindListByID(listID, userID)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error getting the user list", InternalError: err}
	}

	if foundList == nil {
		return nil, &appErrors.BadRequestError{Msg: "The list does not exist"}
	}

	foundList.Name = name

	err = s.repo.UpdateList(foundList)

	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error updating the user list", InternalError: err}
	}

	foundItems, err := s.repo.GetAllListItems(listID, userID)
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

	err = s.repo.BulkUpdateListItems(foundItems)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error bulk updating", InternalError: err}
	}

	return foundList, nil
}
