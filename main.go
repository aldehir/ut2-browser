package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var cmd = &cobra.Command{
	Use:     "ut2-browser",
	PreRunE: preRun,
	RunE:    runServer,
}

var configFile string
var httpAddr string
var httpRoot string

var config Config
var registry *Registry
var state *State
var engine *QueryEngine

var logger *slog.Logger

//go:embed web/dist
var embeddedFrontend embed.FS

func init() {
	cmd.Flags().StringVarP(&configFile, "config", "c", "config.yml", "config file")
	cmd.Flags().StringVar(&httpAddr, "http", ":8200", "http listen address")
	cmd.Flags().StringVar(&httpRoot, "root", "", "http root directory, use embedded if empty")
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func preRun(cmd *cobra.Command, args []string) error {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger = slog.New(handler)

	err := loadConfig()
	if err != nil {
		return err
	}

	registry = &Registry{}
	state = NewState()
	engine = &QueryEngine{
		Registry: registry,
		State:    state,
	}

	return nil
}

func runServer(cmd *cobra.Command, args []string) error {
	addStaticServers()

	var webFS fs.FS
	webFS, _ = fs.Sub(embeddedFrontend, "web/dist")
	if httpRoot != "" {
		webFS = os.DirFS(httpRoot)
	}

	webServer := NewBrowserServer(webFS)
	webServer.Addr = httpAddr

	go func() {
		if err := webServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error("http server error", "err", err)
		}
	}()

	if err := engine.Run(); err != ErrStopped {
		return err
	}

	return nil
}

func loadConfig() error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	return decoder.Decode(&config)
}

func addStaticServers() {
	for _, static := range config.Static {
		for _, server := range static.Servers {
			registry.Register(Registration{
				Address:      server.Address,
				Group:        static.Group,
				Interval:     static.Interval,
				Timeout:      server.Timeout,
				Persist:      true,
				ExternalLink: server.ExternalLink,
			})
		}
	}
}
