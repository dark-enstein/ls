/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package srv

import (
	"context"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/dark-enstein/vault/service"
	"github.com/spf13/cobra"
)

// runCmd represents the service command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("service run called")

		logger := vlog.New(true)

		ctx := context.Background()

		srv := service.New(ctx, logger)
		if err := srv.Run(ctx); err != nil {
			logger.Logger().Fatal().Msgf("error while service is starting: %s\n", err.Error())
		}
	},
}

var port string

func init() {
	//cmd.RootCmd.AddCommand(serviceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	ServiceCmd.Flags().StringVarP(&port, "port", "p", "8080", "Specify port for service to listen on")
}
