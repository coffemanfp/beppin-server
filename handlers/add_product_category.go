package handlers

import (
	"errors"
	"fmt"
	"net/http"

	errs "github.com/coffemanfp/beppin/errors"
	"github.com/coffemanfp/beppin/models"
	"github.com/coffemanfp/beppin/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// AddProductCategory - Adds a product category relation.
func AddProductCategory(c echo.Context) (err error) {
	var m models.ResponseMessage
	var productID, categoryID int

	userIDToken := c.Get("user").(*jwt.Token).Claims.(*models.Claim).User.ID
	productIDParam := c.Param("productid")
	categoryIDParam := c.Param("categoryid")

	if productID, err = utils.Atoi(productIDParam); err != nil || productID == 0 {
		m.Error = fmt.Sprintf("%v: id", errs.ErrInvalidParam)

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	if categoryID, err = utils.Atoi(categoryIDParam); err != nil || categoryID == 0 {
		m.Error = fmt.Sprintf("%v: id", errs.ErrInvalidParam)

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	// Get old product info
	oldProduct, err := Storage.GetProduct(models.Product{ID: int64(productID)})
	if err != nil {
		if errors.Is(err, errs.ErrNotExistentObject) {
			m.Error = fmt.Sprintf("%v: product", errs.ErrNotExistentObject)

			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	if oldProduct.UserID != userIDToken {
		m.Error = http.StatusText(http.StatusUnauthorized)

		return echo.NewHTTPError(http.StatusUnauthorized, m)
	}

	err = Storage.AddProductCategory(int64(productID), int64(categoryID))
	if err != nil {
		if errors.Is(err, errs.ErrNotExistentObject) ||
			errors.Is(err, errs.ErrExistentObject) {
			m.Error = err.Error()

			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	m.Message = "Created."
	return c.JSON(http.StatusCreated, m)
}
