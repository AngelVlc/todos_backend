package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/dtos"
)

func GetUserListsHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	res, err := h.listsSrv.GetUserLists(userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{res, http.StatusOK}
}

func GetUserSingleListHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "id")

	res, err := h.listsSrv.GetUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{res, http.StatusOK}
}

func AddUserListHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)

	dto, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	id, err := h.listsSrv.AddUserList(userID, dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func UpdateUserListHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "id")

	dto, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}

	err = h.listsSrv.UpdateUserList(listID, userID, dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{dto, http.StatusOK}
}

func DeleteUserListHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "id")

	err := h.listsSrv.RemoveUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func GetUserSingleListItemHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	dto, err := h.listItemsSrv.GetListItem(itemID, listID, userID)
	if err != nil {
		return errorResult{err}
	}

	return okResult{dto, http.StatusOK}
}

func AddUserListItemHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")

	i, err := parseListItemBody(r)
	if err != nil {
		return errorResult{err}
	}

	id, err := h.listItemsSrv.AddListItem(listID, userID, i)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func DeleteUserListItemHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	err := h.listItemsSrv.RemoveItem(itemID, listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func UpdateUserListItemHandler(w http.ResponseWriter, r *http.Request, h Handler) HandlerResult {
	userID := getUserIDFromContext(r)
	listID := parseInt32UrlVar(r, "listId")
	itemID := parseInt32UrlVar(r, "itemId")

	i, err := parseListItemBody(r)
	if err != nil {
		return errorResult{err}
	}

	err = h.listsSrv.UpdateUserListItem(itemID, listID, userID, i)
	if err != nil {
		return errorResult{err}
	}
	return okResult{i, http.StatusOK}
}

func parseListBody(r *http.Request) (*dtos.ListDto, error) {
	var dto dtos.ListDto
	err := parseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func parseListItemBody(r *http.Request) (*dtos.ListItemDto, error) {
	var dto dtos.ListItemDto
	err := parseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}
