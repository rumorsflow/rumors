package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/roadrunner-server/errors"
	"github.com/rumorsflow/rumors/v2/internal/common"
	"github.com/rumorsflow/rumors/v2/internal/container"
	"github.com/rumorsflow/rumors/v2/internal/db"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/pkg/repository"
	"github.com/spf13/cobra"
	"os"
)

const createUserPluginName = "create_user"

type CreateUserDTO struct {
	Issuer   string `json:"issuer,omitempty" validate:"required"`
	Username string `json:"username,omitempty" validate:"required,min=3,max=254"`
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=254"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

type CreateUserPlugin struct {
	dto  CreateUserDTO
	repo repository.WriteRepository[*entity.SysUser]
}

func (p *CreateUserPlugin) Init(uow common.UnitOfWork) error {
	r, err := uow.Repository((*entity.SysUser)(nil))
	if err != nil {
		const op = errors.Op("create_user_plugin_init")
		return errors.E(op, err)
	}

	p.repo = r.(repository.ReadWriteRepository[*entity.SysUser])

	return nil
}

func (p *CreateUserPlugin) Serve() chan error {
	errCh := make(chan error, 1)

	go execCreateUser(p.repo, p.dto, errCh)

	return errCh
}

func (p *CreateUserPlugin) Stop(context.Context) error {
	return nil
}

func (p *CreateUserPlugin) Name() string {
	return createUserPluginName
}

func execCreateUser(repo repository.WriteRepository[*entity.SysUser], dto CreateUserDTO, ch chan<- error) {
	const op = errors.Op("create_user_command")

	v := wool.NewValidator()

	if err := v.Validate(dto); err != nil {
		ch <- errors.E(op, err)
		return
	}

	user := entity.SysUser{
		ID:       uuid.New(),
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
	}
	if err := user.GeneratePasswordHash(); err != nil {
		ch <- errors.E(op, err)
		return
	}
	if err := user.GenerateOTPSecret(20); err != nil {
		ch <- errors.E(op, err)
		return
	}

	if err := repo.Save(context.Background(), &user); err != nil {
		ch <- errors.E(op, err)
		return
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      dto.Issuer,
		AccountName: user.Email,
		Secret:      user.OTPSecret,
	})
	if err != nil {
		ch <- errors.E(op, err)
		return
	}

	config := qrterminal.Config{
		Level:          qrterminal.L,
		Writer:         os.Stdout,
		HalfBlocks:     true,
		BlackChar:      qrterminal.BLACK_BLACK,
		WhiteBlackChar: qrterminal.WHITE_BLACK,
		WhiteChar:      qrterminal.WHITE_WHITE,
		BlackWhiteChar: qrterminal.BLACK_WHITE,
		QuietZone:      1,
	}
	qrterminal.GenerateWithConfig(key.String(), config)

	ch <- common.Success
}

func NewCreateCommand() *cobra.Command {
	var dto CreateUserDTO

	cmd := &cobra.Command{
		Use:   "create",
		Short: "New System Administrator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Context().Value("container").(*container.Container).Run(
				&db.Plugin{},
				&CreateUserPlugin{dto: dto},
			)
		},
	}

	cmd.Flags().StringVarP(&dto.Username, "username", "u", "", "Administrator username")
	cmd.Flags().StringVarP(&dto.Email, "email", "e", "", "Administrator email")
	cmd.Flags().StringVarP(&dto.Password, "password", "p", "", "Administrator password")
	cmd.Flags().StringVar(&dto.Issuer, "host", "localhost", "Host for QR code")

	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
