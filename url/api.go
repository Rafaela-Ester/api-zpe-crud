package api

import (
	"context"
	"encore.dev/storage/sqldb"
	"errors"
	_ "gorm.io/gorm"
)

//estrutura dos parametros do usuario

type User struct {
	ID    int
	Name  string
	Email string
	Role  string
}

// estrutura dos parametros para excluir usuario

type UserDelete struct {
	ID   int
	Name string
}

// estrutura dos parametros para atualizar usuario

type UpdateUser struct {
	ID      int
	Name    string
	Email   string
	Role    string
	NewRole string
}

// estrutura dos parametros para ler um usuario

type ReadUser struct {
	ID       string
	ReadUser *User
}

// estrutura de retorno para o usuario

type CreateUserResponse struct {
	Message string
}

// implementado o ednpoint de criação de cadastro de usuario, utilizadno a estrutura de usuario que foi definida em cima

//encore:api public method=POST path=/users/create
func createUser(ctx context.Context, user *User) (*CreateUserResponse, error) {
	if user.Name == "" || user.Email == "" || user.Role == "" {
		return nil, errors.New("Name, Email and Role are required fields")
	}

	if user.Role != "Admin" && user.Role != "Modifier" && user.Role != "Watcher" {
		return nil, errors.New("Invalid role. Allowed roles are Admin, Modifier, and Watcher")
	}

	existingUser, err := existingUser(ctx, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	if existingUser {
		return nil, errors.New("User already exists!")
	}

	_, err = sqldb.Exec(ctx, `
		INSERT INTO users (name, email, role)
		VALUES ($1, $2, $3)
	`, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	response := &CreateUserResponse{
		Message: "User created successfully",
	}

	return response, nil
}

// Estrutura de leitura de usuario, ao inserir o id a api busca se existe ou nao, caso exista retorna na tela os dados do usuario.

//encore:api public method=GET path=/read/:id
func readUser(ctx context.Context, id string) (*ReadUser, error) {
	u := &ReadUser{ID: id}
	user := &User{}

	err := sqldb.QueryRow(ctx, `
		SELECT id, name, email, role FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil, errors.New("User already exists!")
	}

	u.ReadUser = user
	return u, nil
}

type ListUsersResponse struct {
	Users []*User
}

// Retorna os dados de todos os usuarios da api

//encore:api public method=GET path=/users
func ListUsers(ctx context.Context) (*ListUsersResponse, error) {
	rows, err := sqldb.Query(ctx, `
		SELECT id, name, email, role FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	response := &ListUsersResponse{
		Users: users,
	}

	return response, nil
}

type UpdateUserResponse struct {
	Message string
}

// endpoint desenvolvido para realizar o update dos usuarios, a estrtutura foi definida no inicio. Ao iformar
// os dados é preciso enviar a nova função para as verificações

//encore:api public method=POST path=/UpdateUser/update
func updateUser(ctx context.Context, update *UpdateUser) (*UpdateUserResponse, error) {
	if update.Name == "" || update.Email == "" || update.Role == "" || update.NewRole == "" {
		return nil, errors.New("Name, Email, Role and New Role are required fields")
	}

	if update.Role != "Admin" && update.Role != "Modifier" && update.Role != "Watcher" {
		return nil, errors.New("Invalid role. Allowed roles are Admin, Modifier, and Watcher")
	}

	if update.Role != update.NewRole {
		existingUser, err := existingUser(ctx, update.Name, update.Email, update.Role)
		if err != nil {
			return nil, err
		}
		if existingUser {

			if update.Role == "Admin" && (update.NewRole == "Watcher" || update.NewRole == "Modifier") {
				err = updateUsers(ctx, update.Name, update.Email, update.NewRole)
				if err != nil {
					return nil, err
				}
			} else if update.Role == "Modifier" && update.NewRole == "Watcher" {
				err = updateUsers(ctx, update.Name, update.Email, update.NewRole)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, errors.New("User cannot change roles.")
			}

		} else {
			return nil, errors.New("User does not exist")
		}
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("New Role must be different from the current role")
	}

	responseUpdate := &UpdateUserResponse{
		Message: "User updated successfully",
	}
	return responseUpdate, nil
}

type DeleteUserResponse struct {
	Message string
}

//Neste endpoint é recuperado os dados que o foram informados e feito a consulta no banco para a remoção do usuario.

//encore:api public method=POST path=/users/delete
func deleteUser(ctx context.Context, user *UserDelete) (*DeleteUserResponse, error) {
	if user.ID == 0 || user.Name == "" {
		return nil, errors.New("ID and Name are required fields!")
	}

	err, _ := sqldb.Exec(ctx, `
		DELETE FROM users WHERE id = $1 AND name = $2
	`, user.ID, user.Name)
	if err != nil {
		return nil, errors.New("User does not exist")
	}
	response := &DeleteUserResponse{
		Message: "User deleted successfully",
	}
	return response, nil

}

//Execuçoes no banco de dados, essas funções sao chamadas no corpo do endpoint. A funcao de "ExistingUser" é utilizada
//em todas as outras funcoes, ois pe feito a verificação em cada etapa se já existe cadastro.

func existingUser(ctx context.Context, name string, email string, role string) (bool, error) {
	var count int
	err := sqldb.QueryRow(ctx, `
		SELECT COUNT(*) FROM users WHERE name = $1 and  email = $2 and role = $3
	`, name, email, role).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func updateUsers(ctx context.Context, name string, email string, role string) error {
	_, err := sqldb.Exec(ctx, `
		UPDATE users SET role = $3
		WHERE name = $1 AND email = $2
	`, name, email, role)
	return err
}
