package repository

import (
	"database/sql"
	"time"

	"github.com/david-kartopranoto/go-base/entity"
)

//UserSQL repo
type UserSQL struct {
	db *sql.DB
}

//NewUserSQL create new repository
func NewUserSQL(db *sql.DB) *UserSQL {
	return &UserSQL{
		db: db,
	}
}

//Create an user
func (r *UserSQL) Create(e *entity.User) (int64, error) {
	var newID int64
	err := r.db.QueryRow(`
		INSERT INTO app_user (email, password, username, created_at) 
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		e.Email,
		e.Password,
		e.Username,
		time.Now().Format("2006-01-02")).Scan(&newID)

	return newID, err
}

//Get an user
func (r *UserSQL) Get(id int64) (*entity.User, error) {
	return getUser(id, r.db)
}

func getUser(id int64, db *sql.DB) (*entity.User, error) {
	var u entity.User
	err := db.QueryRow("select id, email, password, username, created_at, updated_at from app_user where id = $1",
		id).Scan(&u.ID, &u.Email, &u.Password, &u.Username, &u.CreatedAt, &u.UpdatedAt)

	return &u, err
}

//Update an user
func (r *UserSQL) Update(e *entity.User) error {
	e.UpdatedAt = sql.NullTime{Time: time.Now()}
	_, err := r.db.Exec("update app_user set email = $1, password = $2, username = $3, updated_at = $4 where id = $5",
		e.Email, e.Password, e.Username, e.UpdatedAt.Time.Format("2006-01-02"), e.ID)
	if err != nil {
		return err
	}
	return nil
}

//Search users
func (r *UserSQL) Search(query string) ([]*entity.User, error) {
	rows, err := r.db.Query("SELECT id, email, password, username, created_at, updated_at FROM app_user where username like $1", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Username, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

//List users
func (r *UserSQL) List() ([]*entity.User, error) {
	rows, err := r.db.Query("SELECT id, email, password, username, created_at, updated_at FROM app_user")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Username, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

//Delete an user
func (r *UserSQL) Delete(id int64) error {
	_, err := r.db.Exec("delete from app_user where id = $1", id)
	return err
}
