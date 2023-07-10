package repository

import (
	"scheduler/account"
	"scheduler/repository/database"
	"scheduler/util"
	"time"

	"github.com/bika-c/sqlite"
)

var _ Repository[account.Account, util.UUID] = &Account{}

type Account struct {
	db *database.SQLite
}

func AccountRepo(db *database.SQLite) Account {
	return Account{
		db: db,
	}
}

func (a *Account) Get(id util.UUID) (account.Account, error) {
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`SELECT Email, ID, CreatedAt, UpdatedAt, LastLogin FROM Account WHERE UUID=?`)
	defer stmt.Reset()

	var acc account.Account
	acc.UUID = id

	stmt.BindText(1, id.Str())
	if r, err := stmt.Step(); err != nil {
		return account.EmptyAccount, err
	} else if !r {
		return account.EmptyAccount, ErrAccountDoesNotExist
	} else {
		acc.Email = stmt.ColumnText(0)
		acc.Meta.ID = stmt.ColumnInt(1)
		acc.Meta.CreatedAt, _ = time.Parse(time.DateTime, stmt.ColumnText(2))
		acc.Meta.UpdatedAt, _ = time.Parse(time.DateTime, stmt.ColumnText(3))
		if stmt.ColumnType(4) == sqlite.SQLITE_TEXT {
			acc.Meta.LastLogin, _ = time.Parse(time.DateTime, stmt.ColumnText(4))
		}
	}
	return acc, nil
}

func (a *Account) GetID(uuid util.UUID) (i int) {
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`SELECT ID FROM Account WHERE UUID=?`)
	defer stmt.Reset()

	stmt.BindText(1, uuid.Str())
	if r, err := stmt.Step(); err != nil {
		return -1
	} else if !r {
		return -1
	} else {
		return stmt.ColumnInt(0)
	}
}

func (a *Account) exist(u *account.Account) bool {
	if u.Email == "" && u.UUID.IsUUID() {
		return true
	}
	if !u.UUID.IsUUID() {
		return false
	}

	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`SELECT count(UUID) FROM Account WHERE (UUID=? AND Email=?) OR ID=?`)
	defer stmt.Reset()
	stmt.BindText(1, u.UUID.Str())
	stmt.BindText(2, u.Email)
	stmt.BindInt64(3, int64(u.Meta.ID))
	if r, err := stmt.Step(); err != nil {
		return false
	} else if !r {
		return false
	} else {
		return stmt.ColumnInt(0) == 1
	}
}

func (a Account) Save(u *account.Account) error {
	return a.Update(u)
}

func (a *Account) save(u *account.Account, fn func(*database.Conn, int) error) (e error) {
	if !u.UUID.IsUUID() {
		u.UUID = util.NewUUID()
	}

	conn := a.db.Get()
	defer a.db.Put(conn)
	defer conn.Save()(&e)
	stmt := conn.Prepare(`INSERT INTO Account (UUID, Email, Password) VALUES (?, ?, ?) RETURNING ID`)
	defer stmt.Reset()

	stmt.BindText(1, u.UUID.Str())
	stmt.BindText(2, u.Email)
	stmt.BindText(3, u.Password.String())

	if r, err := stmt.Step(); err != nil {
		return ErrAccountAlreadyExist.Wrap(err)
	} else if !r {
		return ErrAccountCanNotBeCreated
	} else if fn != nil {
		return fn(conn, stmt.ColumnInt(0))
	} else {
		return nil
	}
}

func (a *Account) Update(u *account.Account) error {
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`UPDATE Account SET
			Email = CASE WHEN coalesce(:email, '') = '' THEN
				Email ELSE :email
			END,
			Password = CASE WHEN coalesce(:password, '') = '' THEN
				Password ELSE :password
			END
			WHERE UUID=:uuid RETURNING UUID`)
	defer stmt.Reset()
	stmt.SetText(":email", u.Email)
	stmt.SetText(":password", u.Password.String())
	stmt.SetText(":uuid", u.UUID.Str())

	if r, err := stmt.Step(); err != nil {
		return ErrAccountCanNotBeUpdated.Wrap(err)
	} else if !r {
		return ErrAccountDoesNotExist
	} else {
		return nil
	}
}

func (a *Account) Delete(id util.UUID) error {
	conn := a.db.Get()
	defer a.db.Put(conn)

	stmt := conn.Prepare(`DELETE FROM Account WHERE UUID=? RETURNING ID`)
	defer stmt.Reset()

	stmt.BindText(1, id.Str())

	if r, err := stmt.Step(); err != nil {
		return ErrAccountDoesNotExist
	} else if !r {
		return ErrAccountCanNotBeDeleted
	} else {
		return nil
	}
}

func (a *Account) VerifyPassword(id util.UUID, password account.Password) bool {
	conn := a.db.Get()
	defer a.db.Put(conn)
	stmt := conn.Prepare(`SELECT Password FROM Account WHERE UUID=?`)
	defer stmt.Reset()
	stmt.BindText(1, id.Str())

	if r, err := stmt.Step(); err != nil {
		return false
	} else if !r {
		return false
	} else {
		return password.Compare(account.Password(stmt.ColumnText(0)))
	}
}
