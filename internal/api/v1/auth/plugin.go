package auth

import (
	"bytes"
	"encoding/base64"
	"github.com/alexedwards/flow"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/rumorsflow/rumors/internal/api/middleware/authz"
	"github.com/rumorsflow/rumors/internal/api/middleware/jwt"
	"github.com/rumorsflow/rumors/internal/api/util"
	"github.com/rumorsflow/rumors/internal/models"
	"github.com/rumorsflow/rumors/internal/services/token"
	"github.com/rumorsflow/rumors/internal/storage"
	"image/png"
	"net/http"
)

const PluginName = "/api/v1/auth"

type Plugin struct {
	userStorage  storage.UserStorage
	tokenService token.Service
}

func (p *Plugin) Init(userStorage storage.UserStorage, tokenService token.Service) error {
	p.userStorage = userStorage
	p.tokenService = tokenService
	return nil
}

// Name returns user-friendly plugin name
func (p *Plugin) Name() string {
	return PluginName
}

func (p *Plugin) Register(mux *flow.Mux) {
	mux.HandleFunc(PluginName+"/sign-up", p.signUp, http.MethodPost)
	mux.HandleFunc(PluginName+"/sign-in", p.signIn, http.MethodPost)
	mux.HandleFunc(PluginName+"/otp", authz.Is(p.otp), http.MethodPost)
	mux.HandleFunc(PluginName+"/refresh", p.refresh, http.MethodPost)
}

func (p *Plugin) signUp(w http.ResponseWriter, r *http.Request) {
	var dto SignUpRequest
	util.Bind(r, &dto)

	user := models.User{
		Id:       uuid.NewString(),
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Roles:    []models.Role{models.UserRole},
	}
	if err := user.GeneratePasswordHash(); err != nil {
		panic(err)
	}
	if err := user.GenerateOTPSecret(20); err != nil {
		panic(err)
	}

	host := r.URL.Host
	if host == "" {
		host = "localhost"
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      host,
		AccountName: user.Email,
		Secret:      user.OTPSecret,
	})
	if err != nil {
		panic(err)
	}

	var qr bytes.Buffer
	img, err := key.Image(300, 300)
	if err != nil {
		panic(err)
	}
	if err = png.Encode(&qr, img); err != nil {
		panic(err)
	}

	if err = p.userStorage.Save(r.Context(), &user); err != nil {
		panic(err)
	}

	util.JSON(w, http.StatusCreated, SignUpResponse{
		Uri: key.URL(),
		Qr:  "data:image/png;base64," + base64.StdEncoding.EncodeToString(qr.Bytes()),
	})
}

func (p *Plugin) signIn(w http.ResponseWriter, r *http.Request) {
	var dto SignInRequest
	util.Bind(r, &dto)

	user, err := p.userStorage.FindByUsername(r.Context(), dto.Username)
	if err != nil {
		util.BadRequest(err)
	}

	if err = user.CheckPassword(dto.Password); err != nil {
		util.BadRequest(err)
	}

	session, err := p.tokenService.SessByUser(r.Context(), user, false)
	if err != nil {
		util.BadRequest(err)
	}

	util.OK(w, session)
}

func (p *Plugin) otp(w http.ResponseWriter, r *http.Request) {
	var dto OtpRequest
	util.Bind(r, &dto)

	user, _ := jwt.CtxUser(r.Context())
	if err := user.CheckTOTP(dto.Password); err != nil {
		util.BadRequest(err)
	}

	session, err := p.tokenService.SessByUser(r.Context(), user, true)
	if err != nil {
		util.BadRequest(err)
	}

	util.OK(w, session)
}

func (p *Plugin) refresh(w http.ResponseWriter, r *http.Request) {
	var dto RefreshTokenRequest
	util.Bind(r, &dto)

	session, err := p.tokenService.SessByRefreshToken(r.Context(), dto.RefreshToken)
	if err != nil {
		util.BadRequest(err)
	}

	util.OK(w, session)
}
