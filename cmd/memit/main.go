package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/liamg/memit"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "memit [url] -- [args-for-binary...]",
	Short: "Run a binary from a URL without writing it to disk",
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {

		fmt.Println("Downloading file...")
		url := args[0]
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to retrieve binary: %s\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		fmt.Println("Configuring...")
		cmd, f, err := memit.Command(resp.Body, args[1:]...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create command: %s\n", err)
			os.Exit(1)
		}
		defer f.Close()

		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout

		fmt.Println("Running...")
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute binary: %s\n", err)
			os.Exit(1)
		}
	},
}

func main() {
	rootCmd.Execute()
}
