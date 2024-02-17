package service

import (
	"context"
	"github.com/dark-enstein/vault/cmd/cmd"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/dark-enstein/vault/service"

	"github.com/spf13/cobra"
)

var (
	helpIntroService = `
<Vault Help Goes here>
`
)

type ServiceCMD struct {
	port     string
	debug    bool
	logger   *vlog.Logger
	logLevel int
}

var config = ServiceCMD{}

// serviceCmd represents the service command
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Start Vault Service",
	Long:  `Start Vault Service`,
	Run: func(cmd *cobra.Command, args []string) {
		log := config.logger.Logger()

		// validate argument length
		if len(args) > 0 {
			log.Error().Msgf("do not expect arguments when calling: vault service\n%s\n", helpIntroService)
		}

		// start service
		ctx := context.Background()
		srv := service.New(ctx, config.logger)
		if err := srv.Run(ctx); err != nil {
			log.Fatal().Msgf("error while service is starting: %s\n", err.Error())
		}
	},
}

func init() {
	vault.RootVaultCmd.AddCommand(serviceCmd)

	serviceCmd.PersistentFlags().BoolVarP(&config.debug, "debug", "d", true, "Enable/Disable debug mode")

	serviceCmd.Flags().StringVarP(&config.port, "port", "p", "8080", "Specify the port you would like vault service to listen on")
}
