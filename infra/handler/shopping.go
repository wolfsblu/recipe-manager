package handler

import (
	"context"

	"github.com/wolfsblu/recipe-manager/api"
	"github.com/wolfsblu/recipe-manager/domain"
	"github.com/wolfsblu/recipe-manager/infra/config"
	"github.com/wolfsblu/recipe-manager/infra/env"
	"github.com/wolfsblu/recipe-manager/infra/handler/mapper"
)

type ShoppingHandler struct {
	mapper   *mapper.APIMapper
	Shopping *domain.ShoppingService
}

func NewShoppingHandler(service *domain.ShoppingService) *ShoppingHandler {
	return &ShoppingHandler{
		mapper:   mapper.NewAPIMapper(env.MustGet("BASE_URL")),
		Shopping: service,
	}
}

func (h *ShoppingHandler) GetShoppingLists(ctx context.Context, params api.GetShoppingListsParams) (*api.PaginatedShoppingLists, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}

	paginationReq, err := domain.ValidatePage(params.Cursor.Value, int(params.Limit.Or(domain.DefaultPageSize)))
	if err != nil {
		return nil, err
	}

	result, err := h.Shopping.GetByUser(ctx, user, paginationReq)
	if err != nil {
		return nil, err
	}

	lists, err := h.mapper.ToShoppingLists(result.Data)
	if err != nil {
		return nil, err
	}

	var nextCursor api.OptNilString
	if result.NextCursor != nil {
		nextCursor = api.NewOptNilString(*result.NextCursor)
	}

	return &api.PaginatedShoppingLists{
		Data:       lists,
		NextCursor: nextCursor,
		HasMore:    result.HasMore,
	}, nil
}

func (h *ShoppingHandler) CreateShoppingList(ctx context.Context, req *api.WriteShoppingList) (*api.ReadShoppingList, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	list, err := h.Shopping.Create(ctx, user, req.Name)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToShoppingList(list)
}

func (h *ShoppingHandler) GetShoppingListById(ctx context.Context, params api.GetShoppingListByIdParams) (*api.ReadShoppingList, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	list, err := h.Shopping.GetByID(ctx, user, params.ShoppingListId)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToShoppingList(list)
}

func (h *ShoppingHandler) UpdateShoppingList(ctx context.Context, req *api.WriteShoppingList, params api.UpdateShoppingListParams) (*api.ReadShoppingList, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	list, err := h.Shopping.Update(ctx, user, params.ShoppingListId, req.Name)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToShoppingList(list)
}

func (h *ShoppingHandler) DeleteShoppingList(ctx context.Context, params api.DeleteShoppingListParams) error {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return domain.ErrAuthentication
	}
	return h.Shopping.Delete(ctx, user, params.ShoppingListId)
}

func (h *ShoppingHandler) AddShoppingListItem(ctx context.Context, req *api.WriteShoppingListItem, params api.AddShoppingListItemParams) (*api.ReadShoppingListItem, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	item := h.mapper.FromWriteShoppingListItem(req)
	result, err := h.Shopping.AddItem(ctx, user, params.ShoppingListId, item)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToShoppingListItem(result)
}

func (h *ShoppingHandler) UpdateShoppingListItem(ctx context.Context, req *api.WriteShoppingListItem, params api.UpdateShoppingListItemParams) (*api.ReadShoppingListItem, error) {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return nil, domain.ErrAuthentication
	}
	item := h.mapper.FromWriteShoppingListItem(req)
	result, err := h.Shopping.UpdateItem(ctx, user, params.ShoppingListId, params.ItemId, item)
	if err != nil {
		return nil, err
	}
	return h.mapper.ToShoppingListItem(result)
}

func (h *ShoppingHandler) DeleteShoppingListItem(ctx context.Context, params api.DeleteShoppingListItemParams) error {
	user, ok := ctx.Value(config.CtxKeyUser).(*domain.User)
	if !ok || user == nil {
		return domain.ErrAuthentication
	}
	return h.Shopping.DeleteItem(ctx, user, params.ShoppingListId, params.ItemId)
}
