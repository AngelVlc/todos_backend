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
	listID := parseInt32UrlVar(r, "id")

	l := dtos.GetSingleListResultDto{}
	err := h.listsSrv.GetSingleUserList(listID, userID, &l)
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
	listID := parseInt32UrlVar(r, "id")

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
	listID := parseInt32UrlVar(r, "id")

	err := h.listsSrv.RemoveUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func GetUserSingleListItemHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	l := dtos.GetItemResultDto{}
	err := h.listsSrv.GetUserListItem(itemID, listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func AddUserListItemHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")

	i, err := parseListItemBody(r)
	if err != nil {
		return errorResult{err}
	}

	i.ListID = listID

	id, err := h.listsSrv.AddUserListItem(userID, i)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func DeleteUserListItemHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	err := h.listsSrv.RemoveUserListItem(itemID, listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func UpdateUserListItemHandler(r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	i, err := parseListItemBody(r)
	if err != nil {
		return errorResult{err}
	}

	i.ListID = listID

	err = h.listsSrv.UpdateUserListItem(itemID, listID, userID, i)
	if err != nil {
		return errorResult{err}
	}
	return okResult{i, http.StatusOK}
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

func parseListItemBody(r *http.Request) (*models.ListItem, error) {
	var dto dtos.ListItemDto
	err := parseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	l := dto.ToListItem()
	return &l, nil
}
