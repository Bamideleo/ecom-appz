package repositories

import (
	"database/sql"
	"ecom-appz/internal/models"
	"errors"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll()([]models.Product, error)
	FindByID(id int)(*models.Product, error)
	Update(product *models.Product)error
	Delete(id int)error
	List(page, limit int, sort, order, search string) ([]models.Product, int, error)
	FindWithDetails(id int)(*models.Product, error)
	DeductStock(productID, quantity int) error
}


type productRepository struct{
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository{
	return &productRepository{DB: db}
}

func (r *productRepository) Create(product *models.Product) error  {
	query :=`INSERT INTO products(name, description, price, stock)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, created_at, updated_at
	`
	return r.DB.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
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
			&p.ImageURL,
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
		&p.ImageURL,
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
		product.ImageURL,
		product.ID,
	)

	return err
}

func(r *productRepository) Delete(id int) error{
	_, err := r.DB.Exec("DELETE FROM products WHERE id=$1", id)
	return err
}

func (r *productRepository) List(page, limit int, sort, order, search string) ([]models.Product, int, error){
offset := (page -1) * limit
// Basic validation to prevent SQl injection
allowedSort :=map[string]bool{
	"name": true,
	"price": true,
	"created_at":true,
}
if !allowedSort[sort]{
	sort = "created_at"
}
	if order != "asc"{
		order ="desc"
	}

	baseQuery :=`
	FROM products WHERE ($1 = '' OR name ILIKE '%' || $1 || '%')
	`
	// Count total
	var total int
	countQuery := "SELECT COUNT(*)" + baseQuery
	err := r.DB.QueryRow(countQuery, search).Scan(&total)
	if err !=nil{
		return nil, 0, err
	}

	query :=`
	SELECT id, name, description, price, stock, image_url, created_at, updated_at
	` + baseQuery + `ORDER BY` + sort + ` ` + order + `
	LIMIT $2 OFFSET $3
	`
	rows, err := r.DB.Query(query, search, limit, offset)
	if err != nil{
		return nil, 0, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next(){
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Stock,
			&p.ImageURL,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil{
			return  nil, 0, err
		}
		products = append(products, p)
	}

	return products, total, nil
}

func (r *productRepository) FindWithDetails(id int) (*models.Product, error){
	var product models.Product
	query :=`
	SELECT id, name, description, price, stock, image_url, created_at, 
	updated_at FROM products
	WHERE id = $1
	` 
	err := r.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	// Load categories
	rows, err := r.DB.Query(`
	SELECT c.id, c.name, c.created_at, c.update_at
	FROM categories c
	JOIN product_categories pc ON c.id = pc.category_id
	WHERE pc.product_id = $1
	`, id)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var c models.Category
		if err:= rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err !=nil{
			return nil, err
		}
		product.Categories = append(product.Categories, c)
	}
	return &product, nil
}

func (r *productRepository) DeductStock(productID, quantity int)error{
	query := `
	UPDATE products
	SET stock = stock -$1
	WHERE id = $2 AND stock >= $1
	`

	result, err := r.DB.Exec(query, quantity, productID)
	if err != nil{
		return  err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}
	if rowsAffected == 0{
		return  errors.New("out of stock")
	}
	return nil
}

