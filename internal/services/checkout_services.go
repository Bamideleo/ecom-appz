package services

import (
	"database/sql"
	"ecom-appz/internal/models"
	"ecom-appz/internal/repositories"
	"errors"
)

type CheckoutService struct {
	DB *sql.DB
	CartRepo repositories.CartRepository
	OrderRepo repositories.OrderRepository
	ProductRepo repositories.ProductRepository
}


func NewCheckoutService(
	cart repositories.CartRepository,
	order repositories.OrderRepository,
	product repositories.ProductRepository,
) *CheckoutService{
	return &CheckoutService{
		CartRepo: cart,
		OrderRepo: order,
		ProductRepo: product,
	}
}

func (s *CheckoutService) Checkout(userID string) (*models.Order, error){
	tx, err := s.DB.Begin()
	if err !=nil{
		return nil, err
	}
	
	defer func(){
		if err != nil{
			tx.Rollback()
		}
	}()

	cart, err:= s.CartRepo.GetCartWithItems(userID)
	if err !=nil{
		return nil, err
	}
	if len(cart.Items) == 0{
		return nil, errors.New("carts is empty")
	}
	var total float64
	var orderItems []models.OrderItem
	for _, item := range cart.Items{
		product, err := s.ProductRepo.FindByID(item.ProductID)
		if err != nil{
			return nil, err
		}
		// vlidate stock
		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for product")
		}
		total += float64(item.Quantity) * product.Price
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity: item.Quantity,
			Price: product.Price,
		})
	}
	order := models.Order{
		UserID: userID,
		Status: models.OrderPending,
		TotalAmount: total,
	}
	orderID, err := s.OrderRepo.CreateTx(tx, &order)
	if err != nil{
		return nil, err
	}

	// save order items
	err = s.OrderRepo.AddOrderItems(orderID, orderItems)
	if err !=nil{
		return nil, err
	}
	// Deduct stock
	for _, item :=range cart.Items{
		err:= s.ProductRepo.DeductStockTx(tx, item.ProductID, item.Quantity)
		if err != nil{
			return  nil, err
		}
	}
	// Clear cart
	err = s.CartRepo.ClearCartTx(tx, cart.ID)
	if err != nil{
		return nil, err
	}

	err = tx.Commit()
	if err != nil{
		return nil, err
	}
	
	order.ID = orderID
	order.Items = orderItems

	return  &order, nil
}