package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID              int
	Name            string
	Organization_ID int
}

func CreateUser(db *sql.DB, name string, org Organization) (User, error) {
	result, err := db.Exec("INSERT INTO User (name, organization_id) VALUES (?, ?)",
		name, org.ID)

	if err != nil {
		return User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}
	return User{
		ID: int(id), Name: name, Organization_ID: org.ID,
	}, nil
}

// Get all scheduled meetings for a user
func getScheduledMeetingsForUser(db *sql.DB, userID int) ([]Meeting, error) {
	rows, err := db.Query(`
        SELECT m.meeting_id, m.date, m.time FROM Meeting m 
        JOIN Invitation i ON m.meeting_id = i.meeting_id 
        WHERE i.user_id = ? AND m.is_scheduled = true`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []Meeting
	for rows.Next() {
		var m Meeting
		var dateStr string
		var timeStr string
		if err := rows.Scan(&m.ID, &dateStr, &timeStr); err != nil {
			return nil, err
		}

		m.Date, err = NewDate(dateStr)
		if err != nil {
			return nil, err
		}

		m.Time, err = NewTime(timeStr)
		if err != nil {
			return nil, err
		}

		meetings = append(meetings, m)
	}

	return meetings, nil
}

// Get all invitations for a specific user
func getInvitationsForUser(db *sql.DB, userID int) ([]Invitation, error) {
    rows, err := db.Query("SELECT user_id, meeting_id, is_accepted, responded_to FROM Invitation WHERE user_id = ?", userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var invitations []Invitation
    for rows.Next() {
        var i Invitation
        if err := rows.Scan(&i.User_ID, &i.Meeting_ID, &i.IsAccepted, &i.Responded); err != nil {
            return nil, err
        }
        invitations = append(invitations, i)
    }

    return invitations, nil
}

func AcceptInvitation(db *sql.DB, userID int, meetingID int) error {
	err := UpdateInvitationStatus(db, userID, meetingID, true)
	return err
}

func DeclineInvitation(db *sql.DB, userID int, meetingID int) error {
	err := UpdateInvitationStatus(db, userID, meetingID, false)
	return err
}