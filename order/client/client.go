package client

import (
	"context"
	"log"
	"time"

	"github.com/rasadov/EcommerceAPI/order/models"
	"github.com/rasadov/EcommerceAPI/order/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (client *Client) Close() {
	err := client.conn.Close()
	if err != nil {
		log.Println(err)
	}
}

func (client *Client) PostOrder(
	ctx context.Context,
	accountID uint64,
	products []*models.OrderedProduct,
) (*models.Order, error) {
	var protoProducts []*pb.OrderProduct
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.OrderProduct{
			Id:       p.ID,
			Quantity: p.Quantity,
		})
	}

	r, err := client.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}
	// Create response order
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	err = newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &models.Order{
		ID:         uint(r.Order.GetId()),
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (client *Client) GetOrdersForAccount(ctx context.Context, accountID uint64) ([]models.Order, error) {
	r, err := client.service.GetOrdersForAccount(ctx, &wrapperspb.UInt64Value{
		Value: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create response orders
	var orders []models.Order
	for _, orderProto := range r.Orders {
		newOrder := models.Order{
			ID:         uint(orderProto.Id),
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		err = newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		if err != nil {
			return nil, err
		}

		var products []*models.OrderedProduct
		for _, p := range orderProto.Products {
			products = append(products, &models.OrderedProduct{
				ID:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}

func (client *Client) UpdateOrderStatus(ctx context.Context, orderId uint64, status string) error {
	_, err := client.service.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		OrderId: orderId,
		Status:  status,
	})

	if err != nil {
		return err
	}

	return nil
}
