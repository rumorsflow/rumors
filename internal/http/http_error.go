package http

import (
	"errors"
	"github.com/gowool/wool"
	"github.com/rumorsflow/rumors/v2/internal/repository"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func ErrorTransform(err error) *wool.Error {
	var e *wool.Error
	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, repository.ErrEntityNotFound) {
		e = wool.NewError(http.StatusNotFound, err)
	} else if mongo.IsDuplicateKeyError(err) || errors.Is(err, repository.ErrDuplicateKey) {
		e = wool.NewError(http.StatusConflict, err)
	} else if mongo.IsTimeout(err) {
		e = wool.NewError(http.StatusGatewayTimeout, err)
	} else {
		var ok bool
		if e, ok = err.(*wool.Error); !ok {
			e = wool.NewError(http.StatusInternalServerError, err)
		}
	}
	return e
}
