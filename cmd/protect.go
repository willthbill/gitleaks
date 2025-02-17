package cmd

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/zricethezav/gitleaks/v8/logging"
	"github.com/zricethezav/gitleaks/v8/report"
	"github.com/zricethezav/gitleaks/v8/sources"
)

func init() {
	protectCmd.Flags().Bool("staged", false, "detect secrets in a --staged state")
	protectCmd.Flags().String("log-opts", "", "git log options")
	protectCmd.Flags().StringP("source", "s", ".", "path to source")
	rootCmd.AddCommand(protectCmd)
}

var protectCmd = &cobra.Command{
	Use:    "protect",
	Short:  "protect secrets in code",
	Run:    runProtect,
	Hidden: true,
}

func runProtect(cmd *cobra.Command, args []string) {
	source, err := cmd.Flags().GetString("source")
	if err != nil {
		logging.Fatal().Err(err).Msg("could not get source")
	}
	initConfig(source)

	// setup config (aka, the thing that defines rules)
	cfg := Config(cmd)

	exitCode, _ := cmd.Flags().GetInt("exit-code")
	staged, _ := cmd.Flags().GetBool("staged")
	start := time.Now()
	detector := Detector(cmd, cfg, source)

	// start git scan
	var findings []report.Finding
	gitCmd, err := sources.NewGitDiffCmd(source, staged)
	if err != nil {
		logging.Fatal().Err(err).Msg("")
	}
	findings, err = detector.DetectGit(gitCmd)

	findingSummaryAndExit(detector, findings, exitCode, start, err)
}
