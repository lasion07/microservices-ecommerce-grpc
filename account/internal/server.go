package internal

import (
	"context"
	"fmt"
	"net"

	"github.com/rasadov/EcommerceAPI/account/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(service Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()

	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		service:                           service})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (server *grpcServer) Register(ctx context.Context, request *pb.RegisterRequest) (*wrapperspb.StringValue, error) {
	token, err := server.service.Register(ctx, request.Name, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{
		Value: token,
	}, nil
}

func (server *grpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*wrapperspb.StringValue, error) {
	token, err := server.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{
		Value: token,
	}, nil
}

func (server *grpcServer) GetAccount(ctx context.Context, r *wrapperspb.UInt64Value) (*pb.AccountResponse, error) {
	a, err := server.service.GetAccount(ctx, r.Value)
	if err != nil {
		return nil, err
	}
	return &pb.AccountResponse{Account: &pb.Account{
		Id:   uint64(a.ID),
		Name: a.Name,
	}}, nil
}

func (server *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := server.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var accounts []*pb.Account
	for _, p := range res {
		accounts = append(accounts, &pb.Account{
			Id:   uint64(int(p.ID)),
			Name: p.Name,
		},
		)
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
