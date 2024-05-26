package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
	t := table.NewWriter()
	t.Style().Options.DrawBorder = true
	t.Style().Options.SeparateRows = true
	tTemp := table.Table{}
	tTemp.Render()
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "UniqueId", Align: text.AlignCenter},
		{Name: "TOTP", Align: text.AlignCenter},
		{Name: "Lifetime/s", Align: text.AlignCenter},
		{Name: "Name", Align: text.AlignCenter},
		{Name: "NextTOTP", Align: text.AlignCenter},
	})
	t.AppendHeader(table.Row{"UniqueId", "TOTP", "Lifetime/s", "Name", "NextTOTP"})
	for _, token := range client.cfg.AuthenticatorTokenS {
		tm := time.Now()
		totp, err := token.TOTP(client.cfg.BackupPassword, tm)
		if err != nil {
			log.Fatalf("Failed to generate TOTP: %v", err)
		}
		lifeTime := 30 - tm.Unix()%30
		nextTotp, err := token.TOTP(client.cfg.BackupPassword, tm.Add(time.Second*30))
		if err != nil {
			log.Fatalf("Failed to generate TOTP: %v", err)
		}
		t.AppendRow(table.Row{token.UniqueID, totp, lifeTime, token.Description(), nextTotp})
	}
	fmt.Println(t.Render())
}
