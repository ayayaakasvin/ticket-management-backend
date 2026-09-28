package sqlite

import (
	"context"
	"fmt"

	"github.com/ayayaakasvin/oneflick-ticket/internal/domain"
)

func (s *SQLite) GetCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT category_id, name FROM category ORDER BY category_id`)
	if err != nil {
		return nil, fmt.Errorf("query categories: %w", err)
	}
	defer rows.Close()

	categories := make([]domain.Category, 0)
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, fmt.Errorf("scan category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate categories: %w", err)
	}
	return categories, nil
}
