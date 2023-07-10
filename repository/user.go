package repository

import (
	"scheduler/account"
	"scheduler/repository/database"
	"scheduler/util"
)

// type User interface {
// 	Repository[account.User, util.UUID]
// 	VerifyPassword(UUID util.UUID, password account.Password) bool
// }

type User struct {
	db  *database.SQLite
	acc *Account
}

var _ Repository[account.User, util.UUID] = &User{}

func UserRepo(a *Account) User {
	return User{
		db:  a.db,
		acc: a,
	}
}

func (a *User) Get(id util.UUID) (u account.User, err error) {
	u.Account, err = a.acc.Get(id)
	if err != nil {
		return u, ErrUserDoesNotExist
	}
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`SELECT Name FROM User WHERE AccountID=?`)
	defer stmt.Reset()

	stmt.BindInt64(1, int64(u.Account.Meta.ID))
	if r, err := stmt.Step(); err != nil {
		return account.User{}, err
	} else if !r {
		return account.User{}, ErrUserDoesNotExist
	} else {
		u.Name = stmt.ColumnText(0)
	}
	u.Asset = a.getAsset(id, conn)
	return u, nil
}

func (a *User) getAsset(id util.UUID, conn *database.Conn) account.UserAsset {

	stmt := conn.Prepare(
		`SELECT
			UA.Balance, I.Name AS InstitutionName
		FROM UserAsset AS UA
		INNER JOIN Institution AS I ON UA.InstitutionID = I.ID
		WHERE UA.ID = ?`)
	defer stmt.Reset()

	stmt.BindText(1, id.Str())
	s := account.UserAsset{}
	for {
		if r, err := stmt.Step(); err != nil || !r {
			break
		} else {
			if s.Balance == nil {
				s.Balance = make(map[string]int, 5)
			}
			b := stmt.ColumnInt(0)
			s.Balance[stmt.ColumnText(1)] = b
			s.Sum += b
		}
	}
	return s
}

func (a *User) Save(u *account.User) error {
	return a.acc.save(&u.Account, func(tx *database.Conn, id int) error {
		if id < 0 {
			return ErrUserAlreadyExist
		}
		stmt := tx.Prepare(`INSERT INTO User (Name, AccountID) VALUES (?, ?) RETURNING ID`)
		defer stmt.Reset()

		stmt.BindText(1, u.Name)
		stmt.BindInt64(2, int64(id))

		if r, err := stmt.Step(); err != nil {
			return ErrUserAlreadyExist.Wrap(err)
		} else if !r || stmt.ColumnInt(0) < 0 {
			return ErrUserCanNotBeCreated
		} else {
			return nil
		}
	})
}

func (a *User) Update(u *account.User) error {
	acc := a.acc.GetID(u.UUID)
	if acc < 0 {
		return ErrUserDoesNotExist
	}
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`UPDATE User SET Name=? WHERE AccountID=? RETURNING ID`)
	defer stmt.Reset()

	stmt.BindInt64(1, int64(acc))
	if r, err := stmt.Step(); err != nil {
		return ErrUserDoesNotExist.Wrap(err)
	} else if !r {
		return ErrUserDoesNotExist
	} else {
		return nil
	}
}

func (a *User) Delete(id util.UUID) error {
	return a.acc.Delete(id)
}

func (a *User) VerifyPassword(id util.UUID, password account.Password) bool {
	return a.acc.VerifyPassword(id, password)
}
