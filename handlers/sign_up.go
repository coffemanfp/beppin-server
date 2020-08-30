package handlers

import (
	"errors"
	"fmt"
	"net/http"

	dbm "github.com/coffemanfp/beppin-server/database/models"
	errs "github.com/coffemanfp/beppin-server/errors"
	"github.com/coffemanfp/beppin-server/helpers"
	"github.com/coffemanfp/beppin-server/models"
	"github.com/labstack/echo"
)

// SignUp - Register a user.
func SignUp(c echo.Context) (err error) {
	var m models.ResponseMessage
	var user models.User

	if err = c.Bind(&user); err != nil {
		m.Error = errs.ErrInvalidBody

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	if !user.Validate("signup") {
		m.Error = errs.ErrInvalidBody

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	err = Storage.CreateUser(
		helpers.ShouldParseModelToDBModel(user).(dbm.User),
	)
	if err != nil {
		unwrappedErr := errors.Unwrap(err)

		switch unwrappedErr {
		case errs.ErrExistentObject:
			m.Error = fmt.Sprintf("%v: user", errs.ErrNotExistentObject)
			return echo.NewHTTPError(http.StatusConflict, m)

		case errs.ErrNotExistentObject:
			m.Error = fmt.Sprintf("%v: user", errs.ErrNotExistentObject)
			return echo.NewHTTPError(http.StatusNotFound, m)

		default:
			c.Logger().Error(err)
			m.Error = http.StatusText(http.StatusInternalServerError)

			return echo.NewHTTPError(http.StatusInternalServerError, m)
		}

	}

	m.Message = "Created."
	return c.JSON(http.StatusCreated, m)
}
