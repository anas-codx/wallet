package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/codingdestro/wallet-go/internal/vault"
	"github.com/codingdestro/wallet-go/ui"
	"golang.org/x/term"
)

const vaultPath = "wallet.enc"

func main() {
	listFlag := flag.Bool("l", false, "List all keys in the vault")
	addFlag := flag.String("a", "", "Add or update a secret: -a <key>")
	delFlag := flag.String("d", "", "Delete a secret: -d <key>")
	getFlag := flag.String("v", "", "View and copy a secret: -v <key>")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of wallet:\n")
		fmt.Fprintf(os.Stderr, "  (no flags)  Launch interactive TUI\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// If any flag is set, use the CLI handler
	if *listFlag || *addFlag != "" || *delFlag != "" || *getFlag != "" {
		handleCLI(*listFlag, *addFlag, *delFlag, *getFlag)
		return
	}

	// Default: Launch TUI
	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}

func handleCLI(list bool, add, del, get string) {
	password := promptPassword()
	v, err := vault.Open(vaultPath, password)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if list {
		keys := v.List()
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
		if err := v.Set(add, string(byteVal)); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully added/updated: %s\n", add)
	}

	if del != "" {
		if err := v.Delete(del); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Successfully deleted: %s\n", del)
	}

	if get != "" {
		val, ok := v.Get(get)
		if !ok {
			fmt.Printf("Error: Key [%s] not found.\n", get)
			os.Exit(1)
		}
		if err := clipboard.WriteAll(val); err != nil {
			fmt.Printf("Error: Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Success: Value for [%s] copied to clipboard!\n", get)
	}
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
