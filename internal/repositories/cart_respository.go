package repositories

import (
	"database/sql"
	"errors"
)

type CartRepository interface {
	GetOrCreateCart(userID string) (int, error)
	AddItem(cartID, productID, quantity int) error
	UpdateQuantity(cartID, productID, quantity int) error
	RemoveItem(cartID, productID int) error
}

type cartRepo struct {
	DB *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository{
	return  &cartRepo{DB: db}
}

func (r *cartRepo) GetOrCreateCart(userID string) (int, error){
	var cartID int
	err := r.DB.QueryRow(
		"SELECT id FROM carts WHERE user_id = $1",
		userID,
	).Scan(&cartID)

	if err == sql.ErrNoRows{
		err = r.DB.QueryRow(
			"INSERT INTO carts (user_id) VALUES ($1) RETURNING id",
			userID,
		).Scan(&cartID)
	}
	return  cartID, err
}

func (r *cartRepo)AddItem(cartID, productID, quantity int) error{
	_, err := r.DB.Exec(`
	INSERT INTO cart_items (cart_id, product_id, quantity)
	VALUES ($1, $2, $3)
	ON CONFLICT(cart_id, product_id)
	DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`,
	cartID, productID, quantity,
)
return err
}

func (r *cartRepo) UpdateQuantity(cartID, productID, quantity int)error{
	if quantity <= 0{
		return errors.New(
			"quantity must be greater than zero",
		)
	}

	_, err:= r.DB.Exec(
		"UPDATE cart_items SET quantity = $1 WHERE cart_id = $2 AND product_id = $3",
		quantity, cartID, productID,
	)
	return err
}
func (r *cartRepo) RemoveItem(cartID, productID int) error{
	_, err := r.DB.Exec(
		"DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2",
		cartID, productID,
	)
	return  err
}