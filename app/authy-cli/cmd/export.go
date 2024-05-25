package cmd

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export authy tokens",
	Run:   exportRun,
}

func exportRun(cmd *cobra.Command, args []string) {
	if err := client.cfg.Load(client.cfgFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if client.cfg.Empty() {
		log.Fatal("Please run authy-cli init first")
	}
	for _, token := range client.cfg.AuthenticatorTokenS {
		decrypted, err := token.Decrypt(client.cfg.BackupPassword)
		if err != nil {
			log.Printf("Failed to decrypt token %s: %v", token.Description(), err)
			continue
		}
		params := url.Values{}
		params.Set("secret", decrypted)
		params.Set("digits", strconv.Itoa(token.Digits))
		u := url.URL{
			Scheme:   "otpauth",
			Host:     "totp",
			Path:     token.Description(),
			RawQuery: params.Encode(),
		}
		fmt.Println(u.String())
	}
}
