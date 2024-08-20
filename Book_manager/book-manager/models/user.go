package models

import (
	"book-manager/db"
	"book-manager/utils"
	"errors"
	"log"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (email, password) 
	VALUES ($1, $2) RETURNING id`
	stmt := "saveQuery"
	conn, err := db.DB.Acquire(db.Ctx)
	defer conn.Release()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Conn().Prepare(db.Ctx, stmt, query)
	if err != nil {
		log.Println(err)
		log.Println("111")
		return err
	}

	hashedPassword, err := utils.HashPassword(u.Password)

	if err != nil {
		log.Println(err)
		log.Println("222")
		return err
	}
	var lastInsertID int64
	err = conn.QueryRow(db.Ctx, stmt, u.Email, hashedPassword).Scan(&lastInsertID)

	if err != nil {
		log.Println(err)
		log.Println("333")
		return err
	}
	u.ID = lastInsertID
	return err
}

func (u *User) Authenticate() error {
	query := "SELECT id, password FROM users WHERE email = $1"
	row := db.DB.QueryRow(db.Ctx, query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return errors.New("invalid credentials")
	}

	passIsValid := utils.CheckPasswordHash(u.Password, retrievedPassword)

	if !passIsValid {
		return errors.New("invalid credentials")
	}
	return nil
}
