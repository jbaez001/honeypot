package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"honeypot/internal/config"
	"honeypot/internal/honeypot"
	"honeypot/internal/version"
)

var (
	appVersion     = fmt.Sprintf("honeypot (%s)", version.Version)
	honeypotConfig config.HoneypotConfig
	cfgFile        string
	rootCmd        = &cobra.Command{
		Use:   "honeypot",
		Short: appVersion,
		Long:  appVersion,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting honeypot (%s)\n", version.Version)
			honeypot.Start(honeypotConfig)
		},
	}

	// defaultConfig is loaded only if no config file is found
	defaultConfig = config.HoneypotConfig{
		Name:      "honeypot-1",
		ShoutUrls: nil,
		Honeypots: []config.Honeypot{
			{"FTP", "21", true, false},
			{"HTTP", "80", true, false},
			{"HTTPS", "443", true, false},
			{"TELNET", "23", true, false},
			{"VNC", "5900", true, false},
		},
	}
)

// InitConfig loads the configuration file
func InitConfig(cfgFile string, baseCfg *config.HoneypotConfig) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("honeypot")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Error reading config file", err)
		fmt.Println("Using default config")
		honeypotConfig = defaultConfig
	}

	if err := viper.Unmarshal(&baseCfg); err != nil {
		log.Fatalf("unable to decode config: %v", err)
	}
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	// detect if the user has set the flags
	nameFlag := rootCmd.PersistentFlags().Lookup("name")
	shoutUrlsFlag := rootCmd.PersistentFlags().Lookup("shout-urls")
	honeypotsFlag := rootCmd.PersistentFlags().Lookup("honeypot")

	// skip loading config if flags are set
	if nameFlag.Changed || shoutUrlsFlag.Changed || honeypotsFlag.Changed {
		honeypotStrings, _ := rootCmd.PersistentFlags().GetStringSlice("honeypot")
		honeypotConfig.Honeypots = parseHoneypots(honeypotStrings)
	} else {
		InitConfig(cfgFile, &honeypotConfig)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/honeypot.yaml)")
	rootCmd.PersistentFlags().StringVar(&honeypotConfig.Name, "name", defaultConfig.Name, "name of the honeypot")
	rootCmd.PersistentFlags().StringSliceVar(&honeypotConfig.ShoutUrls, "shout-urls", defaultConfig.ShoutUrls, "list of URLs to shout to")
	rootCmd.PersistentFlags().StringSlice("honeypot", []string{}, "list of honeypots in the format protocol:port:enabled:fragile")

	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("shout-urls", rootCmd.PersistentFlags().Lookup("shout-urls"))
}

func parseHoneypots(honeypotStrings []string) []config.Honeypot {
	var honeypots []config.Honeypot
	for _, hs := range honeypotStrings {
		parts := strings.Split(hs, ":")
		if len(parts) != 4 {
			log.Fatalf("invalid honeypot format: %s", hs)
		}
		enabled, err := strconv.ParseBool(parts[2])
		if err != nil {
			log.Fatalf("invalid enabled value: %s", parts[2])
		}
		fragile, err := strconv.ParseBool(parts[3])
		if err != nil {
			log.Fatalf("invalid fragile value: %s", parts[3])
		}
		honeypots = append(honeypots, config.Honeypot{
			Protocol: parts[0],
			Port:     parts[1],
			Enabled:  enabled,
			Fragile:  fragile,
		})
	}
	return honeypots
}

func main() {
	Execute()
}
