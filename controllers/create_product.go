package controllers

import (
	"net/http"

	"github.com/coffemanfp/beppin-server/database"
	dbm "github.com/coffemanfp/beppin-server/database/models"
	dbu "github.com/coffemanfp/beppin-server/database/utils"
	"github.com/coffemanfp/beppin-server/errors"
	"github.com/coffemanfp/beppin-server/helpers"
	"github.com/coffemanfp/beppin-server/models"
	"github.com/labstack/echo"
)

// CreateProduct - Creates a product.
func CreateProduct(c echo.Context) (err error) {
	var m models.ResponseMessage
	var product models.Product

	if err = c.Bind(&product); err != nil {
		m.Error = "invalid body"

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	if !product.Validate() {
		m.Error = "invalid body"

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	dbProductI, err := helpers.ParseModelToDBModel(product)
	if err != nil {
		c.Logger().Error(err)

		return echo.ErrInternalServerError
	}

	dbProduct := dbProductI.(dbm.Product)

	db, err := database.Get()
	if err != nil {
		c.Logger().Error(err)

		return echo.ErrInternalServerError
	}

	err = dbu.InsertProduct(db, dbProduct)
	if err != nil {
		if err.Error() == errors.ErrNotExistentObject {
			m.Error = err.Error() + " (user)"

			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)

		return echo.ErrInternalServerError
	}

	m.Message = "Created."

	return c.JSON(http.StatusCreated, m)
}