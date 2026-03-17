package models


type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderShipped   OrderStatus = "shipped"
	OrderCompleted OrderStatus = "completed"
	OrderCancelled OrderStatus = "cancelled"
)

func IsValidStatus(status OrderStatus) bool {
	switch status {
	case OrderPending,
		OrderPaid,
		OrderShipped,
		OrderCompleted,
		OrderCancelled:
		return true
	default:
		return false
	}
}

func CanTransition(from, to OrderStatus) bool {
	switch from {
	case OrderPending:
		return to == OrderPaid || to == OrderCancelled
	case OrderPaid:
		return to == OrderShipped || to == OrderCancelled
	case OrderShipped:
		return  to == OrderCompleted
	case OrderCompleted:
		return false
	case OrderCancelled:
		return false
		
	}
	return  false
}