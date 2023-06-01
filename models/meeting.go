package models

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Meeting struct {
	ID          int
	Date        Date
	Time        Time
	IsScheduled bool
}

func NewMeeting(db *sql.DB, dateStr string, timeStr string) (Meeting, error) {
	var meeting Meeting
	date, err := NewDate(dateStr)
	if err != nil {
		return meeting, err
	}
	time, err := NewTime(timeStr)
	if err != nil {
		return meeting, err
	}
	meeting = Meeting{
		Date:        date,
		Time:        time,
		IsScheduled: false,
	}

	err = createMeeting(db, meeting)
	if err != nil {
		return meeting, err
	}

	return meeting, nil
}

func createMeeting(db *sql.DB, meeting Meeting) error {
	// Convert the Date and Time fields to the format expected by MySQL
	dateStr := fmt.Sprintf("%d-%d-%d", meeting.Date.Year, meeting.Date.Month, meeting.Date.Day)
	timeStr := fmt.Sprintf("%d:%d:00", meeting.Time.Hour, meeting.Time.Minute)

	// Insert the meeting into the database
	_, err := db.Exec("INSERT INTO Meeting (date, time, is_scheduled) VALUES (?, ?, ?)", dateStr, timeStr, meeting.IsScheduled)
	if err != nil {
		return err
	}
	return nil
}

func GetMeetingsForUser(db *sql.DB, userID int) ([]Meeting, error) {
	// Query the database for meetings that include the given user
	rows, err := db.Query("SELECT m.meeting_id, m.date, m.time, m.is_scheduled FROM Meeting m JOIN Invitation i ON m.meeting_id = i.meeting_id WHERE i.user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the results and create a slice of Meeting structs
	var meetings []Meeting
	for rows.Next() {
		var m Meeting
		var dateStr string
		var timeStr string
		if err := rows.Scan(&m.ID, &dateStr, &timeStr, &m.IsScheduled); err != nil {
			return nil, err
		}

		// Convert the date and time strings to the Date and Time structs
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

func areAllInvitationsAccepted(db *sql.DB, meetingID int) (bool, error) {
	// Query the database for invitations for the given meeting
	rows, err := db.Query("SELECT is_accepted, responded_to FROM Invitation WHERE meeting_id = ?", meetingID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Check if all invitations have been responded to and accepted
	for rows.Next() {
		var isAccepted bool
		var respondedTo bool
		if err := rows.Scan(&isAccepted, &respondedTo); err != nil {
			return false, err
		}
		if !respondedTo || !isAccepted {
			return false, nil
		}
	}

	return true, nil
}

func ScheduleMeetingIfAllAccepted(db *sql.DB, meeting *Meeting) error {
	// Check if all invitations have been responded to and accepted
	allAccepted, err := areAllInvitationsAccepted(db, meeting.ID)
	if err != nil {
		return err
	}

	// Update the is_scheduled field in the database and the IsScheduled field in the Meeting struct if all invitations have been accepted
	if allAccepted {
		_, err := db.Exec("UPDATE Meeting SET is_scheduled = true WHERE meeting_id = ?", meeting.ID)
		if err != nil {
			return err
		}
		meeting.IsScheduled = true
	}

	return nil
}
