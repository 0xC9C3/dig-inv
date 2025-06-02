package services

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MockServerTransportStream struct {
	header []metadata.MD
}

func (m *MockServerTransportStream) Method() string {
	return "foo"
}

func (m *MockServerTransportStream) SetHeader(md metadata.MD) error {
	m.header = append(m.header, md)
	return nil
}

func (m *MockServerTransportStream) SendHeader(_ metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) SetTrailer(_ metadata.MD) error {
	return nil
}

func (m *MockServerTransportStream) GetHeaders() []metadata.MD {
	return m.header
}

func AddMockServerTransportStreamToContext(ctx context.Context, kv ...string) context.Context {
	ctx = grpc.NewContextWithServerTransportStream(ctx, &MockServerTransportStream{})
	return runtime.NewServerMetadataContext(ctx, runtime.ServerMetadata{
		HeaderMD: metadata.Pairs(kv...),
	})
}
