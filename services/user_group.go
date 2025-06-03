package services

import (
	gw "dig-inv/gen/go"
)

type userGroupServer struct {
	gw.UnimplementedUserGroupServiceServer
}

func NewUserGroupServer() gw.UserGroupServiceServer {
	return &userGroupServer{}
}
