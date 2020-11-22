package handlers

import (
	"errors"
	"fmt"
	"net/http"

	errs "github.com/coffemanfp/beppin/errors"
	"github.com/coffemanfp/beppin/models"
	"github.com/coffemanfp/beppin/utils"
	"github.com/labstack/echo"
)

// UpdateProduct - Updates a product.
func UpdateProduct(c echo.Context) (err error) {
	productIDParam := c.Param("id")
	var m models.ResponseMessage

	productID, err := utils.Atoi(productIDParam)
	if err != nil || productID == 0 {
		m.Error = fmt.Sprintf("%v: id", errs.ErrInvalidParam)

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	var product models.Product

	if err = c.Bind(&product); err != nil {
		m.Error = errs.ErrInvalidBody

		return echo.NewHTTPError(http.StatusBadRequest, m)
	}

	// Check if exists the files
	if len(product.Images) != 0 {
		for _, image := range product.Images {
			// Get old file info
			oldFile, err := Storage.GetFile(models.File{ID: int64(image.ID)})
			if err != nil {
				if errors.Is(err, errs.ErrNotExistentObject) {
					m.Error = fmt.Sprintf("%v: file", errs.ErrNotExistentObject)

					return echo.NewHTTPError(http.StatusNotFound, m)
				}
				c.Logger().Error(err)
				m.Error = http.StatusText(http.StatusInternalServerError)

				return echo.NewHTTPError(http.StatusInternalServerError, m)
			}

			// Check if the file exists in the filesystem
			exists, err := utils.ExistsFile(oldFile.Path)
			if err != nil {
				c.Logger().Error(err)
				m.Error = http.StatusText(http.StatusInternalServerError)

				return echo.NewHTTPError(http.StatusInternalServerError, m)
			}

			if !exists {
				m.Error = fmt.Sprintf("%v: file", errs.ErrNotExistentObject)

				return echo.NewHTTPError(http.StatusNotFound, m)
			}

		}
	}

	updatedProduct, err := Storage.UpdateProduct(
		models.Product{
			ID: int64(productID),
		},
		product,
	)
	if err != nil {
		if errors.Is(err, errs.ErrNotExistentObject) {
			m.Error = fmt.Sprintf("%v: product", errs.ErrExistentObject)

			return echo.NewHTTPError(http.StatusNotFound, m)
		}
		c.Logger().Error(err)
		m.Error = http.StatusText(http.StatusInternalServerError)

		return echo.NewHTTPError(http.StatusInternalServerError, m)
	}

	m.Message = "Updated."
	m.Content = updatedProduct
	m.ContentType = models.TypeProduct
	return c.JSON(http.StatusOK, m)
}
