package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/shivylp/radium"
	"github.com/spf13/cobra"
)

func newServeCmd(cfg *config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start radium in server mode with local configuration",
		Long: `
Start radium in server mode with local configuration.

In this mode, radium runs as a server and other radium CLIs can
connect to it for queries/data. This mode can be used to host a
central source of truth for a team or a community. Pass '--clipboard'
or '-C' flag to enable clipboard monitoring.

You can also use this to start radium in clipboard-monitor-only mode
by passing '--clipboard' or '-C' option and setting '--addr' blank.
`,
	}

	var clipboard bool
	var addr string
	cmd.Flags().BoolVarP(&clipboard, "clipboard", "C", false, "Enable clipboard monitoring")
	cmd.Flags().StringVarP(&addr, "addr", "a", ":8080", "Listen on <ip>:<port>")

	cmd.Run = func(_ *cobra.Command, _ []string) {
		addr = strings.TrimSpace(addr)

		if addr == "" && !clipboard {
			fmt.Printf("both api and clipboard modes are disabled. nothing to do")
			os.Exit(1)
		}

		var wg sync.WaitGroup

		if addr != "" {
			ins := getNewRadiumInstance()
			srv := radium.NewServer(ins)
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				ins.Infof("starting server on '%s'...", addr)
				ins.Errorf("Err: %s", http.ListenAndServe(addr, srv))
				wg.Done()
			}(&wg)
		}

		if clipboard {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				// TODO: run clipboard monitor here
				wg.Done()
			}(&wg)
		}

		wg.Wait()
	}
	return cmd
}
