package api

import (
	"context"
	"encore.dev/storage/sqldb"
	"errors"
)

type User struct {
	Name  string
	Email string
	Role  string
}

//encore:api public method=POST path=/users/create
func createUser(ctx context.Context, user *User) error {
	if user.Name == "" || user.Email == "" || user.Role == "" {
		return errors.New("Name, Email and Role are required fields")
	}

	if user.Role != "Admin" && user.Role != "Modifier" && user.Role != "Watcher" {
		return errors.New("Role diferent!")
	}

	existingUser, err := existingUser(ctx, user.Name, user.Email, user.Role)
	if err != nil {
		return err
	}
	if existingUser {
		return errors.New("User already exists, try again!")
	}

	err = inserteUsers(ctx, user.Name, user.Email, user.Role)
	if err != nil {
		return err
	}
	return errors.New("User created successfully")
}

//encore:api public method=POST path=/users/delete
func deleteUser(ctx context.Context, user *User) error {
	if user.Name == "" || user.Email == "" || user.Role == "" {
		return errors.New("Name, Email and Role are required fields")
	}
	err := deleteUsers(ctx, user.Name, user.Email, user.Role)
	if err != nil {
		return err
	}
	return nil
}

func existingUser(ctx context.Context, name string, email string, role string) (bool, error) {
	var count int
	err := sqldb.QueryRow(ctx, `
		SELECT COUNT(*) FROM users WHERE name = $1 AND email = $2 AND role = $3
	`, name, email, role).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func deleteUsers(ctx context.Context, name string, email string, role string) error {
	_, err := sqldb.Exec(ctx, `
		DELETE FROM users WHERE name = $1 AND email = $2 AND role = $3
	`, name, email, role)
	if err != nil {
		return err
	}

	return nil
}

func inserteUsers(ctx context.Context, name string, email string, role string) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO users (name, email, role)
		VALUES ($1, $2, $3)
	`, name, email, role)
	return err
}
