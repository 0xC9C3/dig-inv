package services

import (
	gw "dig-inv/gen/go"
)

type itemServer struct {
	gw.UnimplementedItemServiceServer
}

func NewItemServer() gw.ItemServiceServer {
	return &itemServer{}
}
