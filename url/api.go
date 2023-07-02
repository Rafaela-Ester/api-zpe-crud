package api

import (
	"context"
	"encore.dev/storage/sqldb"
	"errors"
	_ "gorm.io/gorm"
)

// Estrutura geral de usuario, com todos os dados a serem solicitados.

type User struct {
	ID    int
	Name  string
	Email string
	Role  string
}

// Estrutura para exclusão de usuario, apenas o id e o nome serão solicitados.

type UserDelete struct {
	ID   int
	Name string
}

// Essa estrutura é utilizada para solicitar os dados de um usuario, e uma nova variavel para que seja informado
// a nova função a ser atribuida.

type UpdateUser struct {
	ID      int
	Name    string
	Email   string
	Role    string
	NewRole string
}

// Estrutura dos parametros para ler um usuario, é solicitado o id do usuario que deseja se obter os dados. ReadUser é
// um ponteiro que recebe a estrutura user como parametro.

type ReadUser struct {
	ID       string
	ReadUser *User
}

// Estrutura de retorno para o usuario.

type UserResponse struct {
	Message string
}

// Implementado endpoint de criação de cadastro de usuario, utilizadno a estrutura de usuario que foi definida a cima.
// O endpoint verifica se existe valores nos campos de name, email e role.

//encore:api public method=POST path=/create/users
func createUser(ctx context.Context, user *User) (*UserResponse, error) {
	// Verifica se os campos obrigatorios foram preenchidos.
	if user.Name == "" || user.Email == "" || user.Role == "" {
		return nil, errors.New("Name, Email and Role are required fields")
	}
	// Verifica se a função foi preenchida corretamente.
	if user.Role != "Admin" && user.Role != "Modifier" && user.Role != "Watcher" {
		return nil, errors.New("Invalid role. Allowed roles are Admin, Modifier, and Watcher")
	}

	// Verifica se já existe um usuario com esses dados
	existingUser, err := existingUser(ctx, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}
	if existingUser {
		return nil, errors.New("User already exists!")
	}

	// Se não existir, executa essa query no banco que ira cadastrar o novo usuario.
	_, err = sqldb.Exec(ctx, `
        INSERT INTO users (name, email, role)
        VALUES ($1, $2, $3)
    `, user.Name, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	// Retorna que foi criado com sucesso.
	response := &UserResponse{
		Message: "User created successfully",
	}

	return response, nil
}

type ReadUserResponse struct {
	User *ReadUser
}

// Endpoint para ler usuario e retornar as informações dele.
//Caso não encontre o usuario é retornado a mensagem de usuario não existente.

//encore:api public method=GET path=/read/users/:id
func readUser(ctx context.Context, id string) (*ReadUser, error) {

	u := &ReadUser{ID: id}
	user := &User{}

	// Realiza a consulta no banco de dados para buscar o usuario conforme o id informado.

	err := sqldb.QueryRow(ctx, `
		SELECT id, name, email, role FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Role)
	if err != nil {
		return nil, errors.New("User already exists!")
	}

	// Retorna os dados do usuario.
	u.ReadUser = user
	return u, nil
}

// Endpoint para atualizar a função do usuario, é feito algumas verificações se os dados foram preenchidos, se
// a funcao foi preenchida corretamente e se a nova função é diferente da funcao atual para poder inicar as regras de troca.

//encore:api public method=POST path=/UpdateUser/update
func updateUser(ctx context.Context, update *UpdateUser) (*UserResponse, error) {
	// Verifica se os campos obrigatorios foram preenchidos.
	if update.Name == "" || update.Email == "" || update.Role == "" || update.NewRole == "" {
		return nil, errors.New("Name, Email, Role and New Role are required fields")
	}
	// Verifica se a função foi preenchida corretamente.
	if update.Role != "Admin" && update.Role != "Modifier" && update.Role != "Watcher" {
		return nil, errors.New("Invalid role. Allowed roles are Admin, Modifier, and Watcher")
	}

	// Verifica se a função e nova função são diferentes.
	if update.Role != update.NewRole {
		// Verifica se o usuario existe.
		existingUser, err := existingUser(ctx, update.Name, update.Email, update.Role)
		if err != nil {
			return nil, err
		}
		if existingUser {

			// Verifica se a nova função é admin e a nova função é watcher ou modifier.

			if update.Role == "Admin" && (update.NewRole == "Watcher" || update.NewRole == "Modifier") {
				err = updateUsers(ctx, update.Name, update.Email, update.NewRole)
				if err != nil {
					return nil, err
				}
				// Verifica se a função é modifier e se a nova função é watcher.
			} else if update.Role == "Modifier" && update.NewRole == "Watcher" {
				err = updateUsers(ctx, update.Name, update.Email, update.NewRole)
				if err != nil {
					return nil, err
				}
				// Se for diferente das duas regras acima, retorna que não pode alterar a função.
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
		// Retorna que a nova função deve ser diferente da função atual
		return nil, errors.New("New Role must be different from the current role")
	}

	responseUpdate := &UserResponse{
		Message: "User updated successfully",
	}
	return responseUpdate, nil
}

// Endpoint para deletar um usuario, é solicitado o id e o nome do usuario que deseja deletar. Após isso
//Neste endpoint é recuperado os dados que o foram informados e feito a consulta no banco para a exclusão do usuario.

//encore:api public method=POST path=/users/delete
func deleteUser(ctx context.Context, user *UserDelete) (*UserResponse, error) {
	// Verifica se os dados foram preenchidos
	if user.ID == 0 || user.Name == "" {
		return nil, errors.New("ID and Name are required fields!")
	}

	// Executa a query no banco de dados
	_, err := sqldb.Exec(ctx, `
		DELETE FROM users WHERE id = $1 AND name = $2
	`, user.ID, user.Name)
	if err != nil {
		// Se ocorrer erros, retorna que o usuario não existe.
		return nil, errors.New("User does not exist")
	}
	response := &UserResponse{
		Message: "User deleted successfully",
	}

	return response, nil
}

type ListUsersResponse struct {
	Users []*User
}

// Endpoint para listar todos os usuários.

//encore:api public method=GET path=/users
func listUsers(ctx context.Context) (*ListUsersResponse, error) {
	// Busca no banco de dados todos os usuarios.
	rows, err := sqldb.Query(ctx, `
		SELECT id, name, email, role FROM users
	`)
	if err != nil {
		return nil, errors.New("No users found.")
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
			return nil, err
		}
		// Enquanto encontrar dados, armazena na slice (array)
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.New("No users found.")
	}

	response := &ListUsersResponse{
		Users: users,
	}

	// Retorna todos os usuarios encontrados.
	return response, nil
}

//Execuçoes no banco de dados, essas funções sao chamadas no corpo do endpoint.

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
