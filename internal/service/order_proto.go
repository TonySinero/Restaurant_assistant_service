package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"restaurant-assistant/internal/domain"
	courierProto "restaurant-assistant/pkg/courierProto"
	"restaurant-assistant/pkg/orderservice_fd"
	"time"
)

func (s *OrderService) CreateConnectionCS() (courierProto.CourierServerClient, *grpc.ClientConn, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", s.cfg.GRPCCS.Host, s.cfg.GRPCCS.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("error occurred while creating conn to CS")
		return nil, nil, ctx, err
	}

	orderCreateClient := courierProto.NewCourierServerClient(conn)

	return orderCreateClient, conn, ctx, nil
}

func (s *OrderService) CreateOrderCS(id string, input domain.UpdateOrder) error {
	orderServiceClient, conn, ctx, err := s.CreateConnectionCS()
	if err != nil {
		return err
	}

	order, err := s.repo.GetOrderByID(id)
	if err != nil {
		return err
	}

	restaurant, err := s.repo.GetRestaurantByID(id)
	if err != nil {
		return err
	}

	protoDate, err := ptypes.TimestampProto(order.DeliveryTime)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while updating DeliveryTime")
		return err
	}

	orderCS := &courierProto.OrderCourierServer{
		OrderID:           int64(order.RowID),
		CourierServiceID:  int64(*input.CourierService),
		RestaurantAddress: restaurant.Address,
		RestaurantName:    restaurant.Title,
		DeliveryTime:      protoDate,
		ClientAddress:     order.Address,
		ClientFullName:    order.ClientFullName,
		ClientPhoneNumber: order.ClientPhoneNumber,
		PaymentType:       int64(order.PaymentTypeID),
	}

	if _, err := orderServiceClient.CreateOrder(ctx, orderCS); err != nil {
		log.Error().Err(err).Msg("error occurred while creating order in courier service")
		return err
	}

	defer conn.Close()

	return nil
}

func (s *OrderService) GetAllOrderDeliveryServicesCS() (*courierProto.ServicesResponse, error) {
	orderServiceClient, conn, ctx, err := s.CreateConnectionCS()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	allDeliveryServices, err := orderServiceClient.GetDeliveryServicesList(ctx, &emptypb.Empty{})
	if err != nil {
		log.Error().Err(err).Msg("error occurred while selecting delivery services in Courier service")
		return allDeliveryServices, err
	}

	return allDeliveryServices, nil
}

func (s *OrderService) CreateConnectionFD() (orderservice_fd.OrderServiceFDClient, *grpc.ClientConn, context.Context, error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", s.cfg.GRPCFD.Host, s.cfg.GRPCFD.Port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Error().Err(err).Msg("error occurred while creating conn to FD")
		return nil, nil, ctx, err
	}

	orderUpdateClient := orderservice_fd.NewOrderServiceFDClient(conn)

	return orderUpdateClient, conn, ctx, nil
}

func (s *OrderService) UpdateOrderFD(id string, input domain.UpdateOrder) error {
	orderClientFD, conn, ctx, err := s.CreateConnectionFD()
	if err != nil {
		return err
	}

	updateOrder := &orderservice_fd.UpdateOrderMessage{
		OrderUUID: id,
		Status:    int64(*input.Status),
	}

	if _, err := orderClientFD.UpdateOrder(ctx, updateOrder); err != nil {
		log.Error().Err(err).Msg("error occurred while updating order in FD")
		return err
	}

	defer conn.Close()

	return nil
}
