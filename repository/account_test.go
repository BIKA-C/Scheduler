package repository

import (
	"fmt"
	"math/rand"
	"scheduler/account"
	"scheduler/repository/database"
	"scheduler/util"
	"testing"
)

func setupRepo(b *testing.B) (Account, []util.UUID) {
	b.Helper()
	db := database.NewSQLite(b.TempDir()+"/test.db", 40, database.DefaultPragma)
	database.Setup(db)
	a := Account{
		db: db,
	}
	uuid := make([]util.UUID, b.N)
	for i := 0; i < b.N; i++ {
		uuid[i] = util.NewUUID()
		acc := account.Account{
			Password: "12345678",
			Email:    fmt.Sprintf("acc%d@test.bench", i),
			UUID:     uuid[i],
		}
		if err := a.save(&acc, func(_ *database.Conn, id int) error {
			if id < 0 {
				b.Fatal("insert failed")
			}
			return nil
		}); err != nil {
			b.Fatal(err)
		}
	}
	c := db.Get()
	defer db.Put(c)
	s := c.PrepareTransient("SELECT count(*) FROM Account")
	defer s.Finalize()
	if r, err := s.Step(); !r || err != nil {
		b.Fatal(err, r, "Error in fetching count")
	} else if s.ColumnInt(0) != b.N {
		b.Fatal(s.ColumnInt(0))
	}

	return a, uuid
}

func BenchmarkAccRepo_Read(b *testing.B) {
	a, uuid := setupRepo(b)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := rand.Intn(b.N)
			if _, err := a.Get(uuid[i]); err != nil {
				b.Error("err:", err.Error())
			}
		}
	})
}

func BenchmarkAccRepo_Write(b *testing.B) {
	a, uuid := setupRepo(b)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			i := rand.Intn(b.N)
			acc := account.Account{
				UUID:     uuid[i],
				Password: "1234567890",
				Email:    "email" + string(uuid[i]),
			}
			if err := a.Update(&acc); err != nil {
				b.Error("err:", err.Error())
			}
		}
	})
}

func BenchmarkAccRepo_Base(b *testing.B) {
	a, _ := setupRepo(b)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			conn := a.db.Get()
			s := conn.Prepare("SELECT 1")
			if r, err := s.Step(); err != nil {
				b.Fatal(err)
			} else if r {
				_ = s.ColumnInt(0)
			} else {
				b.Fatal("No result")
			}
			s.Reset()
			a.db.Put(conn)
		}
	})
}
