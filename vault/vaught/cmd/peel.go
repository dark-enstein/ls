/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/spf13/cobra"
)

// peelCmd represents the cli command
var peelCmd = &cobra.Command{
	Use:   "peel",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := vlog.New(debug)
		if len(id) == 0 {
			logger.Logger().Fatal().Msg("please pass in the id of the token to be retrieved")
			return
		}

		fmt.Printf("Peek=ling record with ID %s\n", id)

		ctx := context.Background()

		ic, err := jsonDecode(DefaultConfigLoc, logger)
		if err != nil {
			logger.Logger().Fatal().Msgf("error occurred while reading config: %s", err)
			return
		}

		manager, err := ic.Manager(ctx, logger)
		if err != nil {
			logger.Logger().Fatal().Msgf("error setting up manager: %s", err)
			return
		}

		token, err := manager.GetTokenByID(ctx, id)
		if err != nil {
			logger.Logger().Fatal().Msgf("error retrieving token: %s", err)
			return
		}

		jsonByte, err := json.Marshal(token)
		if err != nil {
			logger.Logger().Fatal().Msgf("error marshalling token into json: %s", err)
			return
		}

		fmt.Println("All Tokens:")
		fmt.Println(string(jsonByte))
	},
}

func init() {
	//cmd.RootCmd.AddCommand(CliCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cliCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cliCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
