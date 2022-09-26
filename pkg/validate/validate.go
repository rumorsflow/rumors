package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
	"strings"
)

var _ echo.Validator = (*customValidator)(nil)

type customValidator struct {
	validator *validator.Validate
}

func New() *customValidator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			for _, name = range []string{"form", "query", "param", "header"} {
				if name = fld.Tag.Get(name); len(name) > 0 {
					return name
				}
			}
			return ""
		}
		return name
	})

	return &customValidator{validator: v}
}

func (cv *customValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	return nil
}
