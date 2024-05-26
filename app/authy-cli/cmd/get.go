package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/jimyag/authy-go"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get TOTP by UniqueID",
	Run:   getRun,
}

func getRun(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(client.cfgFile); os.IsNotExist(err) {
		syncCmd.Run(cmd, args)
	}
	if err := client.cfg.Load(client.cfgFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var retToken authy.AuthenticatorToken
	for _, token := range client.cfg.AuthenticatorTokenS {
		if token.UniqueID == args[0] {
			retToken = token
			break
		}
	}
	if retToken.UniqueID == "" {
		log.Fatal("Token not found")
	}
	totp, err := retToken.TOTP(client.cfg.BackupPassword, time.Now())
	if err != nil {
		log.Fatalf("Failed to generate TOTP: %v", err)
	}
	fmt.Println(totp)
}
