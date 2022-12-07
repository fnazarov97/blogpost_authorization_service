package postgres

import (
	"blockpost/genprotos/authorization"
	"errors"
	"time"
)

// AddUser ...
func (p Postgres) AddUser(id string, req *authorization.CreateUserRequest) error {

	_, err := p.DB.Exec(`Insert into "user"("id", "username", "password", "user_type", created_at) 
						 VALUES($1, $2, $3, $4, now())
						`, id, req.Username, req.Password, req.UserType)
	if err != nil {
		return err
	}
	return nil
}

// GetUserByID ...
func (p Postgres) GetUserByID(id string) (*authorization.User, error) {
	res := &authorization.User{}
	var deletedAt *time.Time
	var updatedAt *string
	err := p.DB.QueryRow(`SELECT 
		"id",
		"username",
		"password",
		"user_type",
		"created_at",
		"updated_at",
		"deleted_at"
    FROM "user" WHERE id = $1`, id).Scan(
		&res.Id,
		&res.Username,
		&res.Password,
		&res.UserType,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("user not found")
	}

	return res, err
}

// GetUserList ...
func (p Postgres) GetUserList(offset, limit int, search string) (*authorization.GetUserListResponse, error) {
	resp := &authorization.GetUserListResponse{
		Users: make([]*authorization.User, 0),
	}
	rows, err := p.DB.Queryx(`SELECT
		"id",
		"username",
		"password",
		"user_type",
		"created_at",
		"updated_at"
	FROM "user" WHERE deleted_at IS NULL AND ("username" ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		a := &authorization.User{}

		var updatedAt *string

		err := rows.Scan(
			&a.Id,
			&a.Username,
			&a.Password,
			&a.UserType,
			&a.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		resp.Users = append(resp.Users, a)
	}

	return resp, err
}

// UpdateUser ...
func (p Postgres) UpdateUser(entity *authorization.UpdateUserRequest) error {

	res, err := p.DB.NamedExec(`UPDATE "user" SET "password"=:p, "updated_at"=now() WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": entity.Id,
		"p":  entity.Password,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("user not found")
}

// DeleteUser ...
func (p Postgres) DeleteUser(id string) error {
	res, err := p.DB.Exec(`UPDATE "user" SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("user not found")
}

// GetUserByUsername ...
func (p Postgres) GetUserByUsername(username string) (*authorization.User, error) {
	res := &authorization.User{}
	var deletedAt *time.Time
	var updatedAt *string
	err := p.DB.QueryRow(`SELECT 
		"id",
		"username",
		"password",
		"user_type",
		"created_at",
		"updated_at",
		"deleted_at"
    FROM "user" WHERE "username" = $1`, username).Scan(
		&res.Id,
		&res.Username,
		&res.Password,
		&res.UserType,
		&res.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if deletedAt != nil {
		return res, errors.New("user not found")
	}

	return res, err
}
