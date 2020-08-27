package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	dbm "github.com/coffemanfp/beppin-server/database/models"
	errs "github.com/coffemanfp/beppin-server/errors"
	"github.com/coffemanfp/beppin-server/models"
	"github.com/coffemanfp/beppin-server/utils"
	"github.com/labstack/echo"
)

// UpdateAvatar Updates the user avatar.
func UpdateAvatar(c echo.Context) (err error) {
	var avatar models.Avatar
	var m models.ResponseMessage
	var userID uint64

	userIDParam := c.Param("id")

	if userID, err = utils.ParseUint(userIDParam, 64); err != nil || userID == 0 {
		m.Error = fmt.Sprintf("%v: id", errs.ErrInvalidParam)

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	if err = c.Bind(&avatar); err != nil {
		m.Error = errs.ErrInvalidBody

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	if !avatar.Validate() {
		m.Error = errs.ErrInvalidBody

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	avatarURL := avatar.URL

	// If the user will save their avatar in our file system
	if avatar.Data != "" && avatar.URL == "" {
		avatarURL, err = avatar.Save(strconv.Itoa(int(userID)))
		if err != nil {
			c.Logger().Error(err)
			m.Error = http.StatusText(http.StatusInternalServerError)

			return echo.NewHTTPError(http.StatusInternalServerError, m)
		}
	}

	err = Storage.UpdateAvatar(
		avatarURL,
		dbm.User{ID: int64(userID)},
	)
	if err != nil {
		if errors.Is(err, errs.ErrNotExistentObject) {
			m.Error = fmt.Sprintf("%v: user", errs.ErrNotExistentObject)

			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}
	if err != nil {
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	m.Message = "Updated."

	return c.JSON(http.StatusOK, m)
}
