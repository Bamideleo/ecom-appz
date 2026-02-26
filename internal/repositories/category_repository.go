package repositories

import (
	"database/sql"
	"ecom-appz/internal/models"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	FindAll() ([]models.Category, error)
	Update(category *models.Category) error
	Delete(id int) error

	AttachProduct(categoryID, productID int) error
	DetachProduct(categoryID, productID int) error
}

type categoryRepository struct{
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository{
	return &categoryRepository{DB:db}
}

func (r *categoryRepository) Create(category *models.Category) error{
	query := `
	INSERT INTO categories (name)
	VALUES ($1)
	RETURNING id, created_at, updated_at
	`
	return  r.DB.QueryRow(query, category.Name).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *categoryRepository) FindAll() ([]models.Category, error){
	rows, err := r.DB.Query(`
	SELECT id, name, created_at, update_at
	FROM categories
	`)
	if err != nil{
		return  nil, err
	}

	defer rows.Close()
	var categories []models.Category

	for rows.Next(){
		var c models.Category
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil{
			return  nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}


func (r *categoryRepository) Update(category *models.Category) error{
	_, err := r.DB.Exec(`
	UPDATE categories SET name=$1, update_at=NOW()
	WHERE id=$2
	`, category.Name, category.ID)
	return  err
}

func (r *categoryRepository) Delete(id int) error{
	_, err := r.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	return err
}

func (r *categoryRepository) AttachProduct(categoryID, productID int)error{
	_, err := r.DB.Exec(`
	INSERT INTO product_categories (product_id, category_id)
	VALUES ($1, $2)
	`, productID, categoryID)
	return err
}

func (r *categoryRepository) DetachProduct(categoryID, productID int) error{
	_, err := r.DB.Exec(`
		DELETE FROM product_categories
		WHERE product_id =$1 AND category_id=$2
	`, productID, categoryID)

	return err
}