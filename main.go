package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"syscall"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/codingdestro/wallet-go/pkg/utils"
	"github.com/codingdestro/wallet-go/ui"
	"golang.org/x/term"
)

func main() {
	listFlag := flag.Bool("l", false, "List all keys in the vault")
	addFlag := flag.String("a", "", "Add or update a secret: -a <key>")
	delFlag := flag.String("d", "", "Delete a secret: -d <key>")
	getFlag := flag.String("v", "", "View and copy a secret: -v <key>")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  (no flags)  Launch interactive TUI\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *listFlag {
		handleList()
		return
	}
	if *addFlag != "" {
		handleAdd(*addFlag)
		return
	}
	if *delFlag != "" {
		handleDelete(*delFlag)
		return
	}
	if *getFlag != "" {
		handleGet(*getFlag)
		return
	}

	p := tea.NewProgram(ui.InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
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

func loadVault(password string) map[string]string {
	var data map[string]string
	err := utils.LoadEncryptedJSON("wallet.enc", password, &data)
	if err != nil {
		fmt.Println("Error: Access Denied. Invalid Password.")
		os.Exit(1)
	}
	return data
}

func handleList() {
	if !utils.FileExists("wallet.enc") {
		fmt.Println("Error: vault.enc not found.")
		os.Exit(1)
	}
	password := promptPassword()
	data := loadVault(password)

	if len(data) == 0 {
		fmt.Println("Vault is empty.")
		return
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Println("\nVault Secrets:")
	for _, k := range keys {
		fmt.Printf("• %s\n", k)
	}
}

func handleAdd(key string) {
	password := promptPassword()
	var data map[string]string
	if utils.FileExists("wallet.enc") {
		data = loadVault(password)
	} else {
		data = make(map[string]string)
	}

	fmt.Printf("Enter Value for [%s]: ", key)
	byteVal, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Printf("Error reading value: %v\n", err)
		os.Exit(1)
	}

	data[strings.TrimSpace(key)] = string(byteVal)
	err = utils.SaveEncryptedJSON("wallet.enc", password, data)
	if err != nil {
		fmt.Printf("Error saving vault: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully added/updated: %s\n", key)
}

func handleDelete(key string) {
	if !utils.FileExists("wallet.enc") {
		fmt.Println("Error: vault.enc not found.")
		os.Exit(1)
	}
	password := promptPassword()
	data := loadVault(password)

	if _, ok := data[key]; !ok {
		fmt.Printf("Error: Key [%s] not found.\n", key)
		os.Exit(1)
	}

	delete(data, key)
	err := utils.SaveEncryptedJSON("wallet.enc", password, data)
	if err != nil {
		fmt.Printf("Error saving vault: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully deleted: %s\n", key)
}

func handleGet(key string) {
	if !utils.FileExists("wallet.enc") {
		fmt.Println("Error: vault.enc not found.")
		os.Exit(1)
	}
	password := promptPassword()
	data := loadVault(password)

	val, ok := data[key]
	if !ok {
		fmt.Printf("Error: Key [%s] not found.\n", key)
		os.Exit(1)
	}

	err := clipboard.WriteAll(val)
	if err != nil {
		fmt.Printf("Error: Value retrieved but failed to copy to clipboard: %v\n", err)
		fmt.Println("Value (visible for 5s):", val)
		return
	}
	fmt.Printf("Success: Value for [%s] copied to clipboard!\n", key)
}
