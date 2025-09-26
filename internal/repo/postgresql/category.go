package postgresql

import (
	"context"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/models"
)

func (p *PostgreSQL) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT category_id, name
		FROM category`)
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return  categories, nil
}