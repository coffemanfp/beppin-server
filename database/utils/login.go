package utils

import (
	"database/sql"
	"errors"
	"fmt"

	errs "github.com/coffemanfp/beppin/errors"
	"github.com/coffemanfp/beppin/models"
)

// Login - Select a user by his username and password, and checks if exists.
func Login(dbtx DBTX, userToLogin models.User) (user models.User, match bool, err error) {
	if dbtx == nil {
		err = errs.ErrClosedDatabase
		return
	}

	match = true

	query := `
		SELECT
			users.id, files.id, files.path, language, username, theme, currency
		FROM
			users
		LEFT JOIN
			files
		ON
			users.avatar_id = files.id
		WHERE
			username = $1 AND password = $2 OR email = $3 AND password = $2
	`

	stmt, err := dbtx.Prepare(query)
	if err != nil {
		err = fmt.Errorf("failed to prepare the login (%s) user statement: %v", userToLogin.Username, err)
		return
	}
	defer stmt.Close()

	var nullData nullUserData

	err = stmt.QueryRow(
		userToLogin.Username,
		userToLogin.Password,
		userToLogin.Email,
	).Scan(
		&user.ID,
		&nullData.AvatarID,
		&nullData.AvatarPath,
		&user.Language,
		&user.Username,
		&user.Theme,
		&user.Currency,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			match = false
			return
		}

		err = fmt.Errorf("failed to login (%s) user: %v", userToLogin.Username, err)
	}

	nullData.setResults(&user)
	if user.Avatar != nil {
		user.Avatar.SetURL()
	}
	return
}
