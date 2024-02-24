package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dark-enstein/vault/internal/store"
	"github.com/dark-enstein/vault/internal/tokenize"
	"github.com/dark-enstein/vault/internal/vlog"
	"github.com/dark-enstein/vault/service"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)

var (
	DefaultRootConfig = "~/.vault/cli"
	DefaultConfigLoc  = filepath.Join(DefaultRootConfig, "config")
	DefaultCipherLoc  = filepath.Join(DefaultRootConfig, ".cipher")
	DefaultStoreLoc   = filepath.Join(DefaultRootConfig, ".store")
	DefaultGobLoc     = filepath.Join(DefaultRootConfig, ".gob")
)

// initCmd represents the command for initializing the vault system
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes your local system for on-demand secret encryption and storage using vault",
	Long: `The 'init' command sets up your local system for on-demand secret encryption and storage using vault.
This command supports multiple storage backends, allowing for flexible configurations based on the environment and requirements.

Supported storage options include:
- File: Utilizes the local file system for persistent storage.
- Gob: Employs GOB file storage for serialization of Go data structures.
- Redis: Connects to a Redis server for distributed storage and caching.
- In-memory map: Uses a concurrent-safe map for in-memory storage, ideal for temporary data and testing.

Examples:
Initialize the service with file storage:
  vault init --store file --fileLoc /path/to/store

Initialize the service with Redis:
  vault init --store redis --connectionString "redis://user:password@localhost:6379"

Each storage option offers specific flags for customization, providing the flexibility to adapt to various deployment scenarios.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing vault cli")

		// check if config exists, override
		logger := vlog.New(true)

		_ = context.Background()
		var err error

		DefaultRootConfig, err = homedir.Expand(DefaultRootConfig)
		if err != nil {
			logger.Logger().Fatal().Msgf("error occurred while setting up cli: %w", err)
		}

		ic := InstanceConfig{
			ID:        xid.New().String(),
			CipherLoc: DefaultCipherLoc,
			StoreType: storeStr,
			Debug:     debug,
			LastUse:   time.Now().UnixNano(),
		}

		switch storeStr {
		case service.STORE_FILE:
			log.Info().Msg("Using File storage")
			ic.StoreLoc, err = filepath.Abs(fileLoc)
			if err != nil {
				log.Fatal().Msgf("error with path: %s\n", err)
			}
		case service.STORE_GOB:
			log.Info().Msg("Using Gob storage")
			ic.StoreLoc, err = filepath.Abs(gobLoc)
			if err != nil {
				log.Fatal().Msgf("error with path: %s\n", err)
			}
		case service.STORE_REDIS:
			log.Info().Msg("Using Redis storage")
			_, err = redis.ParseURL(redisConnString)
			if err != nil {
				log.Fatal().Msgf("error with redis string: %s\n", err)
			}
			ic.RedisString = redisConnString
		case service.STORE_MAP:
			log.Info().Msg("Using In-memory map storage")
		}
		if err != nil {
			logger.Logger().Fatal().Msgf("error while setting up service: %s", err)
		}

		// persist to disk at config loc
		err = jsonEncode(DefaultConfigLoc, ic, logger)
		if err != nil {
			logger.Logger().Fatal().Msgf("error occurred while setting up cli: %s", err)
		}

		logger.Logger().Info().Msgf("Successfully set up Vault CLI. You can begin using the other commands.")
	},
}

var storeStr string
var redisConnString string
var gobLoc string
var fileLoc string
var debug bool

func init() {
	initCmd.Flags().StringVarP(&storeStr, "store", "s", "file", "Specify the storage backend for the service. Options: file, gob, redis, in-memory syncmap.")
	initCmd.Flags().StringVarP(&redisConnString, "connectionString", "c", store.DefaultRedisConnectionString, "Specify the Redis connection string.")
	initCmd.Flags().StringVarP(&gobLoc, "gobLoc", "g", DefaultGobLoc, "Specify the disk location for the gob store.")
	initCmd.Flags().StringVarP(&fileLoc, "fileLoc", "f", DefaultStoreLoc, "Specify the disk location for the file store.")
	//initCmd.Flags().StringVarP(&configPath, "configPath", "c", DefaultConfigLoc, "Specify the disk location for the config file.")
	initCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable or disable debug mode.")
}

//var configPath string

type InstanceConfig struct {
	ID          string `json:"id"`
	CipherLoc   string `json:"cipher_loc"`
	StoreLoc    string `json:"store_loc"`
	RedisString string `json:"redis_string"`
	StoreType   string `json:"store_type"`
	Debug       bool   `json:"debug"`
	LastUse     int64  `json:"last_login"`
}

func (ic *InstanceConfig) Manager(ctx context.Context, logger *vlog.Logger) (*tokenize.Manager, error) {
	log := logger.Logger()
	switch storeStr {
	case service.STORE_FILE:
		log.Info().Msg("Using File storage")
		return tokenize.NewManager(ctx, logger, tokenize.WithStore(store.NewFile(ic.StoreLoc, logger))), nil
	case service.STORE_GOB:
		log.Info().Msg("Using Gob storage")
		gob, err := store.NewGob(ctx, ic.StoreLoc, logger, false)
		if err != nil {
			log.Fatal().Msgf("error while creating storage backend: %s", err)
		}
		return tokenize.NewManager(ctx, logger, tokenize.WithStore(gob)), nil
	case service.STORE_REDIS:
		log.Info().Msg("Using Redis storage")
		r, err := store.NewRedis(ic.RedisString, logger)
		if err != nil {
			log.Fatal().Msgf("error while creating storage backend: %s", err)
		}
		return tokenize.NewManager(ctx, logger, tokenize.WithStore(r)), nil
	case service.STORE_MAP:
		log.Info().Msg("Using In-memory map storage")
		return tokenize.NewManager(ctx, logger, tokenize.WithStore(store.NewSyncMap(ctx, logger))), nil
	default:
		return nil, errors.New("unrecognized store parameter")
	}
}

func jsonEncode(path string, ic InstanceConfig, logger *vlog.Logger) error {
	log := logger.Logger()

	log.Info().Msgf("setting up vault config at location: %s", path)

	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Error().Msgf("encountered error while opening file at location %s: %s", path, err)
		return err
	}
	defer fd.Close()

	jsonBytes, err := json.Marshal(&ic)
	if err != nil {
		log.Error().Msgf("encountered error while encoding json %s: %s", path, err)
		return err
	}

	_, err = fd.Write(jsonBytes)
	if err != nil {
		log.Error().Msgf("encountered error while writing config to file location %s: %s", path, err)
		return err
	}

	return nil
}
func jsonDecode(path string, logger *vlog.Logger) (*InstanceConfig, error) {
	log := logger.Logger()

	log.Info().Msgf("setting up vault config at location: %s", path)

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Error().Msgf("encountered error while reading config file at location %s: %s", path, err)
		return nil, err
	}

	var ic = InstanceConfig{}
	err = json.Unmarshal(fileBytes, &ic)
	if err != nil {
		log.Error().Msgf("json config invalid: %s", err)
		return nil, err
	}

	return &ic, nil
}
