package user

import (
	"context"
	"fmt"
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

const (
	qrPluginName = "qr"
	filter       = "index=0&size=1&field.0.0=username&value.0.0=%s&field.0.1=email&value.0.1=%s"
)

type QrDTO struct {
	Username string
	Email    string
	Host     string
}

type QRPlugin struct {
	dto  QrDTO
	repo repository.ReadRepository[*entity.SysUser]
}

func (p *QRPlugin) Init(uow common.UnitOfWork) error {
	r, err := uow.Repository((*entity.SysUser)(nil))
	if err != nil {
		const op = errors.Op("qr_plugin_init")
		return errors.E(op, err)
	}

	p.repo = r.(repository.ReadWriteRepository[*entity.SysUser])

	return nil
}

func (p *QRPlugin) Serve() chan error {
	errCh := make(chan error, 1)

	go execQr(p.repo, p.dto, errCh)

	return errCh
}

func (p *QRPlugin) Stop(context.Context) error {
	return nil
}

func (p *QRPlugin) Name() string {
	return qrPluginName
}

func execQr(repo repository.ReadRepository[*entity.SysUser], dto QrDTO, ch chan<- error) {
	const op = errors.Op("qr_user_command")

	users, err := repo.Find(context.Background(), db.BuildCriteria(fmt.Sprintf(filter, dto.Username, dto.Email)))
	if err != nil {
		ch <- errors.E(op, err)
		return
	}

	if len(users) != 1 {
		ch <- errors.E(op, "user not found")
		return
	}

	user := users[0]

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      dto.Host,
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

func NewQRCommand() *cobra.Command {
	var dto QrDTO

	cmd := &cobra.Command{
		Use:   "qr",
		Short: "OTP QR code",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Context().Value("container").(*container.Container).Run(
				&db.Plugin{},
				&QRPlugin{dto: dto},
			)
		},
	}

	cmd.Flags().StringVarP(&dto.Username, "username", "u", "", "Administrator username")
	cmd.Flags().StringVarP(&dto.Email, "email", "e", "", "Administrator email")
	cmd.Flags().StringVar(&dto.Host, "host", "localhost", "Host for QR code")

	return cmd
}
