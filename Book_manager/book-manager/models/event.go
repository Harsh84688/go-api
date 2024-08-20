package models

import (
	"book-manager/cache"
	"book-manager/db"
	"encoding/json"
	"log"
	"time"
)

type Event struct {
	ID          int64
	Name        string `binding:"required"`
	Author      string `binding:"required"`
	Description string `binding:"required"`
	UserId      int64
}

var events = []Event{}

func (e *Event) Save() error {
	query := `
	INSERT INTO events (name,author, description,user_id) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	stmt := "insertQuery"
	conn, err := db.DB.Acquire(db.Ctx)
	defer conn.Release()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Conn().Prepare(db.Ctx, stmt, query)
	if err != nil {
		log.Println(err)
		return err
	}
	var lastInsertID int64
	err = conn.QueryRow(db.Ctx, stmt, e.Name, e.Author, e.Description, e.UserId).Scan(&lastInsertID)

	if err != nil {
		log.Println(err)
		log.Println("111")
		return err
	}

	e.ID = lastInsertID
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(db.Ctx, query)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Author, &event.Description, &event.UserId)

		if err != nil {
			log.Println(err)
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func cacheOutOp(idc string, event *Event) (*Event, error) {
	data, err := cache.GetCachedResult(idc)
	if err != nil {
		return nil, err
	}
	if data != "" {
		log.Println(data)
		err = json.Unmarshal([]byte(data), &event)

		if err != nil {
			return nil, err
		}
		if event.ID == 0 {
			return nil, nil
		}
		return event, nil
	}
	return nil, nil
}

func cacheInOp(idc string, event Event, ctime time.Duration) {
	datac, err := json.Marshal(event)
	if err != nil {
		panic("error marshalling event")
	}

	err = cache.SetCachedResult(idc, string(datac), ctime)

	if err != nil {
		panic("error caching event")
	}

}

func GetEventByID(id int64, idc string) (*Event, error) {
	var event Event

	data, err := cacheOutOp(idc, &event)

	if err != nil {
		return nil, err
	}

	if data != nil {
		return data, nil
	}
	query := "SELECT * FROM events WHERE id = $1"
	row := db.DB.QueryRow(db.Ctx, query, id)

	err = row.Scan(&event.ID, &event.Name, &event.Author, &event.Description, &event.UserId)

	go cacheInOp(idc, event, 10*time.Minute)

	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (e Event) Update(idc string) error {

	go cache.DeleteCachedResult(idc)
	query := `
	UPDATE events
	SET name = $1, author = $2   ,  description = $3,  user_id = $4
	WHERE id = $5`

	stmt := "updateQuery"
	conn, err := db.DB.Acquire(db.Ctx)
	defer conn.Release()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Conn().Prepare(db.Ctx, stmt, query)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Exec(db.Ctx, stmt, e.Name, e.Author, e.Description, e.UserId, e.ID)

	return err
}

func (e Event) Delete(idc string) error {
	go cache.DeleteCachedResult(idc)
	query := "DELETE FROM events WHERE id = $1"
	stmt := "deleteQuery"
	conn, err := db.DB.Acquire(db.Ctx)
	defer conn.Release()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Conn().Prepare(db.Ctx, stmt, query)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = conn.Exec(db.Ctx, stmt, e.ID)
	return err
}
