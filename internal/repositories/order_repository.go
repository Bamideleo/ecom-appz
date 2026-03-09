package repositories

import (
	"database/sql"
	"ecom-appz/internal/models"
)

type OrderRepository interface {
	Create(order *models.Order) (int, error)
	AddOrderItems(oderID int, items []models.OrderItem)error
	GetByID(orderID int) (*models.Order, error)
	GetUserOrders(userId string) ([]models.Order, error)
}


type orderRepository struct{
	DB *sql.DB
}


func NewORdeRepository(db *sql.DB)OrderRepository{
	return  &orderRepository{DB: db}
}

func (r *orderRepository) Create(order *models.Order)(int, error){
	var orderID int
	err := r.DB.QueryRow(`
	INSERT INTO orders (user_id, status, total_amount)
	VALUES ($1, $2, $3)
	RETURNING id
	`,
	order.UserID,
	order.Status,
	order.TotalAmount,
).Scan(&orderID)

	if err != nil{
		return 0, err
	}
	return  orderID, nil
}

func (r *orderRepository) AddOrderItems(orderID int, items []models.OrderItem) error{
	for _, item := range items{
		_, err := r.DB.Exec(`
			INSERT INTO order_items
			(order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`,
		orderID,
		item.ProductID,
		item.Quantity,
		item.Price,
	)

	if err != nil{
		return err
	}
	}
	return  nil
}

func (r *orderRepository) GetByID(orderID int) (*models.Order, error){
	var order models.Order

	err := r.DB.QueryRow(`
	SELECT id, user_id, status, total_amount, created_at, updated_at
	FROM orders
	WHERE id=$1
	`, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.CreateAt,
		&order.UpdatedAt,
	)
	if err !=nil{
		return nil, err
	}
	rows, err := r.DB.Query(`
	SELECT product_id, quantity, price FROM order_items
	WHERE order_id=$1
	`, orderID)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var item models.OrderItem
		err := rows.Scan(
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)

		if err != nil{
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return &order, nil
}
func (r *orderRepository) GetUserOrders(userId string)([]models.Order, error){
	rows, err := r.DB.Query(`
	SELECT id, status, total_amount, created_at, updated_at
	FROM orders
	WHERE user_id = $1
	ORDER BY created_at DESC
	`, userId)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	var orders []models.Order
	for rows.Next(){
	var order models.Order
	err := rows.Scan(
		&order.ID,
		&order.Status,
		&order.TotalAmount,
		&order.CreateAt,
		&order.UpdatedAt,
	)
	
	if err !=nil{
		return  nil, err
	}

	order.UserID = userId
	// fetch order items
	items, err :=r.getOrderItems(order.ID)
	if err != nil{
		return nil, err
	}

	order.Items =items
	orders = append(orders, order)

	}
	return  orders, nil
}

func (r* orderRepository) getOrderItems(orderID int) ([]models.OrderItem, error){
	rows, err := r.DB.Query(`
	SELECT product_id, quantity, price
	FROM order_items
	WHERE order_id = $1
	`,orderID)

	if err !=nil{
		return nil, err
	}
	defer rows. Close()
	var items []models.OrderItem

	for rows.Next(){
		var item models.OrderItem

		err := rows.Scan(
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err !=nil{
			return nil, err
		}
		item.OrderID = orderID
		items = append(items, item)
	}
	return items, nil
}

