package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/AngelVlc/todos/dtos"
	appErrors "github.com/AngelVlc/todos/errors"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/wire"
	"github.com/jinzhu/gorm"
)

// ListsHandler is the handler for the lists endpoints
func ListsHandler(r *http.Request, db *gorm.DB) handlerResult {
	switch r.Method {
	case http.MethodGet:
		return processListsGET(r, db)
	case http.MethodPost:
		return processListsPOST(r, db)
	case http.MethodDelete:
		return processListsDELETE(r, db)
	case http.MethodPut:
		return processListsPUT(r, db)
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}
}

func processListsGET(r *http.Request, db *gorm.DB) handlerResult {
	listID := getListIDFromURL(r.URL)
	userID := getUserIDFromContext(r)

	listSrv := wire.InitListsService(db)
	if listID == 0 {
		r := []dtos.GetListsResultDto{}
		err := listSrv.GetUserLists(userID, &r)
		if err != nil {
			return errorResult{err}
		}
		return okResult{r, http.StatusOK}
	}

	l := dtos.GetSingleListResultDto{}
	err := listSrv.GetSingleUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func processListsPOST(r *http.Request, db *gorm.DB) handlerResult {
	l, err := parseListBody(r)
	userID := getUserIDFromContext(r)

	if err != nil {
		return errorResult{err}
	}

	listSrv := wire.InitListsService(db)

	id, err := listSrv.AddUserList(userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func processListsPUT(r *http.Request, db *gorm.DB) handlerResult {
	listID := getListIDFromURL(r.URL)
	userID := getUserIDFromContext(r)

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}
	listSrv := wire.InitListsService(db)
	err = listSrv.UpdateUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func processListsDELETE(r *http.Request, db *gorm.DB) handlerResult {
	listID := getListIDFromURL(r.URL)
	userID := getUserIDFromContext(r)

	listSrv := wire.InitListsService(db)
	err := listSrv.RemoveUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func getListIDFromURL(u *url.URL) int32 {
	var r int32

	if len(u.Path) > len("/lists") {
		listID := u.Path[len("/lists/"):]
		i, _ := strconv.ParseInt(listID, 10, 32)
		r = int32(i)
	}

	return r
}

func parseListBody(r *http.Request) (models.List, error) {
	decoder := json.NewDecoder(r.Body)
	var dto dtos.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.List{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	l := dto.ToList()

	return l, nil
}
