package util

import (
	"encoding/json"
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/mongo-ext"
	"github.com/rumorsflow/rumors/internal/pkg/validator"
	"net/http"
)

func GetId(r *http.Request) string {
	return flow.Param(r.Context(), "id")
}

func GetCriteria(r *http.Request) mongoext.Criteria {
	criteria, _ := mongoext.GetC(r.URL.RawQuery, "filters")
	return criteria
}

func Bind(r *http.Request, dto any) {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(err)
	}

	if err := validator.ValidateCtx(r.Context(), dto); err != nil {
		panic(err)
	}
}
