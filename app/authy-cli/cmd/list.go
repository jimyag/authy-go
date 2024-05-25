package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/jimyag/authy-go"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all TOTP",
	Run:   listRun,
}

func listRun(cmd *cobra.Command, args []string) {
	if err := client.cfg.Load(client.cfgFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if client.cfg.Empty() {
		log.Fatal("Please run authy-cli init first")
	}

	for _, token := range client.cfg.AuthenticatorTokenS {
		decrypted, err := token.Decrypt(client.cfg.BackupPassword)
		if err != nil {
			log.Fatalf("Failed to decrypt token %s: %v", token.Description(), err)
		}
		totp, err := authy.GenerateTOTP([]byte(decrypted), time.Now(), token.Digits, 30)
		if err != nil {
			log.Fatalf("Failed generate TOTP %s", err)
		}
		_, _ = fmt.Printf("%s: %s \t:%s\n", token.UniqueID, totp, token.Description())
	}
}