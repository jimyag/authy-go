package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all TOTP",
	Run:   listRun,
}

func listRun(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(client.cfgFile); os.IsNotExist(err) {
		syncCmd.Run(cmd, args)
	}
	if err := client.cfg.Load(client.cfgFile); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	for _, token := range client.cfg.AuthenticatorTokenS {
		totp, err := token.TOTP(client.cfg.BackupPassword, time.Now())
		if err != nil {
			log.Fatalf("Failed to generate TOTP: %v", err)
		}
		_, _ = fmt.Printf("%s: %s \t:%s\n", token.UniqueID, totp, token.Description())
	}
}
