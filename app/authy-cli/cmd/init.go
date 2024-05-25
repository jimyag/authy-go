package cmd

import (
	"bufio"
	"context"
	"fmt"
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
	fmt.Print("\nWhat is your Authy Id ? (digits only): ")
	if !sc.Scan() {
		fmt.Println("Please provide an authyId code")
		return
	}
	authyId, err := strconv.Atoi(strings.TrimSpace(sc.Text()))
	if err != nil {
		fmt.Println("Please provide an authyId code")
		return
	}
	client.cfg.UserID = uint64(authyId)
	regStart, err := client.authyCli.RequestDeviceRegistration(context.Background(), client.cfg.UserID, authy.ViaMethodPush)
	if err != nil {
		log.Fatalf("Failed to request device registration :%s\n", err)
		return
	}

	if !regStart.Success {
		_, _ = fmt.Printf("authy did not accept the device registration request: %+v\n", regStart)
		return
	}

	// Poll for a while until the user has responded to the device registration
	var regPIN string
	timeout := time.Now().Add(5 * time.Minute)
	for {
		if timeout.Before(time.Now()) {
			fmt.Println("gave up waiting for user to respond to Authy device registration request")
			return
		}

		_, _ = fmt.Printf("Checking device registration status (%s until we give up)\n", time.Until(timeout).Truncate(time.Second))
		regStatus, err := client.authyCli.CheckDeviceRegistration(context.Background(), client.cfg.UserID, regStart.RequestID)
		if err != nil {
			_, _ = fmt.Printf("Failed to check device registration status err %v\n", err)
			return
		}

		if regStatus.Status == authy.RegistrationStatusAccepted {
			regPIN = regStatus.PIN
			break
		} else if regStatus.Status != authy.RegistrationStatusPending {
			_, _ = fmt.Printf("invalid status while waiting for device registration: %s\n", regStatus.Status)
			return
		}
		time.Sleep(5 * time.Second)
	}

	// We have the registration PIN, complete the registration
	regComplete, err := client.authyCli.CompleteDeviceRegistration(context.Background(), client.cfg.UserID, regPIN)
	if err != nil {
		_, _ = fmt.Printf("Failed to complete device registration, err %v\n", err)
		return
	}

	if regComplete.Device.SecretSeed == "" {
		_, _ = fmt.Println("something went wrong completing the device registration")
		return
	}

	_, _ = fmt.Printf("Please provide your Authy TOTP backup password: ")
	pp, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Failed to read the password")
		return
	}
	fmt.Println()
	client.cfg.BackupPassword = strings.TrimSpace(string(pp))

	client.cfg.UserID = regComplete.AuthyID
	client.cfg.Seed = regComplete.Device.SecretSeed
	client.cfg.DeviceID = regComplete.Device.ID
	client.cfg.APIKey = regComplete.Device.APIKey

	respApps, err := client.authyCli.QueryAuthenticatorApps(context.Background(), client.cfg.UserID, client.cfg.DeviceID, client.cfg.Seed)
	if err != nil {
		log.Fatalf("Could not fetch authenticator apps: %v\n", err)
	}
	if !respApps.Success {
		log.Fatalf("Failed to fetch authenticator apps: %+v\n", respApps)
	}

	// Fetch the actual tokens now
	tokens, err := client.authyCli.QueryAuthenticatorTokens(context.Background(), client.cfg.UserID, client.cfg.DeviceID, client.cfg.Seed)
	if err != nil {
		log.Fatalf("Could not fetch authenticator tokens: %v\n", err)
	}
	if !tokens.Success {
		log.Fatalf("Failed to fetch authenticator tokens: %+v", tokens)
	}

	client.cfg.AuthenticatorTokenS = tokens.AuthenticatorTokens
	if err = client.cfg.Save(client.cfgFile); err != nil {
		log.Fatalf("Failed to save config: %v\n", err)
	}
}
