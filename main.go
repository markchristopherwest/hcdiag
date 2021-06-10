package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/host-diagnostics/agent"
)

func main() {
	os.Exit(realMain())
}

func realMain() (returnCode int) {
	// TODO(mkcp): rename to support-bundler
	l := configureLogging("host-diagnostics")

	// Parse our CLI flags
	flags := Flags{}
	err := flags.parseFlags(os.Args[1:])
	if err != nil {
		return 64
	}

	// FIXME(mkcp): pipe our configuration data into the agent
	var config agent.Config
	if flags.Config != "" {
		config, err := agent.ParseHCL(flags.Config)
		if err != nil {
			log.Fatalf("Failed to load configuration: %s", err)
		}
		l.Debug("Config is", "config", config)
	}

	cfg := mergeAgentConfig(config, flags)
	a := agent.NewAgent(cfg, l)

	// Run the agent
	// NOTE(mkcp): Are there semantic returnCodes we can send based on the agent error type?
	errs := a.Run()
	if 0 < len(errs) {
		return 1
	}
	return 0
}

// configureLogging takes a logger name, sets the default configuration, grabs the LOG_LEVEL from our ENV vars, and
//  returns a configured and usable logger.
func configureLogging(loggerName string) hclog.Logger {
	// Create logger, set default and log level
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name: loggerName,
	})
	hclog.SetDefault(appLogger)
	if logStr := os.Getenv("LOG_LEVEL"); logStr != "" {
		if level := hclog.LevelFromString(logStr); level != hclog.NoLevel {
			appLogger.SetLevel(level)
			appLogger.Debug("Logger configuration change", "LOG_LEVEL", hclog.Fmt("%s", logStr))
		}
	}
	return hclog.Default()
}

// Flags stores our CLI inputs.
type Flags struct {
	OS           string
	Serial       bool
	Dryrun       bool
	Consul       bool
	Nomad        bool
	TFE          bool
	Vault        bool
	AllProducts  bool
	Includes     []string
	IncludeSince time.Duration
	Outfile      string
	Config       string
}

type CSVFlag struct {
	Values *[]string
}

func (s CSVFlag) String() string {
	if s.Values == nil {
		return ""
	}
	return strings.Join(*s.Values, ",")
}

func (s CSVFlag) Set(v string) error {
	*s.Values = strings.Split(v, ",")
	return nil
}

func (f *Flags) parseFlags(args []string) error {
	flags := flag.NewFlagSet("hc-bundler", flag.ExitOnError)
	flags.BoolVar(&f.Dryrun, "dryrun", false, "Performing a dry run will display all commands without executing them")
	flags.BoolVar(&f.Serial, "serial", false, "Run products in sequence rather than concurrently")
	flags.BoolVar(&f.Consul, "consul", false, "Run Consul diagnostics")
	flags.BoolVar(&f.Nomad, "nomad", false, "Run Nomad diagnostics")
	flags.BoolVar(&f.TFE, "tfe", false, "Run Terraform Enterprise diagnostics")
	flags.BoolVar(&f.Vault, "vault", false, "Run Vault diagnostics")
	flags.BoolVar(&f.AllProducts, "all", false, "Run all available product diagnostics")
	flags.StringVar(&f.OS, "os", "auto", "Override operating system detection")
	flags.StringVar(&f.Outfile, "outfile", "support", "Output file name")
	flags.StringVar(&f.Config, "config", "", "Path to HCL configuration file")
	flags.DurationVar(&f.IncludeSince, "include-since", time.Duration(0), "How long ago until now to include files. Examples: 72h, 25m, 45s, 120h1m90s")
	flags.Var(&CSVFlag{&f.Includes}, "includes", "files or directories to include (comma-separated, file-*-globbing available if 'wrapped-*-in-single-quotes')\ne.g. '/var/log/consul-*,/var/log/nomad-*'")

	return flags.Parse(args)
}

// FIXME(mkcp): Don't love how this fits together yet
// mergeAgentConfig merges flags into the agent.Config, prioritizing flags over HCL config.
func mergeAgentConfig(config agent.Config, flags Flags) agent.Config {
	// Convert our flag input to agent configuration
	from := time.Unix(0, flags.IncludeSince.Nanoseconds())
	to := time.Now()
	config.OS = flags.OS
	config.Serial = flags.Serial
	config.Dryrun = flags.Dryrun
	config.Consul = flags.AllProducts || flags.Consul
	config.Nomad = flags.AllProducts || flags.Nomad
	config.TFE = flags.AllProducts || flags.TFE
	config.Vault = flags.AllProducts || flags.Vault
	config.Includes = flags.Includes
	config.IncludeFrom = from
	config.IncludeTo = to
	config.Outfile = flags.Outfile
	return config
}
