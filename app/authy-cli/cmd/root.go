package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/jimyag/authy-go"
)

var rootCmd = &cobra.Command{
	Use:   "authy-cli",
	Short: "Authy CLI",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var client *Client

const (
	appName    = "authy-cli"
	configPath = "authy-go.json"
)

type Client struct {
	authyCli *authy.Client
	cfg      *Config
	cfgFile  string
}

func New() (*Client, error) {
	authyCli, err := authy.New()
	if err != nil {
		return nil, err
	}
	cli := &Client{
		authyCli: &authyCli,
		cfg:      &Config{},
		cfgFile:  "",
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get home directory: %w", err)
	}
	// ~/.config/authy-cli/authy-go.json
	cfgPath := filepath.Join(homePath, ".config", appName, configPath)
	cli.cfgFile = cfgPath
	_, err = os.Stat(cfgPath)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(filepath.Dir(cfgPath), 0755); err != nil {
			return nil, fmt.Errorf("could not create config directory: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("could not stat config file: %w", err)
	}
	return cli, nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)
	var err error
	client, err = New()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// We'll persist this to the filesystem so we don't need to
// re-register the device every time
type Config struct {
	UserID              uint64                     `json:"user_id"`
	DeviceID            uint64                     `json:"device_id"`
	Seed                string                     `json:"seed"`
	APIKey              string                     `json:"api_key"`
	BackupPassword      string                     `json:"backup_password"`
	AuthenticatorTokenS []authy.AuthenticatorToken `json:"authenticator_tokens"`
}

func (c *Config) Save(file string) error {
	if file == "" {
		return fmt.Errorf("file is empty")
	}
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(c)
}

func (c *Config) Load(file string) error {
	if file == "" {
		return fmt.Errorf("file is empty")
	}
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(c)
}

func (c *Config) Empty() bool {
	return c.UserID == 0 ||
		c.DeviceID == 0 ||
		c.APIKey == "" ||
		c.Seed == "" ||
		c.BackupPassword == "" ||
		len(c.AuthenticatorTokenS) == 0
}
