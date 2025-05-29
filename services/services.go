package services

/*func runServer(ctx context.Context, network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			grpclog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	s := grpc.NewServer()
	gw.RegisterYourServiceServer(s, NewEchoServer())

	go func() {
		defer s.GracefulStop()
		<-ctx.Done()
	}()
	return s.Serve(l)
}*/
