package services

import (
	"context"
	"dig-inv/ent/assetclass"
	gw "dig-inv/gen/go"
	"dig-inv/log"
	"dig-inv/store"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"time"
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

	classes, err := client.AssetClass.Query().Where(
		assetclass.DeletedAtIsNil()).Order(
		assetclass.ByOrder(sql.OrderDesc()),
	).All(ctx)
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
			Provider:    class.Provider,
			Description: class.Description,
			Order:       int32(class.Order),
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
		SetOrder(int(class.Order)).
		SetProvider(class.Provider).
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

func (a assetClassServer) UpdateAssetClass(ctx context.Context, class *gw.AssetClass) (*gw.AssetClass, error) {
	client, err := store.GetClient()
	if err != nil {
		grpclog.Errorf("Failed to get store client: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get store client: %v", err)
	}

	user := ctx.Value(AuthenticatedSubjectKey).(string)

	assetClassUuid, err := uuid.Parse(class.Id)
	if err != nil {
		grpclog.Errorf("Invalid UUID format for asset class ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid asset class ID format: %v", err)
	}

	updatedClass, err := client.AssetClass.UpdateOneID(assetClassUuid).
		SetName(class.Name).
		SetDescription(class.Description).
		SetOrder(int(class.Order)).
		SetProvider(class.Provider).
		SetIcon(class.Icon).
		SetColor(class.Color).
		SetUpdatedBy(user).
		Save(ctx)
	if err != nil {
		grpclog.Errorf("Failed to update asset class: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update asset class: %v", err)
	}
	log.S.Debugw("Updated asset class", "id", updatedClass.ID, "name", updatedClass.Name)
	return &gw.AssetClass{
		Id:          updatedClass.ID.String(),
		Name:        updatedClass.Name,
		Description: updatedClass.Description,
		Order:       int32(updatedClass.Order),
		Icon:        updatedClass.Icon,
		Color:       updatedClass.Color,
	}, nil
}

func (a assetClassServer) DeleteAssetClass(ctx context.Context, elementId *gw.ElementId) (*gw.EmptyMessage, error) {
	client, err := store.GetClient()
	if err != nil {
		grpclog.Errorf("Failed to get store client: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get store client: %v", err)
	}

	assetClassUuid, err := uuid.Parse(elementId.Id)
	if err != nil {
		grpclog.Errorf("Invalid UUID format for asset class ID: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid asset class ID format: %v", err)
	}

	user := ctx.Value(AuthenticatedSubjectKey).(string)

	_, err = client.AssetClass.UpdateOneID(assetClassUuid).
		SetDeletedAt(time.Now()).
		SetUpdatedBy(user).
		Save(ctx)
	if err != nil {
		grpclog.Errorf("Failed to delete asset class: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete asset class: %v", err)
	}
	log.S.Debugw("Deleted asset class", "id", assetClassUuid)

	return &gw.EmptyMessage{}, nil
}

func NewAssetClassServer() gw.AssetClassServiceServer {
	return &assetClassServer{}
}
