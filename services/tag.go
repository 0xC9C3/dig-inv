package services

import (
	gw "dig-inv/gen/go"
)

type tagServer struct {
	gw.UnimplementedTagServiceServer
}

func NewTagServer() gw.TagServiceServer {
	return &tagServer{}
}
