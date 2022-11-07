package util

import (
	"encoding/json"
	"github.com/alexedwards/flow"
	"github.com/rumorsflow/rumors/internal/pkg/validator"
	"net/http"
)

func GetId(r *http.Request) string {
	return flow.Param(r.Context(), "id")
}

func Bind(r *http.Request, dto any) {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(err)
	}

	if err := validator.ValidateCtx(r.Context(), dto); err != nil {
		panic(err)
	}
}
