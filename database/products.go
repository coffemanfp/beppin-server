package database

import (
	"fmt"

	"github.com/coffemanfp/beppin-server/database/models"
	dbu "github.com/coffemanfp/beppin-server/database/utils"
	errs "github.com/coffemanfp/beppin-server/errors"
)

// ProductStorage reprensents all implementations for product utils.
// type ProductStorage interface {
// 	CreateProduct(product models.Product) error
// 	GetProduct(productToFind models.Product) (models.Product, error)
// 	GetProducts(limit, offset int) (models.Products, error)
// 	UpdateProduct(productToUpdate, product models.Product) error
// 	DeleteProduct(product models.Product) error
// }

func (dS defaultStorage) CreateProduct(product models.Product) (err error) {
	exists, err := dS.ExistsUser(models.User{ID: product.UserID})
	if err != nil {
		return
	}

	if !exists {
		err = fmt.Errorf("failed to check (%d) user: %w", product.UserID, errs.ErrNotExistentObject)
		return
	}

	err = dbu.InsertProduct(dS.db, product)
	return
}

func (dS defaultStorage) GetProduct(productToFind models.Product) (product models.Product, err error) {
	product, err = dbu.SelectProduct(dS.db, productToFind)
	return
}

func (dS defaultStorage) GetProducts(limit, offset int) (products models.Products, err error) {
	products, err = dbu.SelectProducts(dS.db, limit, offset)
	return
}

func (dS defaultStorage) UpdateProduct(productToUpdate, product models.Product) (err error) {
	productToUpdate, err = dS.GetProduct(
		models.Product{
			ID: productToUpdate.ID,
		},
	)
	if err != nil {
		return
	}

	product = fillProductEmptyFields(product, productToUpdate)

	err = dbu.UpdateProduct(dS.db, productToUpdate, product)
	return
}

func (dS defaultStorage) DeleteProduct(productToDelete models.Product) (err error) {
	err = dbu.DeleteProduct(dS.db, productToDelete)
	return
}

func fillProductEmptyFields(product, previousProductData models.Product) models.Product {

	switch "" {
	case product.Name:
		product.Name = previousProductData.Name
	case product.Description:
		product.Description = previousProductData.Description
	}

	if product.Categories == nil {
		product.Categories = previousProductData.Categories
	}

	return product
}