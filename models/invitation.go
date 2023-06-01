package models

import (
	"database/sql"
)

type Invitation struct {
	User_ID    int
	Meeting_ID int
	IsAccepted bool
	Responded  bool
}

func CreateInvitations(db *sql.DB, meetingID int, users []User) error {
    for _, user := range users {
        err := CreateInvitation(db, meetingID, user.ID)
        if err != nil {
            return err
        }
    }
    return nil
}

func CreateInvitation(db *sql.DB, meetingID int, userID int) error {
	_, err := db.Exec("INSERT INTO Invitation (user_id, meeting_id, is_accepted, response) VALUES (?, ?, ?, ?)",
	 userID, meetingID, false, false)
	if err != nil {
		return err
	}
	return nil
}

// Update the status of an invitation
func UpdateInvitationStatus(db *sql.DB, userID int, meetingID int, is_accepted bool) error {
	_, err := db.Exec("UPDATE Invitation SET is_accepted = ?, responded_to = true WHERE user_id = ? AND meeting_id = ?", 
	is_accepted, userID, meetingID)
	if err != nil {
		return err
	}
	return nil
}
