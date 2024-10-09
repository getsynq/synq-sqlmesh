package cmd

import (
	sqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"context"
	"fmt"
	"github.com/getsynq/synq-sqlmesh/build"
	"github.com/getsynq/synq-sqlmesh/git"
	"github.com/getsynq/synq-sqlmesh/process"
	"github.com/getsynq/synq-sqlmesh/sqlmesh"
	"github.com/getsynq/synq-sqlmesh/synq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "synq-sqlmesh",
	Short: "Small utility to collect SQLMesh metadata information and upload it to Synq",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of synq-sqlmesh",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("synq-sqlmesh %s (%s)", strings.TrimSpace(build.Version), strings.TrimSpace(build.Time))
	},
}

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect metadata information from SQLMesh and store to the file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		gitContext := git.CollectGitContext(cmd.Context(), SQLMeshProjectDir)

		err := WithSQLMesh(func(baseUrl url.URL) error {
			logrus.Info("SQLMesh base URL:", baseUrl.String())

			output, err := sqlmesh.CollectMetadata(baseUrl, createFileContentGlobFilter())
			if err != nil {
				return err
			}
			output.GitContext = gitContext

			if err := synq.DumpMetadata(output, args[0]); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Collect metadata information from SQLMesh and send to Synq API",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		gitContext := git.CollectGitContext(cmd.Context(), SQLMeshProjectDir)

		err := WithSQLMesh(func(baseUrl url.URL) error {
			logrus.Info("SQLMesh base URL:", baseUrl.String())

			output, err := sqlmesh.CollectMetadata(baseUrl, createFileContentGlobFilter())
			if err != nil {
				return err
			}
			output.GitContext = gitContext

			if SynqApiToken == "" {
				return fmt.Errorf("SYNQ_TOKEN environment variable is not set")
			}

			if err := synq.UploadMetadata(cmd.Context(), output, SynqApiEndpoint, SynqApiToken); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	},
}

var uploadAuditCmd = &cobra.Command{
	Use:   "upload_audit",
	Short: "Sends to Synq output of `audit` command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		gitContext := git.CollectGitContext(cmd.Context(), SQLMeshProjectDir)
		output := &sqlmeshv1.IngestExecutionRequest{
			Command:    []string{"sqlmesh", "audit"},
			GitContext: gitContext,
		}
		output.UploaderVersion = strings.TrimSpace(fmt.Sprintf("synq-sqlmesh/%s", build.Version))
		output.UploaderBuildTime = strings.TrimSpace(build.Time)

		for _, fileArg := range args {
			err := sqlmesh.CollectAuditLog(output, fileArg)
			if err != nil {
				logrus.WithError(err).Error("Failed to collect audit log")
				os.Exit(0)
			}
		}

		if SynqApiToken == "" {
			logrus.Error("SYNQ_TOKEN environment variable is not set")
			os.Exit(0)
		}

		if err := synq.UploadExecutionLog(cmd.Context(), output, SynqApiEndpoint, SynqApiToken); err != nil {
			logrus.WithError(err).Error("Failed to upload execution log")
			os.Exit(0)
		}
	},
}

var uploadRunCmd = &cobra.Command{
	Use:   "upload_run",
	Short: "Sends to Synq output of `run` command",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		gitContext := git.CollectGitContext(cmd.Context(), SQLMeshProjectDir)
		output := &sqlmeshv1.IngestExecutionRequest{
			Command:    []string{"sqlmesh", "run"},
			GitContext: gitContext,
		}
		output.UploaderVersion = strings.TrimSpace(fmt.Sprintf("synq-sqlmesh/%s", build.Version))
		output.UploaderBuildTime = strings.TrimSpace(build.Time)

		for _, fileArg := range args {
			err := sqlmesh.CollectAuditLog(output, fileArg)
			if err != nil {
				logrus.WithError(err).Error("Failed to collect audit log")
				os.Exit(0)
			}
		}

		if SynqApiToken == "" {
			logrus.Error("SYNQ_TOKEN environment variable is not set")
			os.Exit(0)
		}

		if err := synq.UploadExecutionLog(cmd.Context(), output, SynqApiEndpoint, SynqApiToken); err != nil {
			logrus.WithError(err).Error("Failed to upload execution log")
			os.Exit(0)
		}
	},
}

func createFileContentGlobFilter() sqlmesh.GlobFilter {
	if SQLMeshCollectFileContent {
		return sqlmesh.NewGlobFilter(SQLMeshCollectFileContentIncludePattern, SQLMeshCollectFileContentExcludePattern)
	}
	return sqlmesh.NewExcludeEverythingGlobFilter()
}

func WithSQLMesh(f func(baseUrl url.URL) error) error {
	baseUrl := url.URL{
		Host:   fmt.Sprintf("%s:%d", SQLMeshUiHost, SQLMeshUiPort),
		Scheme: "http",
	}
	if SQLMeshUiStart {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()
		sqlMeshProcess, err := process.ExecuteCommand(ctx, SQLMesh, []string{"ui", "--host", SQLMeshUiHost, "--port", fmt.Sprintf("%d", SQLMeshUiPort)}, process.WithDir(SQLMeshProjectDir))
		if err != nil {
			return err
		}

		sqlmesh.WaitForSQLMeshToStart(baseUrl)

		err = f(baseUrl)
		_ = sqlMeshProcess.Kill()
		if err != nil {
			return err
		}
		return nil
	} else {
		return f(baseUrl)
	}
}

var SynqApiEndpoint string = "https://developer.synq.io/"
var SynqApiToken string = os.Getenv("SYNQ_TOKEN")
var SQLMesh string = "sqlmesh"
var SQLMeshProjectDir string = "."
var SQLMeshUiStart bool = true
var SQLMeshUiHost string = "localhost"
var SQLMeshUiPort int = 8080
var SQLMeshCollectFileContent = false
var SQLMeshCollectFileContentIncludePattern = "external_models.yaml,models/**.sql,models/**.py,audits/**.sql,tests/**.yaml"
var SQLMeshCollectFileContentExcludePattern = "*.log"

func init() {
	rootCmd.PersistentFlags().StringVar(&SynqApiToken, "synq-token", SynqApiToken, "Synq API token")
	rootCmd.PersistentFlags().StringVar(&SynqApiEndpoint, "synq-endpoint", SynqApiEndpoint, "Synq API endpoint URL")
	rootCmd.PersistentFlags().StringVar(&SQLMesh, "sqlmesh-cmd", SQLMesh, "SQLMesh launcher location")
	rootCmd.PersistentFlags().StringVar(&SQLMeshProjectDir, "sqlmesh-project-dir", SQLMeshProjectDir, "Location of SQLMesh project directory")
	rootCmd.PersistentFlags().BoolVar(&SQLMeshUiStart, "sqlmesh-ui-start", SQLMeshUiStart, "Launch and control SQLMesh UI process automatically")
	rootCmd.PersistentFlags().StringVar(&SQLMeshUiHost, "sqlmesh-ui-host", SQLMeshUiHost, "SQLMesh UI host")
	rootCmd.PersistentFlags().IntVar(&SQLMeshUiPort, "sqlmesh-ui-port", SQLMeshUiPort, "SQLMesh UI port")
	rootCmd.PersistentFlags().BoolVar(&SQLMeshCollectFileContent, "sqlmesh-collect-file-content", SQLMeshCollectFileContent, "If content of the project files should be collected")
	rootCmd.PersistentFlags().StringVar(&SQLMeshCollectFileContentIncludePattern, "sqlmesh-collect-file-content-include", SQLMeshCollectFileContentIncludePattern, "File patterns to include content")
	rootCmd.PersistentFlags().StringVar(&SQLMeshCollectFileContentExcludePattern, "sqlmesh-collect-file-content-exclude", SQLMeshCollectFileContentExcludePattern, "File patterns to exclude content")

	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(uploadAuditCmd)
	rootCmd.AddCommand(uploadRunCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
