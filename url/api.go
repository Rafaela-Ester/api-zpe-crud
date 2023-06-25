package api

import (
	"context"
	"encore.dev/storage/sqldb"
)

//type Role string

//const (
/*Admin    Role = "Admin"
Modifier Role = "Modifier"
Watcher  Role = "Watcher"*/

type User struct {
	Name  string
	Email string
	//Role  Role
}

// encore:api public method=POST path=/users/{name}/{email}
func createUser(ctx context.Context, user *User) error {
	err := insert(ctx, user.Name, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func insert(ctx context.Context, name string, email string) error {
	_, err := sqldb.Exec(ctx, `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
	`, name, email)
	return err
}
