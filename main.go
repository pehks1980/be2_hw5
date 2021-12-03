package main

import (
	"fmt"
	"pehks1980/shard/model"
	"pehks1980/shard/mypool"
	"pehks1980/shard/myshard"
	"time"
)

func main() {

	con := model.Con{
		M: myshard.NewManager(2),
		P: mypool.NewPool(),
	}
	con.M.Add(&myshard.Shard{"port=8181 user=test password=test dbname=test sslmode=disable",
		"port=8182 user=test password=test dbname=test sslmode=disable",
		0})
	con.M.Add(&myshard.Shard{"port=8191 user=test password=test dbname=test sslmode=disable",
		"port=8192 user=test password=test dbname=test sslmode=disable", 1})

	users := []*model.User{
		{1, "Joe Biden", 78, 10, con},
		{10, "Jill Biden", 69, 1, con},
		{13, "Donald Trump", 74, 25, con},
		{24, "Melania Trump", 78, 13, con},
	}
	for _, user := range users {
		err := user.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create user %v: %w", user, err))
		}
	}

	activities := []*model.Activity{
		{1, "Eating", time.Now(), con},
		{10, "Hiking", time.Now(), con},
		{13, "Playing Black Jack", time.Now(), con},
		{24, "Shopping", time.Now(), con},
	}
	for _, act := range activities {
		err := act.Create()
		if err != nil {
			fmt.Println(fmt.Errorf("error on create activity %v: %w", act, err))
		}
	}

	actUserId := []int{1, 10, 13, 24}
	for _, id := range actUserId {
		var act model.Activity
		act, err := act.Read(id, con)
		if err != nil {
			fmt.Println(fmt.Errorf("error on create activity %v: %w", act, err))
		}
		fmt.Println(fmt.Errorf("activity %v", act))

	}

	var act model.Activity
	err := act.Delete(24, con)
	if err != nil {
		fmt.Println(fmt.Errorf("error on delete activity %v: %w", act, err))
	}
	/*
		users := []*model.User{
			{1, "Joe Biden", 788, 10, m,p},
			{10, "Jill Biden", 698, 1,m,p},
			{13, "Donald Trump", 748, 25, m, p},
			{24, "Melania Trump", 788, 13, m, p},
		}
		for _, user := range users {
			err := user.Update()
			if err != nil {
				fmt.Println(fmt.Errorf("error on update user %v: %w", user, err))
			}
		}
	*/
}
