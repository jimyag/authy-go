# authy-go

[![GoDoc](https://godoc.org/github.com/jimyag/authy-go?status.svg)](https://godoc.org/github.com/jimyag/authy-go)

This is a Go library that allows you to access your [Authy](https://authy.com) TOTP tokens.

It was created to facilitate exports of your TOTP database, because Authy does not provide any way to access or port your TOTP tokens to another client.

It also somewhat documents Authy's protocol/encryption, since public materials on that are somewhat scarce.

Please be careful. You can get your Authy account suspended very easily by using this package. It does not hide itself or mimic the official clients.

This tool draws inspiration from [alexzorin/authy](https://github.com/alexzorin/authy), combining their strengths and introducing additional enhancements.

## authy-cli

This program will enroll itself as an additional device on your Authy account and export all of your TOTP tokens in [Key URI Format](https://github.com/google/google-authenticator/wiki/Key-Uri-Format).

It is also able to save the TOTP database in a JSON file encrypted with your Authy backup password, which can be used for backup purposes, and to read it back in order to decrypt it.

### Installation

Pre-built binaries are available from the [releases page](https://github.com/jimyag/authy-go/releases).

Alternatively, it can be compiled from source, which requires [Go 1.18 or newer](https://golang.org/doc/install):

```shell
go install github.com/jimyag/authy-go/app/authy-cli@latest
```

### Usage

To use authy-cli to fetch TOTP tokens and export them in Google Authenticator Key URI format:

1. Run authy-cli:

```bash
authy-cli
```

2. Provide your Authy ID and Backup Password:
The program will prompt you for:

- Your Authy ID: This is the unique identifier for your Authy account.
- Your Authy Backup Password: This is required to decrypt your TOTP secrets.

3. Device Registration:

- The program will send a device registration request using the push method.
- This will send a push notification to your existing Authy apps (Android, iOS, Desktop, or Chrome).
- You will need to approve the request from one of your other Authy apps.

4. Save Authentication Credential:

- If the registration is successful, the program will save its authentication credential (a random value) to $HOME/.config/authy-cli/authy-go.json for future use.
- Make sure to delete this file and de-register the device after you're finished.

5. Fetch and Decrypt TOTP Database:

- If the program successfully fetches your encrypted TOTP database, it will prompt you for your Authy backup password.
This password is required to decrypt the TOTP secrets.

6. Export TOTP Tokens:

- The program will dump all of your TOTP tokens in URI format, which you can use to import into other applications such as Google Authenticator.
Example output:

```bash
otpauth://totp/Issuer:AccountName?secret=BASE32SECRET&issuer=Issuer
```

### commands

- sync : Sync your Authy ID and Backup Password, and request device registration.
- export : Export all TOTP tokens in URI format
- list : List all TOTP tokens
- get : Get TOTP token by UniqueID

### Important Notes

1. Security: Always keep your Authy ID and Backup Password secure. Do not share these credentials.
2. Cleanup: Remember to delete $HOME/.config/authy-cli/authy-go.json and de-register the device from your Authy account after you have finished using this tool to avoid any security risks.
