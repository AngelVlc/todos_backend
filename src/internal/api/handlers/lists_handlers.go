package handlers

import (
	"net/http"

	"github.com/AngelVlc/todos/internal/api/dtos"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/handler"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/helpers"
	"github.com/AngelVlc/todos/internal/api/shared/infrastructure/results"
)

func GetUserListsHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	res, err := h.ListsSrv.GetUserLists(userID)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{res, http.StatusOK}
}

func GetUserSingleListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "id")

	res, err := h.ListsSrv.GetUserList(listID, userID)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{res, http.StatusOK}
}

func AddUserListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)

	dto, err := parseListBody(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	id, err := h.ListsSrv.AddUserList(userID, dto)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{id, http.StatusCreated}
}

func UpdateUserListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "id")

	dto, err := parseListBody(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	err = h.ListsSrv.UpdateUserList(listID, userID, dto)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{dto, http.StatusOK}
}

func DeleteUserListHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "id")

	err := h.ListsSrv.RemoveUserList(listID, userID)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{nil, http.StatusNoContent}
}

func GetUserSingleListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "listId")
	itemID := helpers.ParseInt32UrlVar(r, "itemId")

	dto, err := h.ListItemsSrv.GetListItem(itemID, listID, userID)
	if err != nil {
		return results.ErrorResult{err}
	}

	return results.OkResult{dto, http.StatusOK}
}

func AddUserListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "listId")

	i, err := parseListItemBody(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	id, err := h.ListItemsSrv.AddListItem(listID, userID, i)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{id, http.StatusCreated}
}

func DeleteUserListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "listId")
	itemID := helpers.ParseInt32UrlVar(r, "itemId")

	err := h.ListItemsSrv.RemoveListItem(itemID, listID, userID)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{nil, http.StatusNoContent}
}

func UpdateUserListItemHandler(w http.ResponseWriter, r *http.Request, h handler.Handler) handler.HandlerResult {
	userID := helpers.GetUserIDFromContext(r)
	listID := helpers.ParseInt32UrlVar(r, "listId")
	itemID := helpers.ParseInt32UrlVar(r, "itemId")

	i, err := parseListItemBody(r)
	if err != nil {
		return results.ErrorResult{err}
	}

	err = h.ListItemsSrv.UpdateListItem(itemID, listID, userID, i)
	if err != nil {
		return results.ErrorResult{err}
	}
	return results.OkResult{i, http.StatusOK}
}

func parseListBody(r *http.Request) (*dtos.ListDto, error) {
	var dto dtos.ListDto
	err := helpers.ParseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func parseListItemBody(r *http.Request) (*dtos.ListItemDto, error) {
	var dto dtos.ListItemDto
	err := helpers.ParseBody(r, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}
