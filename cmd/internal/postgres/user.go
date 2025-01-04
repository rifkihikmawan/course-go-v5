package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"new-go-project/cmd/service"
	"strings"
)

func NewUserStore(
	db *Client,
) service.UserStore {
	return &userStore{
		db: db,
	}
}

type userStore struct {
	db *Client
}

// transactional ke database:
// BEGIN
// EXECUTION
// COMMIT

func (s *userStore) CreateUser(ctx context.Context, user *service.User) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (
		id,
		first_name,
		middle_name,
		last_name,
		type,
		status
	) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.MiddleName,
		user.LastName,
		user.Type,
		user.Status,
	)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *userStore) GetUsers(ctx context.Context) ([]service.User, error) {
	query := `SELECT
		id,
		first_name,
		middle_name,
		last_name,
		type,
		status
	FROM users`

	rows, err := s.db.Leader.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []service.User
	for rows.Next() {
		var user service.User
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.MiddleName,
			&user.LastName,
			&user.Type,
			&user.Status,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *userStore) GetUserById(ctx context.Context, id string) (*service.User, error) {
	var user service.User
	err := s.db.Leader.QueryRowContext(ctx, `
		SELECT id, first_name, middle_name, last_name, type, status
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Type,
		&user.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	return &user, nil
}

func (s *userStore) ActivateUser(ctx context.Context, id string) (*service.User, error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update status to active
	_, err = tx.ExecContext(ctx, `
		UPDATE users 
		SET status = $1
		WHERE id = $2
	`, service.UserStatusActive, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update user status: %w", err)
	}

	// Get updated user
	var user service.User
	err = tx.QueryRowContext(ctx, `
		SELECT id, first_name, middle_name, last_name, type, status
		FROM users
			WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Type,
		&user.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

func (s *userStore) DeleteUser(ctx context.Context, id string) error {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
		DELETE FROM users 
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("user not found")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *userStore) UpdateUserName(ctx context.Context, id string, firstName, middleName, lastName *string) (*service.User, error) {
	tx, err := s.db.Leader.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Prepare update fields
	updates := map[string]interface{}{
		"updated_at": "NOW()",
	}
	if firstName != nil {
		updates["first_name"] = *firstName
	}
	if middleName != nil {
		updates["middle_name"] = middleName
	}
	if lastName != nil {
		updates["last_name"] = *lastName
	}

	// Build query
	var fields []string
	var values []interface{}
	paramCount := 1

	for field, value := range updates {
		if str, ok := value.(string); ok && str == "NOW()" {
			fields = append(fields, field+" = "+str)
		} else {
			fields = append(fields, fmt.Sprintf("%s = $%d", field, paramCount))
			values = append(values, value)
			paramCount++
		}
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(fields, ", "),
		paramCount,
	)
	values = append(values, id)

	result, err := tx.ExecContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		return nil, nil
	}

	// Get updated user
	var user service.User
	err = tx.QueryRowContext(ctx, `
		SELECT id, first_name, middle_name, last_name, type, status, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.MiddleName,
		&user.LastName,
		&user.Type,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	return &user, tx.Commit()
}
