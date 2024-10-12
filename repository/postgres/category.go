package postgres

import (
	"database/sql"
	"errors"
	"financeapp/domain/category"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepo {
	return &CategoryRepo{
		db: db,
	}
}

func (cr CategoryRepo) Add(c *category.Category) (*category.Category, error) {
	tx, err := cr.db.Beginx()
	defer tx.Rollback()
	if err != nil {
		return nil, err
	}
	_, err = tx.Exec(
		"INSERT INTO categories (id, user_id, name) VALUES ($1, $2, $3)",
		c.ID, c.UserID, c.Name,
	)
	if err != nil {
		pgxError, ok := err.(*pgconn.PgError)
		if ok {
			switch pgxError.Code {
			case "23505":
				return nil, category.ErrorDuplicateCategoryName
			}
		}
		return nil, err
	}
	tx.Commit()
	return c, nil
}

func (cr CategoryRepo) GetByUserID(userID uuid.UUID) ([]category.Category, error) {
	rows, err := cr.db.Queryx("SELECT * FROM categories WHERE user_id=$1", userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, category.ErrorCategoryNotFound
		}
		return nil, err
	}
	categories := make([]category.Category, 0)
	for rows.Next() {
		c := category.Category{}
		err := rows.Scan(&c.ID, &c.UserID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (cr CategoryRepo) GetByID(userID, categoryID uuid.UUID) (*category.Category, error) {
	row := cr.db.QueryRowx(
		"SELECT * FROM categories WHERE id=$1 and user_id=$2",
		categoryID,
		userID,
	)
	c := category.Category{}
	err := row.Scan(&c.ID, &c.UserID, &c.Name)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
