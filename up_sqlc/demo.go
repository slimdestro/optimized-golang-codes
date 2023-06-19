/**
  // Sample code for : squirrel-golang
  // @slimdestro
*/

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/Masterminds/squirrel"
)

type User struct {
	ID                int
	FirstName         string
	LastName          string
	EmailID           string
	PhoneNumber       string
	HashedPass        string
	Country           string
	City              string
	Zip               int
	StatusDescription string
}

func main() {
	db, err := sql.Open("mysql", "(ommitted):(ommitted)@tcp(localhost:3306)/destro.svc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background() 
	qb := squirrel.Select("*").From("users"). 
		Where(squirrel.Eq{"email_id": "mukul@gmail.com"})
		// emailID mukul@gmail.com doesnt exist

	query, args, err := qb.ToSql()
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.EmailID,
			&user.PhoneNumber,
			&user.HashedPass,
			&user.Country,
			&user.City,
			&user.Zip, 
			&user.StatusDescription,
		)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		fmt.Println(user)
	}
}