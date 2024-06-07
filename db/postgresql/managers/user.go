package managers

import (
	"database/sql"

	pbu "user-service/services/genproto"

	"github.com/google/uuid"
)

type UserManager struct {
	conn *sql.DB
}

func NewUserManager(conn *sql.DB) *UserManager {
	return &UserManager{conn: conn}
}

func (um *UserManager) CreateUser(user *pbu.UserCreatedReq) (string, error) {
	newUUID := uuid.NewString()
	query := "INSERT INTO users (id, name, email, phone) VALUES ($1, $2, $3, $4)"
	_, err := um.conn.Exec(query, newUUID, user.Name, user.Email, user.Phone)
	if err != nil {
		return "", err
	}
	return newUUID, nil
}

func (um *UserManager) CreateCart(id string) error {
	query := "INSERT INTO cart(id, user_id, items) VALUES ($1, $2, '[]')"
	newUUID := uuid.NewString()
	_, err := um.conn.Exec(query, newUUID, id)
	return err
}

func (um *UserManager) GetUserByID(id string) (*pbu.UserGetByIDResp, error) {
	query := "SELECT id, name, email, phone, created_at, updated_at, deleted_at FROM users WHERE id = $1 AND deleted_at = 0"
	row := um.conn.QueryRow(query, id)
	var user pbu.UserGetByIDResp
	err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (um *UserManager) GetAllUsers() (*pbu.UsersGetAllResp, error) {
	query := "SELECT id, name, email, phone, created_at, updated_at, deleted_at FROM users WHERE deleted_at = 0"
	rows, err := um.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users pbu.UsersGetAllResp
	users.Total = 0
	for rows.Next() {
		var user pbu.UserGetByIDResp
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
		if err != nil {
			return nil, err
		}
		users.Users = append(users.Users, &user)
		users.Total++
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (um *UserManager) UpdateUser(user *pbu.UserUpdatedReq) error {
	query := "UPDATE users SET name = $1, email = $2, phone = $3, updated_at = NOW() WHERE id = $4 AND deleted_at = 0"
	_, err := um.conn.Exec(query, user.Name, user.Email, user.Phone, user.Id)
	return err
}

func (um *UserManager) DeleteUser(id string) error {
	query := "UPDATE users SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1 AND deleted_at = 0"
	_, err := um.conn.Exec(query, id)
	return err
}
