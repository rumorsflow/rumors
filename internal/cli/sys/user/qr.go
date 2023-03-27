package user

import (
	"errors"
	"fmt"
	"github.com/mdp/qrterminal/v3"
	"github.com/pquerna/otp/totp"
	"github.com/rumorsflow/rumors/v2/internal/repository/db"
	"github.com/rumorsflow/rumors/v2/pkg/di"
	"github.com/rumorsflow/rumors/v2/pkg/mongodb"
	"github.com/spf13/cobra"
	"os"
)

const filter = "index=0&size=1&field.0.0=username&value.0.0=%s&field.0.1=email&value.0.1=%s"

func NewQRCommand() *cobra.Command {
	var username, email, host string

	cmd := &cobra.Command{
		Use:   "qr",
		Short: "OTP QR code",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := di.Activators(
				mongodb.Activator("mongo"),
				db.SysUserActivator(),
			); err != nil {
				return err
			}

			userRepo, err := db.GetSysUserRepository(cmd.Context())
			if err != nil {
				return err
			}

			users, err := userRepo.Find(cmd.Context(), db.BuildCriteria(fmt.Sprintf(filter, username, email)))
			if err != nil {
				return err
			}

			if len(users) != 1 {
				return errors.New("user not found")
			}

			user := users[0]

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

	cmd.Flags().StringVarP(&username, "username", "u", "", "Administrator username")
	cmd.Flags().StringVarP(&email, "email", "e", "", "Administrator email")
	cmd.Flags().StringVar(&host, "host", "localhost", "Host for QR code")

	return cmd
}
