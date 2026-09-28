package postgresql

import (
	"context"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

func (p *PostgreSQL) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := p.conn.QueryContext(ctx, `
		SELECT category_id, name
		FROM category`)
	if err != nil {
		return nil, err
	}

	var categories []domain.Category
	for rows.Next() {
		var category domain.Category
		rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}

		categories = append(categories, category)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("scan error: %v", err)
	}

	return categories, nil
}
