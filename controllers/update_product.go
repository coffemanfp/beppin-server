package controllers

import (
	"net/http"

	"github.com/coffemanfp/beppin-server/database"
	dbm "github.com/coffemanfp/beppin-server/database/models"
	dbu "github.com/coffemanfp/beppin-server/database/utils"
	"github.com/coffemanfp/beppin-server/errors"
	"github.com/coffemanfp/beppin-server/helpers"
	"github.com/coffemanfp/beppin-server/models"
	"github.com/coffemanfp/beppin-server/utils"
	"github.com/labstack/echo"
)

// UpdateProduct - Updates a product.
func UpdateProduct(c echo.Context) (err error) {
	productIDParam := c.Param("id")
	var m models.ResponseMessage

	productID, err := utils.Atoi(productIDParam)
	if err != nil || productID == 0 {
		m.Error = "id param not valid"

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	var product models.Product

	if err = c.Bind(&product); err != nil {
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

	err = dbu.UpdateProduct(db, productID, dbProduct)
	if err != nil {
		if err.Error() == errors.ErrNotExistentObject {
			m.Error = err.Error() + " (product)"
			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)

		return echo.ErrInternalServerError
	}

	m.Message = "Updated."

	return c.JSON(http.StatusOK, m)
}