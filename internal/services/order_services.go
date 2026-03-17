package services

import (
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"errors"
)

type OrderService struct {
	Repo repositories.OrderRepository
}

func(s *OrderService) UpdateOrderStatus(orderID int, newStatus models.OrderStatus) error{
	if !models.IsValidStatus(newStatus) {
		return errors.New("Invalid status")
	}

	order, err := s.Repo.GetByID(orderID)
	if err != nil{
		return err
	}

	if !models.CanTransition(order.Status, newStatus){
		return errors.New("invalid status transition")
	}
	return s.Repo.UpdateStatus(orderID, newStatus)
}