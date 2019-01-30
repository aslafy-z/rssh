package server

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"os/user"

	"github.com/rs/zerolog/log"

	api "github.com/Xide/rssh/pkg/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type apiFlags struct {
	BindAddr      string
	BindPort      uint16
	EtcdEndpoints []string
	Config        string
	RootDomain    string
}

// Splits {"a,b", "c"} into {"a", "b", "c"}
// Temporary fix (hopefully) because Cobra doesn't
// handle separators if they are not followed by a whitespace.
func splitParts(maybeParted []string) []string {
	r := []string{}
	for _, x := range maybeParted {
		if strings.Contains(x, ",") {
			for _, newKey := range strings.Split(x, ",") {
				r = append(r, newKey)
			}
		} else {
			r = append(r, x)
		}
	}
	return r
}

// Taken from https://www.socketloop.com/tutorials/golang-use-regular-expression-to-validate-domain-name
func isValidDomain(d string) bool {
	re := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	return re.MatchString(d)
}

func parseArgs(flags *apiFlags) func() {
	return func() {

		flags.BindAddr = viper.GetString("addr")
		flags.RootDomain = viper.GetString("domain")
		if !isValidDomain(flags.RootDomain) {
			log.Fatal().Str("domain", flags.RootDomain).Msg("Invalid domain name.")
		}
		port, err := strconv.ParseUint(viper.GetString("port"), 10, 16)
		if err != nil {
			log.Fatal().
				Str("port", viper.Get("addr").(string)).
				Msg(fmt.Sprintf("Could not parse %s as an integer.", viper.Get("addr").(string)))
		}
		flags.EtcdEndpoints = splitParts(viper.GetStringSlice("etcd"))
		flags.BindPort = uint16(port)
	}
}

func initConfig(flags *apiFlags) func() {
	return func() {
		cnf := viper.GetString("config")
		if cnf != "" {
			viper.SetConfigFile(cnf)
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				log.Warn().
					Str("error", err.Error()).
					Msg("Ignoring current directory as config file source.")
			} else {
				viper.AddConfigPath(cwd)
			}

			user, err := user.Current()
			if err != nil {
				log.Warn().
					Str("error", err.Error()).
					Msg("Could not find current user informations, ignoring configuration file")
				return
			}
			viper.AddConfigPath(user.HomeDir)
			viper.SetConfigName(".rssh")
		}

		if err := viper.ReadInConfig(); err == nil {
			log.Info().Str("file", viper.ConfigFileUsed()).Msg("Configuration file loaded")
		} else {
			log.Warn().Str("error", err.Error()).Msg("Could not load configuration file.")
		}
	}
}

func NewCommand() *cobra.Command {
	flags := &apiFlags{}

	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run the RSSH public server.",
		Long:  `Run the RSSH public server.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			initConfig(flags)()
			parseArgs(flags)()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			httpAPI, err := api.NewDispatcher(
				flags.BindAddr,
				flags.BindPort,
				flags.RootDomain,
			)
			if err != nil {
				return err
			}

			executor, err := api.NewExecutor(flags.EtcdEndpoints)
			if err != nil {
				return err
			}
			return httpAPI.Run(executor)
		},
	}

	cmd.PersistentFlags().StringVarP(
		&flags.BindAddr,
		"addr",
		"a",
		"0.0.0.0",
		"HTTP API bind address",
	)

	cmd.PersistentFlags().StringVarP(
		&flags.RootDomain,
		"domain",
		"d",
		"",
		"Domain the RSSH public server will be known as.",
	)

	cmd.PersistentFlags().Uint16VarP(
		&flags.BindPort,
		"port",
		"p",
		8080,
		"HTTP API port",
	)

	cmd.PersistentFlags().StringSliceVarP(
		&flags.EtcdEndpoints,
		"etcd",
		"e",
		[]string{"http://127.0.0.1:2379"},
		"Comma separated list of the Etcd hosts to discover",
	)

	cmd.PersistentFlags().StringVarP(
		&flags.Config,
		"config",
		"c",
		"",
		"Server configuration file to use",
	)

	viper.BindPFlags(cmd.PersistentFlags())
	return cmd
}
