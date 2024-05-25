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
	decrypted, err := retToken.Decrypt(client.cfg.BackupPassword)
	if err != nil {
		log.Fatalf("Failed to decrypt token %s: %v", retToken.Description(), err)
	}
	totp, err := authy.GenerateTOTP([]byte(decrypted), time.Now(), retToken.Digits, 30)
	if err != nil {
		log.Fatalf("Failed generate TOTP %s", err)
	}
	fmt.Println(totp)
}
