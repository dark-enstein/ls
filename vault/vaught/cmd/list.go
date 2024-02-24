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

// listCmd represents the CLI command for listing all stored tokens
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all stored tokens in the vault",
	Long: `The 'list' command retrieves and displays all tokens currently stored in the vault. 
This command is useful for getting an overview of all the secrets managed by the vault system. 

Usage:

  vault list

This will output all the tokens stored, formatted as JSON for easy reading and integration with other tools. Ensure you have the appropriate permissions and the vault is correctly configured before running this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listing records in vault")

		ctx := context.Background()

		logger := vlog.New(debug)

		ic, err := jsonDecode(DefaultConfigLoc, logger)
		if err != nil {
			logger.Logger().Fatal().Msgf("error occurred while reading config: %s", err)
		}

		manager, err := ic.Manager(ctx, logger)
		if err != nil {
			logger.Logger().Fatal().Msgf("error setting up manager: %s", err)
		}

		tokens, err := manager.GetAllTokens(ctx)
		if err != nil {
			logger.Logger().Fatal().Msgf("error retrieving tokens: %s", err)
		}

		bytesResult, err := json.Marshal(&tokens)
		if err != nil {
			logger.Logger().Fatal().Msgf("error marshalling tokens into json: %s", err)
		}

		fmt.Println("All Tokens:")
		fmt.Println(string(bytesResult))
	},
}

func init() {
	initCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable or disable debug mode.")
}
