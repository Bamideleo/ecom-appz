package repositories

import (
	"database/sql"
	"ecom-appz/internal/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll()([]models.Product, error)
	FindByID(id int)(*models.Product, error)
	Update(product *models.Product)error
	Delete(id int)error
}


type productRepository struct{
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository{
	return &productRepository{DB: db}
}

func (r *productRepository) Create(product *models.Product) error  {
	query :=`INSERT INTO products(name, description, price, stock)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`
	return r.DB.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
} 

func (r *productRepository) FindAll()([]models.Product, error){
	rows, err := r.DB.Query(`
	SELECT id, name, description, price, stock, created_at, updated_at FROM
	products 
	`)
	if err !=nil{
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next(){
		var p models.Product
		if err:=rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err !=nil{
			return nil, err
		}
		products = append(products, p)
	}
	return  products, nil
}
func (r *productRepository) FindByID(id int)(*models.Product, error){
	var p models.Product

	query := `
	SELECT id, name description, price, stock, created_at, updated_at
	FROM products WHERE id = $1
	`
	err := r.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.Stock,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	return &p, err
}

func (r *productRepository) Update(product *models.Product)error{
	query := `
	UPDATE products SET name=$1, description=$2, price=$3, stock=$4,
	updated_at=NOW()
	WHERE id =$5
	`
	_, err :=r.DB.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ID,
	)

	return err
}

func(r *productRepository) Delete(id int) error{
	_, err := r.DB.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}

