package user

import (
	"github.com/google/uuid"
	"github.com/gowool/wool"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/rumorsflow/rumors/v2/internal/entity"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/spf13/cobra"
	"os"
)

type CreateUserDTO struct {
	Username string `json:"username,omitempty" validate:"required,min=3,max=254"`
	Email    string `json:"email,omitempty" validate:"required,email,min=3,max=254"`
	Password string `json:"password,omitempty" validate:"required,min=8,max=64"`
}

func NewCreateCommand() *cobra.Command {
	var dto CreateUserDTO
	var host string

	v := wool.NewValidator()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "New System Administrator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := di.Activators(
				mongodb.Activator("mongo"),
				db.SysUserActivator(),
			); err != nil {
				return err
			}

			if err := v.ValidateCtx(cmd.Context(), dto); err != nil {
				return err
			}

			user := entity.SysUser{
				ID:       uuid.New(),
				Username: dto.Username,
				Email:    dto.Email,
				Password: dto.Password,
			}
			if err := user.GeneratePasswordHash(); err != nil {
				return err
			}
			if err := user.GenerateOTPSecret(20); err != nil {
				return err
			}

			userRepo, err := db.GetSysUserRepository(cmd.Context())
			if err != nil {
				return err
			}

			if err = userRepo.Save(cmd.Context(), &user); err != nil {
				return err
			}

			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      host,
				AccountName: user.Email,
				Secret:      user.OTPSecret,
			})
			if err != nil {
				return err
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

			return nil
		},
	}

	cmd.Flags().StringVarP(&dto.Username, "username", "u", "", "Administrator username")
	cmd.Flags().StringVarP(&dto.Email, "email", "e", "", "Administrator email")
	cmd.Flags().StringVarP(&dto.Password, "password", "p", "", "Administrator password")
	cmd.Flags().StringVar(&host, "host", "localhost", "Host for QR code")

	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}
