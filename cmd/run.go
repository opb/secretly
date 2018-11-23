package main

import (
	"log"

	"github.com/opb/secretly/kvp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var runEnvFiles bool

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&runEnvFiles, "use-files", "f", false, "Use files instead of AWS Secrets Manager")
}

var runCmd = &cobra.Command{
	Use:   "run <secretsNames...> -- <command> [<arg...>]",
	Short: "Grab secrets from AWS SecretsManager (or local files) and load into the env",
	Long: `Use secretly as a prefix for a command to load secrets from
AWS SecretsManager (or local files) into the Env variable context for 
the specified command`,
	Args: func(cmd *cobra.Command, args []string) error {
		dashIdx := cmd.ArgsLenAtDash()
		if dashIdx == -1 {
			return errors.New("please separate secretsNames and command with '--'. See usage")
		}
		if err := cobra.MinimumNArgs(1)(cmd, args[:dashIdx]); err != nil {
			return errors.Wrap(err, "at least one secretsNames must be specified")
		}
		if err := cobra.MinimumNArgs(1)(cmd, args[dashIdx:]); err != nil {
			return errors.Wrap(err, "must specify command to run. See usage")
		}
		return nil
	},
	Run: execRun,
}

func execRun(cmd *cobra.Command, args []string) {
	dashIdx := cmd.ArgsLenAtDash()
	secretsNames := args[:dashIdx]
	commandAndArgs := args[dashIdx:]
	command := commandAndArgs[0]

	var commandArgs []string
	if len(commandAndArgs) > 2 {
		commandArgs = commandAndArgs[2:]
	}

	var provider kvp.Provider
	if runEnvFiles{
		provider = kvp.FileProvider{}
	}else{
		provider = kvp.SMProvider{}
	}
	envs, err := kvp.MergedEnvPairs(secretsNames, provider)
	if err != nil {
		log.Fatalln(err)
	}

	execPlatform(command, commandArgs, envs)
}
