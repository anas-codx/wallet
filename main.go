package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/codingdestro/wallet-go/internal/platform"
	"github.com/codingdestro/wallet-go/internal/vault"
	"github.com/codingdestro/wallet-go/pkg/utils"
	"github.com/codingdestro/wallet-go/ui"
	"golang.org/x/term"
)

func getVaultPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		// Fallback to home directory if config dir is not available
		homeDir, _ := os.UserHomeDir()
		configDir = homeDir
	}

	appDir := filepath.Join(configDir, "wallet")
	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		_ = os.MkdirAll(appDir, 0700)
	}

	return filepath.Join(appDir, "wallet.enc")
}

func main() {
	vaultPath := getVaultPath()
	listFlag := flag.Bool("l", false, "List all keys in the vault")
	addFlag := flag.String("a", "", "Add or update a secret: -a <key>")
	delFlag := flag.String("d", "", "Delete a secret: -d <key>")
	getFlag := flag.String("v", "", "View and copy a secret: -v <key>")
	genFlag := flag.Int("g", 0, "Generate a secure password of length <n>")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of wallet:\n")
		fmt.Fprintf(os.Stderr, "  (no flags)  Launch interactive TUI\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *genFlag > 0 {
		handleGenerate(*genFlag)
		return
	}

	// Dependency Injection Setup
	cb := &platform.SystemClipboard{}
	factory := func(password string) (vault.SecretStore, error) {
		return vault.NewStore(vaultPath, password)
	}

	// If any flag is set, use the CLI handler
	if *listFlag || *addFlag != "" || *delFlag != "" || *getFlag != "" {
		handleCLI(*listFlag, *addFlag, *delFlag, *getFlag, factory, cb)
		return
	}

	// Default: Launch TUI with Dependencies Injected
	p := tea.NewProgram(ui.NewModel(factory, cb))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}

func handleCLI(list bool, add, del, get string, factory ui.StoreFactory, cb platform.Clipboard) {
	password := promptPassword()
	store, err := factory(password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if list {
		keys := store.List()
		if len(keys) == 0 {
			fmt.Println("Vault is empty.")
			return
		}
		fmt.Println("\nVault Secrets:")
		for _, k := range keys {
			fmt.Printf("• %s\n", k)
		}
	}

	if add != "" {
		fmt.Printf("Enter Value for [%s]: ", add)
		byteVal, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			fmt.Printf("Error reading value: %v\n", err)
			os.Exit(1)
		}
		if err := store.Set(add, string(byteVal)); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully added/updated: %s\n", add)
	}

	if del != "" {
		if err := store.Delete(del); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully deleted: %s\n", del)
	}

	if get != "" {
		val, ok := store.Get(get)
		if !ok {
			fmt.Printf("Error: Key [%s] not found.\n", get)
			os.Exit(1)
		}
		if err := cb.Write(val); err != nil {
			fmt.Printf("Error: Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Success: Value for [%s] copied to clipboard!\n", get)
	}
}

func handleGenerate(length int) {
	cb := &platform.SystemClipboard{}
	config := utils.PasswordConfig{
		Length:           length,
		IncludeDigits:    true,
		IncludeSymbols:   true,
		IncludeUppercase: true,
	}
	pwd, err := utils.GeneratePassword(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := cb.Write(pwd); err != nil {
		fmt.Printf("Error: Failed to copy to clipboard: %v\n", err)
		fmt.Println("Generated Password:", pwd)
		return
	}
	fmt.Printf("Success: Generated %d-character password and copied to clipboard!\n", length)
}

func promptPassword() string {
	fmt.Print("Enter Vault Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}
	return string(bytePassword)
}
