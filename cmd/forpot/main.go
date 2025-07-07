package main

import (
	"forpot/internal/app"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	var port int
	cmd := &cobra.Command{
		Use:   "forpot [user@]host",
		Short: "Port forwarding app",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			app.RunPortForwarding(args[0], port)
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 22, "Set ssh port")
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
