package model

import (
	"crypto/rand"
	"database/sql"
	"math/big"

	"pehks1980/shard/mypool"
	"pehks1980/shard/myshard"
	"time"
)

type Con struct {
	M *myshard.Manager
	P *mypool.Pool
}

type User struct {
	UserId int
	Name   string
	Age    int
	Spouse int
	Con    Con
}

func (u *User) connection() (*sql.DB, error) {
	s, err := u.Con.M.ShardById(u.UserId)
	if err != nil {
		return nil, err
	}
	sqldb, err1 := u.Con.P.Connection(s.Address)
	if err1 != nil {
		return nil, err1
	}
	return sqldb, nil
}

func (u *User) Create() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO "users" VALUES ($1, $2, $3, $4)`, u.UserId, u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Read() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	r := c.QueryRow(`SELECT "name", "age", "spouse" FROM "users" WHERE "user_id" = $1`, u.UserId)
	return r.Scan(
		&u.Name,
		&u.Age,
		&u.Spouse,
	)
}

func (u *User) Update() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "users" SET "name" = $2, "age" = $3, "spouse" = $4 WHERE "user_id" = $1`, u.UserId,
		u.Name, u.Age, u.Spouse)
	return err
}

func (u *User) Delete() error {
	c, err := u.connection()
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "users" WHERE "user_id" = $1`, u.UserId)
	return err
}

type Activity struct {
	UserId int
	Name   string
	Date   time.Time
	Con    Con
}

func (a *Activity) connection(mode int) (*sql.DB, error) {
	s, err := a.Con.M.ShardById(a.UserId)
	if err != nil {
		return nil, err
	}
	var sqldb *sql.DB
	// task 2 get data from replica in case of read operation
	if mode == 1 {
		sqldb, err = a.Con.P.Connection(s.AddressRepl)
	} else {
		sqldb, err = a.Con.P.Connection(s.Address)
	}

	// task 3 probablity based choose connection for reading P = 0.5 for master or replica
	if mode == 1 {
		n, _ := rand.Int(rand.Reader, big.NewInt(100))
		if n.Int64() < 50 {
			sqldb, err = a.Con.P.Connection(s.AddressRepl)
		} else {
			sqldb, err = a.Con.P.Connection(s.Address)
		}
	}

	if err != nil {
		return nil, err
	}
	return sqldb, nil
}

func (a *Activity) Create() error {
	c, err := a.connection(0)
	if err != nil {
		return err
	}
	_, err = c.Exec(`INSERT INTO activities VALUES ($1, $2, $3::timestamp)`, a.UserId, a.Name, a.Date)
	return err
}

func (a *Activity) Read(userid int, con Con) (Activity, error) {
	a.Con = con
	a.UserId = userid
	mode := 1 // means 1 - read from relica 0 - read from master
	c, err := a.connection(mode)
	if err != nil {
		return Activity{}, err
	}
	r := c.QueryRow(`SELECT "name", "date"::timestamp FROM "activities" WHERE "user_id" = $1`, userid)
	var act Activity
	_ = r.Scan(&act.Name, &act.Date)
	act.UserId = userid
	return act, nil
}

func (a *Activity) Update() error {
	c, err := a.connection(0)
	if err != nil {
		return err
	}
	_, err = c.Exec(`UPDATE "activities" SET "name" = $2, "age" = $3 WHERE "user_id" = $1`, a.UserId,
		a.Name, a.Date)
	return err
}

func (a *Activity) Delete(userid int, con Con) error {
	a.UserId = userid
	a.Con = con
	c, err := a.connection(0)
	if err != nil {
		return err
	}
	_, err = c.Exec(`DELETE FROM "activities" WHERE "user_id" = $1`, userid)
	return err
}
