/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
)

var (
	h2         bool
	h3         bool
	remoteName bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gurl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		run(args...)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.Flags().BoolVar(&h2, "http2", false, "Use HTTP 2")
	rootCmd.Flags().BoolVar(&h3, "http3", false, "Use HTTP v3 only")
	rootCmd.Flags().BoolVarP(&remoteName, "remote-name", "O", false, "Write output to a file named as the remote file")
}

func run(addresses ...string) {
	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
	for _, address := range addresses {
		resp, err := client.Get(address)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		buffer := &bytes.Buffer{}
		io.Copy(buffer, resp.Body)
		if remoteName {
			strs := strings.Split(address, "/")
			fileName := strs[len(strs)-1]
			file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			_, err = io.Copy(file, buffer)
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("%s\n", buffer.Bytes())
		}
	}
}
