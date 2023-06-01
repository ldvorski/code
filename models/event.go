package models

import (
	"database/sql"
)

type Event struct {
	ID       int
	Name     string
	Date     Date
	Duration Time
}

// Create a new event
func createEvent(db *sql.DB, name string, date Date, duration Time)	(Event, error) {
	result, err := db.Exec("INSERT INTO Event (name, date, duration) VALUES (?, ?, ?)", name, date, duration)
	if err != nil {
		return Event{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Event{}, err
	}

	return Event{
		ID: int(id), Name: name, Date: date, Duration: duration,
	},nil
}

// Allow a user to join an event
func joinEvent(db *sql.DB, userID int, eventID int) error {
	_, err := db.Exec("INSERT INTO Participation (user_id, event_id) VALUES (?, ?)", userID, eventID)
	if err != nil {
		return err
	}
	return nil
}

// Allow a user to leave an event
func leaveEvent(db *sql.DB, userID int, eventID int) error {
	_, err := db.Exec("DELETE FROM Participation WHERE user_id = ? AND event_id = ?", userID, eventID)
	if err != nil {
		return err
	}
	return nil
}
