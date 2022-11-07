package errorhandler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	rrerr "github.com/roadrunner-server/errors"
	"github.com/rumorsflow/contracts/config"
	"github.com/rumorsflow/rumors/internal/api/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
)

const (
	RootPluginName = "http"
	PluginName     = "error_handler"
)

type Plugin struct {
	log *zap.Logger
}

func (p *Plugin) Init(cfg config.Configurer, log *zap.Logger) error {
	const op = rrerr.Op("http error handler plugin init")

	if !cfg.Has(RootPluginName) {
		return rrerr.E(op, rrerr.Disabled)
	}

	p.log = log

	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Handle(w http.ResponseWriter, err error) {
	var he *util.HTTPError

	if e, ok := err.(*util.HTTPError); ok {
		he = e
	} else {
		if ve, ok := err.(validator.ValidationErrors); ok {
			he = util.NewHTTPError(http.StatusUnprocessableEntity, nil)
			data := make(map[string][]string)
			for _, fe := range ve {
				if _, ok = data[fe.Field()]; ok {
					data[fe.Field()] = append(data[fe.Field()], fe.Error())
				} else {
					data[fe.Field()] = []string{fe.Error()}
				}
			}
			he.Data = data
		} else if errors.Is(err, mongo.ErrNoDocuments) {
			he = util.NewHTTPError(http.StatusNotFound, err)
		} else if mongo.IsDuplicateKeyError(err) {
			he = util.NewHTTPError(http.StatusConflict, err)
		} else if mongo.IsTimeout(err) {
			he = util.NewHTTPError(http.StatusGatewayTimeout, err)
		} else {
			he = util.NewHTTPError(http.StatusInternalServerError, err)
		}
	}

	if util.IsDebug(p.log) && he.Internal != nil {
		he.Developer = he.Internal.Error()
	}

	util.JSON(w, he.Code, he)
}
