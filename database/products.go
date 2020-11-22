package database

import (
	"fmt"

	dbu "github.com/coffemanfp/beppin/database/utils"
	errs "github.com/coffemanfp/beppin/errors"
	"github.com/coffemanfp/beppin/models"
)

func (dS defaultStorage) CreateProduct(product models.Product) (createdProduct models.Product, err error) {
	exists, err := dS.ExistsUser(models.User{ID: product.UserID})
	if err != nil {
		return
	}

	if !exists {
		err = fmt.Errorf("failed to check (%d) user: %w", product.UserID, errs.ErrNotExistentObject)
		return
	}

	if product.Images == nil && len(product.Images) == 0 {
		createdProduct, err = dbu.InsertProduct(dS.db, product)
		return
	}

	createdProduct, err = dbu.InsertProduct(dS.db, product)
	if err != nil {
		return
	}
	for _, file := range product.Images {
		exists, err = dbu.ExistsFile(dS.db, models.File{ID: file.ID})
		if err != nil {
			return
		}

		if !exists {
			err = fmt.Errorf("failed to check (%d) file: %w", file.ID, errs.ErrNotExistentObject)
			return
		}

		err = dbu.InsertProductFile(dS.db, createdProduct.ID, file.ID)
		if err != nil {
			return
		}
	}

	return
}

func (dS defaultStorage) GetProduct(productToFind models.Product) (product models.Product, err error) {
	product, err = dbu.SelectProduct(dS.db, productToFind)
	if err != nil {
		return
	}

	files, err := dbu.SelectProductFiles(dS.db, productToFind)
	if err != nil {
		return
	}

	product.Images = files
	return
}

func (dS defaultStorage) GetProducts(limit, offset int) (products models.Products, err error) {
	products, err = dbu.SelectProducts(dS.db, limit, offset)

	var files models.Files
	for i := 0; i < len(products); i++ {
		files, err = dbu.SelectProductFiles(dS.db, products[i])
		if err != nil {
			return
		}

		products[i].Images = files
	}
	return
}

func (dS defaultStorage) UpdateProduct(productToUpdate, product models.Product) (updatedProduct models.Product, err error) {
	updatedProduct, err = dbu.UpdateProduct(dS.db, productToUpdate, product)
	return
}

func (dS defaultStorage) DeleteProduct(productToDelete models.Product) (id int, err error) {
	id, err = dbu.DeleteProduct(dS.db, productToDelete)
	return
}
