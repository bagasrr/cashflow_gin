package handler

import (
	"cashflow_gin/api"
	"cashflow_gin/services"
	"context"
)

type GroupAPI struct {
	Service services.GroupService
}

func (c *GroupAPI) GetGroups(ctx context.Context, request api.GetGroupsRequestObject) (api.GetGroupsResponseObject, error) {
	return api.GetGroups201JSONResponse{}, nil
}

func (c *GroupAPI) CreateGroup(ctx context.Context, request api.CreateGroupRequestObject) (api.CreateGroupResponseObject, error) {
	return api.CreateGroup201JSONResponse{}, nil
}

func (c *GroupAPI) DeleteGroup(ctx context.Context, request api.DeleteGroupRequestObject) (api.DeleteGroupResponseObject, error) {
	return api.DeleteGroup200JSONResponse{}, nil
}

func (c *GroupAPI) GetGroupById(ctx context.Context, request api.GetGroupByIdRequestObject) (api.GetGroupByIdResponseObject, error) {
	return api.GetGroupById200JSONResponse{}, nil
}

func (c *GroupAPI) UpdateGroup(ctx context.Context, request api.UpdateGroupRequestObject) (api.UpdateGroupResponseObject, error) {
	return api.UpdateGroup201JSONResponse{}, nil
}
