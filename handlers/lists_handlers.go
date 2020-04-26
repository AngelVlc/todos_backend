package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	"github.com/AngelVlc/todos/models"
	"github.com/AngelVlc/todos/wire"
)

func GetUserLists(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listSrv := wire.InitListsService(h.Db)
	res := []dtos.GetListsResultDto{}
	err := listSrv.GetUserLists(userID, &res)
	if err != nil {
		return errorResult{err}
	}
	return okResult{res, http.StatusOK}
}

func GetUserSingleList(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	listSrv := wire.InitListsService(h.Db)

	l := dtos.GetSingleListResultDto{}
	err = listSrv.GetSingleUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func AddUserList(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	listSrv := wire.InitListsService(h.Db)

	id, err := listSrv.AddUserList(userID, l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func UpdateUserList(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	listSrv := wire.InitListsService(h.Db)
	err = listSrv.UpdateUserList(listID, userID, l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func DeleteUserList(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	listSrv := wire.InitListsService(h.Db)
	err = listSrv.RemoveUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func parseListBody(r *http.Request) (*models.List, error) {
	var dto dtos.ListDto
	err := parseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	l := dto.ToList()
	return &l, nil
}
