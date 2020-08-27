package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/coffemanfp/beppin-server/config"
	"github.com/coffemanfp/beppin-server/database"
	errs "github.com/coffemanfp/beppin-server/errors"
	"github.com/coffemanfp/beppin-server/handlers"
	"github.com/coffemanfp/beppin-server/models"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetProducts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name            string
		WithData        bool
		ExpectedContent interface{}
	}{
		{
			Name:            "without_products",
			ExpectedContent: models.Products{},
		},
		{
			Name:     "with_products",
			WithData: true,
			ExpectedContent: models.Products{
				models.Product{
					UserID:      exampleProducts[0].UserID,
					Name:        fmt.Sprintf("Product at %d ", time.Now().Unix()),
					Description: exampleProducts[0].Description,
					Categories:  exampleProducts[0].Categories,
				},
			},
		},
	}

	for _, ts := range tests {
		ts := ts

		t.Run(ts.Name, func(t *testing.T) {
			t.Parallel()

			// Setup server
			e := echo.New()
			e.Logger.Debug()

			setStorage(t)

			if ts.WithData {

				if !existsLanguage(t, exampleLanguage) {
					insertLanguage(t, exampleLanguage)
				}
				if !existsUser(t, exampleUser) {
					insertUser(t, exampleUser)
				}
				insertProduct(t, ts.ExpectedContent.(models.Products)[0])
			}

			e.GET("/", handlers.GetProducts)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			var m models.ResponseMessage
			m = decodeResponseMessage(t, rec)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, models.TypeProducts, m.ContentType)

			if ts.WithData {
				var exists bool

				for _, productContent := range m.Content.([]interface{}) {
					if productContent.(map[string]interface{})["name"] == ts.ExpectedContent.(models.Products)[0].Name {
						exists = true
					}
				}
				assert.True(t, exists)
			}
		})
	}
}

func TestFailedGetProducts(t *testing.T) {
	t.Parallel()
	invalidParamLimit := fmt.Sprintf("%v: limit", errs.ErrInvalidParam)
	invalidParamOffset := fmt.Sprintf("%v: offset", errs.ErrInvalidParam)

	tests := []struct {
		Name               string
		QueryParams        url.Values
		ExpectedStatusCode int
		ExpectedError      string
		WithDatabase       bool
	}{
		{
			Name: "limit_negative_number",
			QueryParams: url.Values{
				"limit": []string{
					"-1",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamLimit,
		},
		{
			Name: "limit_super_negative_number",
			QueryParams: url.Values{
				"limit": []string{
					"-986544567890",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamLimit,
		},
		{
			Name: "limit_letters",
			QueryParams: url.Values{
				"limit": []string{
					"a",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamLimit,
		},
		{
			Name: "limit_super_letters",
			QueryParams: url.Values{
				"limit": []string{
					"ajhkklaskldjkasksjdlfkjsdlfkjlasdkjfljsdf",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamLimit,
		},
		{
			Name: "limit_super_greater_max",
			QueryParams: url.Values{
				"limit": []string{
					strconv.Itoa(int(config.GetSettings().MaxElementsPerPagination)) + "09876545678909876545678987678",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamLimit,
		},
		{
			Name: "offset_negative_number",
			QueryParams: url.Values{
				"offset": []string{
					"-1",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamOffset,
		},
		{
			Name: "offset_super_negative_number",
			QueryParams: url.Values{
				"offset": []string{
					"-986544567890",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamOffset,
		},
		{
			Name: "offset_letters",
			QueryParams: url.Values{
				"offset": []string{
					"a",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamOffset,
		},
		{
			Name: "offset_super_letters",
			QueryParams: url.Values{
				"offset": []string{
					"ajhkklaskldjkasksjdlfkjsdlfkjlasdkjfljsdf",
				},
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedError:      invalidParamOffset,
		},
		{
			Name:               "without_database",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedError:      http.StatusText(http.StatusInternalServerError),
		},
	}

	for _, ts := range tests {
		ts := ts
		t.Run(ts.Name, func(t *testing.T) {
			t.Parallel()
			// Setup server
			e := echo.New()
			e.Logger.Debug()

			e.GET("/", handlers.GetProducts)

			if ts.WithDatabase {
				var storage database.Storage
				storage, err := database.NewDefault()

				assert.Nil(t, err)

				handlers.Storage = storage
			} else {
				handlers.Storage = database.New(nil)
			}

			// Now the request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.URL.RawQuery = ts.QueryParams.Encode()

			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, ts.ExpectedStatusCode, rec.Code)
			assertResponseError(t, ts.ExpectedError, decodeResponseMessage(t, rec))
		})
	}
}
