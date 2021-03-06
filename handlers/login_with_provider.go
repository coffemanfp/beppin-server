package handlers

import (
	"fmt"
	"net/http"

	errs "github.com/coffemanfp/beppin/errors"
	"github.com/coffemanfp/beppin/models"
	"github.com/labstack/echo"
	"github.com/stretchr/gomniauth"
)

// LoginWithProvider login the user with a provider.
func LoginWithProvider(c echo.Context) (err error) {
	var m models.ResponseMessage

	providerParam := c.Param("provider")
	if providerParam == "" {
		m.Error = fmt.Sprintf("%v: provider", errs.ErrInvalidParam)

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	provider, err := gomniauth.Provider(providerParam)
	if err != nil {
		err = fmt.Errorf("failed to get (%s) provider: %v", providerParam, err)
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	loginURL, err := provider.GetBeginAuthURL(nil, nil)
	if err != nil {
		err = fmt.Errorf("failed to get begin auth url for (%s) provider: %v", providerParam, err)
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	c.Redirect(http.StatusTemporaryRedirect, loginURL)
	return
}
