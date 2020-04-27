package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
	"github.com/AngelVlc/todos/models"
)

func GetUserListsHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	res := []dtos.GetListsResultDto{}
	err := h.listsSrv.GetUserLists(userID, &res)
	if err != nil {
		return errorResult{err}
	}
	return okResult{res, http.StatusOK}
}

func GetUserSingleListHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	l := dtos.GetSingleListResultDto{}
	err = h.listsSrv.GetSingleUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func AddUserListHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	id, err := h.listsSrv.AddUserList(userID, l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func UpdateUserListHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	err = h.listsSrv.UpdateUserList(listID, userID, l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func DeleteUserListHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	listID, err := parseInt32UrlVar(r, "id")
	if err != nil {
		return errorResult{err}
	}

	err = h.listsSrv.RemoveUserList(listID, userID)
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
