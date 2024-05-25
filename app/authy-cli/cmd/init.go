package cmd

import (
	"bufio"
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/jimyag/authy-go"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init Authy config",
	Run:   initRun,
}

func initRun(cmd *cobra.Command, args []string) {
	sc := bufio.NewScanner(os.Stdin)
	log.Printf("\nWhat is your Authy Id ? (digits only): ")
	if !sc.Scan() {
		log.Fatalln("Please provide an authyId code ")
	}
	authyId, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
	if err != nil {
		log.Fatalln("Please provide an authyId code")
	}
	client.cfg.UserID = uint64(authyId)
	regStart, err := client.authyCli.RequestDeviceRegistration(context.Background(), client.cfg.UserID, authy.ViaMethodPush)
	if err != nil {
		log.Fatalln("Failed to request device registration ,", err)
	}

	if !regStart.Success {
		log.Fatalln("authy did not accept the device registration request:", regStart)
	}

	// Poll for a while until the user has responded to the device registration
	var regPIN string
	timeout := time.Now().Add(5 * time.Minute)
	for {
		if timeout.Before(time.Now()) {
			log.Fatalln("gave up waiting for user to respond to Authy device registration request")
		}

		log.Printf("Checking device registration status (%s until we give up)\n", time.Until(timeout).Truncate(time.Second))
		regStatus, err := client.authyCli.CheckDeviceRegistration(context.Background(), client.cfg.UserID, regStart.RequestID)
		if err != nil {
			log.Fatalln("Failed to check device registration status ,", err)
		}

		if regStatus.Status == authy.RegistrationStatusAccepted {
			regPIN = regStatus.PIN
			break
		} else if regStatus.Status != authy.RegistrationStatusPending {
			log.Fatalln("invalid status while waiting for device registration , status:", regStatus.Status)
		}
		time.Sleep(5 * time.Second)
	}

	// We have the registration PIN, complete the registration
	regComplete, err := client.authyCli.CompleteDeviceRegistration(context.Background(), client.cfg.UserID, regPIN)
	if err != nil {
		log.Fatalln("Failed to complete device registration,", err)
	}

	if regComplete.Device.SecretSeed == "" {
		log.Fatalln("something went wrong completing the device registration")
	}

	log.Println("Please provide your Authy TOTP backup password: ")
	pp, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("Failed to read the password")
	}
	client.cfg.BackupPassword = strings.TrimSpace(string(pp))
	client.cfg.UserID = regComplete.AuthyID
	client.cfg.Seed = regComplete.Device.SecretSeed
	client.cfg.DeviceID = regComplete.Device.ID
	client.cfg.APIKey = regComplete.Device.APIKey

	respApps, err := client.authyCli.QueryAuthenticatorApps(context.Background(), client.cfg.UserID, client.cfg.DeviceID, client.cfg.Seed)
	if err != nil {
		log.Fatalln("Could not fetch authenticator apps ,", err)
	}
	if !respApps.Success {
		log.Fatalln("Failed to fetch authenticator apps ,", respApps)
	}

	// Fetch the actual tokens now
	tokens, err := client.authyCli.QueryAuthenticatorTokens(context.Background(), client.cfg.UserID, client.cfg.DeviceID, client.cfg.Seed)
	if err != nil {
		log.Fatalln("Could not fetch authenticator tokens ,", err)
	}
	if !tokens.Success {
		log.Fatalln("Failed to fetch authenticator tokens ,", tokens)
	}

	client.cfg.AuthenticatorTokenS = tokens.AuthenticatorTokens
	if err = client.cfg.Save(client.cfgFile); err != nil {
		log.Fatalln("Failed to save config, ", err)
	}
}
