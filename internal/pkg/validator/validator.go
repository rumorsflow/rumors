package validator

import (
	"context"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

var v *validator.Validate

func init() {
	v = validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" || name == "" {
			return ""
		}
		return name
	})
}

func Validate(i any) error {
	return v.Struct(i)
}

func ValidateCtx(ctx context.Context, i any) error {
	return v.StructCtx(ctx, i)
}
