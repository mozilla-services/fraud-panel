package db

import (
	"database/sql"
	"fmt"
)

func UserLogIn(name, email string) (sessionID string, err error) {
	var userID uint64
	err = h.QueryRow(`SELECT id FROM member WHERE email = $1`, email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = h.QueryRow(`INSERT INTO member(name, email) VALUES ($1, $2)
						ON CONFLICT ON CONSTRAINT member_email_key DO NOTHING
						RETURNING id `,
				name, email).Scan(&userID)
			if err != nil {
				return
			}
		} else {
			return
		}
	}
	err = h.QueryRow(`INSERT INTO session(member_id) VALUES ($1) RETURNING id`,
		userID).Scan(&sessionID)
	if err != nil {
		return
	}
	if sessionID == "" {
		err = fmt.Errorf("failed to generate a valid session ID")
	}
	return
}
