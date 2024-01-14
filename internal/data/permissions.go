package data

import (
	"context"
	"database/sql"
	"time"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i :=type range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT p.code
		FROM premissions p
		INNER JOIN users_permissions up ON up.permission_id = p.id
		INNER JOIN user u ON up.user_id = u.id
		WHERE u.id = $id`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string

		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}