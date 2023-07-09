package repository

import (
	"database/sql"
	"scheduler/account"
	"scheduler/repository/database"
	"scheduler/util"
)

type UserRepository struct {
	db  *sql.DB
	acc *sqliteAccRepo
}

var _ Repository[account.User, util.UUID] = &UserRepository{}

func UserRepo(a AccountRepository) *UserRepository {
	return &UserRepository{
		db:  nil,
		acc: a.(*sqliteAccRepo),
	}
}

func (a *UserRepository) Get(id util.UUID) (u account.User, err error) {
	u.Account, err = a.acc.Get(id)
	if err != nil {
		return u, ErrUserDoesNotExist
	}
	err = a.db.QueryRow(`SELECT Name FROM User WHERE AccountID=?`, u.Meta.ID).Scan(
		&u.Name,
	)
	if err != nil {
		return u, ErrUserDoesNotExist
	}
	return
}

func (a *UserRepository) Save(u *account.User) error {
	return a.acc.save(&u.Account, func(tx *database.Conn, id int) error {
		if id < 0 {
			return ErrUserAlreadyExist
		}
		stmt := tx.Prepare(`INSERT INTO User (Name, AccountID) VALUES (?, ?)`)
		defer stmt.Reset()

		stmt.BindText(1, u.Name)
		stmt.BindInt64(2, int64(id))

		if r, err := stmt.Step(); err != nil {
			return ErrUserCanNotBeCreated.Wrap(err)
		} else if !r || stmt.ColumnInt(0) != 1 {
			return ErrAccountDoesNotExist
		} else {
			return nil
		}
	})
}

func (a *UserRepository) Update(u *account.User) error {
	acc := a.acc.GetID(u.UUID)
	if acc < 0 {
		return ErrUserDoesNotExist
	}
	_, err := a.db.Exec(`UPDATE User SET Name=?`, u.Name)
	return err
}

func (a *UserRepository) Delete(id util.UUID) error {
	return a.acc.Delete(id)
}
