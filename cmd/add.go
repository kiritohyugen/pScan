/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"kiritohyugen/cobra/pScan/scan"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1>...<hostN>",
	Short:        "Add new host(s) to the list",
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// hostsFile, err := cmd.Flags().GetString("hosts-file")
		// if err != nil {
		// 	return err
		// }
		hostsFile := viper.GetString("hosts-file")

		return addAction(os.Stdout, hostsFile, args)
	},
}

func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addAction(out io.Writer, hostsFile string, args []string) error {

	hl := &scan.HostsList{}
	// fmt.Println("Hosts file path:", hostsFile) // Debugging output

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	for _, h := range args {
		if err := hl.Add(h); err != nil {
			return err
		}

		fmt.Fprintln(out, "Added host:", h)
	}

	return hl.Save(hostsFile)
}
