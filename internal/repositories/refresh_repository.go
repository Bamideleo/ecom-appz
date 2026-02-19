package repositories

import (
	"database/sql"
	"ecom-appz/internal/models"
)

type RefreshRepository interface {
	Store(token *models.RefreshToken) error
	Delete(token string) error
	Find(token string) (*models.RefreshToken, error)
}


type refreshRepository struct{
	DB *sql.DB
}

func NewRefreshRepository(db *sql.DB) *refreshRepository{
	return &refreshRepository{DB: db}
}

func (r *refreshRepository) Store(token *models.RefreshToken)error{
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at)
	VALUES ($1, $2, $3)
	`
	_, err := r.DB.Exec(query, token.UserId, token.Token, token.ExpiresAt)
	return  err
}

func (r* refreshRepository) Delete(token string)error{
	_, err := r.DB.Exec("DELETE FROM refresh_token WHERE token=$1", token)
	return err
}

func (r *refreshRepository)Find(token string)(*models.RefreshToken, error){
	var t models.RefreshToken

	query :=`
	SELECT id, user_id, token, 
	expires_at FROM refresh_tokens WHERE token=$1`
	err := r.DB.QueryRow(query, token).Scan(&t.ID, &t.UserId, &t.Token, &t.ExpiresAt)
	return &t, err
}