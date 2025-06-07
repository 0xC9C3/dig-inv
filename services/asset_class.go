package services

import (
	"context"
	gw "dig-inv/gen/go"
	"dig-inv/log"
	"dig-inv/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

type assetClassServer struct {
	gw.AssetClassServiceServer
}

func (a assetClassServer) GetAssetClasses(ctx context.Context, message *gw.EmptyMessage) (*gw.AssetClasses, error) {
	client, err := store.GetClient()
	if err != nil {
		grpclog.Errorf("Failed to get store client: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get store client: %v", err)
	}

	classes, err := client.AssetClass.Query().All(ctx)
	if err != nil {
		grpclog.Errorf("Failed to query asset classes: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to query asset classes: %v", err)
	}

	log.S.Debugw("Retrieved asset classes", "count", len(classes))

	res := make([]*gw.AssetClass, 0, len(classes))
	for _, class := range classes {
		res = append(res, &gw.AssetClass{
			Id:          class.ID.String(),
			Name:        class.Name,
			Description: class.Description,
			Icon:        class.Icon,
			Color:       class.Color,
		})
	}

	return &gw.AssetClasses{
		Classes: res,
	}, nil
}

func (a assetClassServer) CreateAssetClass(ctx context.Context, class *gw.AssetClass) (*gw.AssetClass, error) {
	client, err := store.GetClient()
	if err != nil {
		grpclog.Errorf("Failed to get store client: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get store client: %v", err)
	}

	user := ctx.Value(AuthenticatedSubjectKey).(string)

	newClass, err := client.AssetClass.Create().
		SetName(class.Name).
		SetDescription(class.Description).
		SetIcon(class.Icon).
		SetColor(class.Color).
		SetCreatedBy(user).
		SetUpdatedBy(user).
		Save(ctx)
	if err != nil {
		grpclog.Errorf("Failed to create asset class: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create asset class: %v", err)
	}
	log.S.Debugw("Created new asset class", "id", newClass.ID, "name", newClass.Name)
	return &gw.AssetClass{
		Id:          newClass.ID.String(),
		Name:        newClass.Name,
		Description: newClass.Description,
	}, nil

}

func NewAssetClassServer() gw.AssetClassServiceServer {
	return &assetClassServer{}
}
