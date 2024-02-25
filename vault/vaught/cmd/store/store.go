/*
Copyright © 2024 Ayobami Bamigboye <ayo@greystein.com>
*/
package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/vault/internal/model"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/dark-enstein/vault/vaught/cmd/helper"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type StoreOptions struct {
	id    string
	val   string
	debug bool
}

const (
	FlagID    = "id"
	FlagValue = "value"
)

// NewStoreCmd represents the cli command
func NewStoreCmd() *cobra.Command {

	sop := &StoreOptions{}

	storeCmd := &cobra.Command{
		Use:   "store",
		Short: "Stores a new token with the specified value",
		Long: `The 'store' command securely stores a new token in the vault with a specified value, identified by a unique ID. 
This command is essential for adding new secrets to the vault, providing a safe way to manage sensitive information.

Usage:

  vault store --id <token-id> --value <secret-value>

Replace '<token-id>' with the unique identifier for the new token, and '<secret-value>' with the actual secret information you wish to store. The command securely processes and stores the token in the configured storage backend, ensuring the confidentiality and integrity of your secret data.

Examples:
Store a new token:
  vault store --id "1234abcd" --value "mySecretData"

Ensure to initialize the vault using 'vault init' before storing any tokens to set up the necessary configurations and storage backend.`,
		Run: func(cmd *cobra.Command, args []string) {
			debug, err := cmd.Flags().GetBool("debug")
			if err != nil {
				log.Error().Msgf("error retrieving persistent flag: %s: %w", "debug", err)
			}
			logger := vlog.New(debug)
			ctx := context.Background()
			bytes, err := sop.Run(ctx, logger)

			// Resolve persistent flags

			sop.debug = debug
			if err != nil {
				// revisit this
				log.Fatal().Msgf(err.Error())
			}

			fmt.Printf("Stored token with id: %s\n", sop.id)
			fmt.Println(string(bytes))
		},
	}

	storeCmd.Flags().StringVarP(&sop.id, FlagID, "i", "", "specify token ID to be stored")
	storeCmd.Flags().StringVarP(&sop.id, FlagValue, "v", "", "specify token ID to be stored")
	storeCmd.MarkFlagsRequiredTogether(FlagID, FlagValue)
	return storeCmd
}

func (sop *StoreOptions) Run(ctx context.Context, logger *vlog.Logger) ([]byte, error) {
	fmt.Println("Storing token with id:", sop.id)
	var err error

	// check if config exists, override
	ic := helper.NewInstanceConfig()
	err = ic.JsonDecode()
	if err != nil {
		return nil, err
	}

	// initialize token manager
	manager, err := ic.Manager(ctx)
	if err != nil {
		logger.Logger().Debug().Msgf("error retrieving tokens from store: %s", err)
		return nil, err
	}

	token, err := manager.Tokenize(ctx, sop.id, sop.val)
	if err != nil {
		logger.Logger().Fatal().Msgf("error retrieving token: %s", err)
		return nil, err
	}

	// reset request struct
	sop.id = ""

	// capture token into struct
	tokenResp := &model.Child{
		Key:   sop.id,
		Value: token,
	}

	jsonByte, err := json.Marshal(tokenResp)
	if err != nil {
		logger.Logger().Fatal().Msgf("error marshalling token into json: %s", err)
		return nil, err
	}

	logger.Logger().Info().Msgf("Successfully set up Vault CLI. You can begin using the other commands.")

	return jsonByte, nil
}
