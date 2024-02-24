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

// peekCmd represents the CLI command for viewing a specific token
var peekCmd = &cobra.Command{
	Use:   "peek",
	Short: "Displays details of a specific token by ID",
	Long: `The 'peek' command retrieves and displays the details of a specific token stored in the vault, identified by its ID. 
This is useful for quickly viewing the properties of a particular token without needing to list all tokens.

Usage:

  vault peek --id <token-id>

Replace '<token-id>' with the actual ID of the token you wish to view. The token's details will be displayed in JSON format, providing comprehensive information about the token's attributes.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := vlog.New(debug)
		if len(id) == 0 {
			logger.Logger().Fatal().Msg("please pass in the id of the token to be retrieved")
			return
		}

		fmt.Printf("Peeking record with ID %s\n", id)

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

var id string

func init() {
	initCmd.Flags().StringVarP(&id, "id", "i", "", "Specify token ID to be retrieved.")
}
