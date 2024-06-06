package cmd

import (
	"context"
	"fmt"
	"github.com/getsynq/synq-sqlmesh/build"
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
		err := WithSQLMesh(func(baseUrl url.URL) error {
			logrus.Info("SQLMesh base URL:", baseUrl.String())

			output, err := sqlmesh.CollectMetadata(baseUrl)
			if err != nil {
				return err
			}

			if err := synq.DumpMetadata(output, args[0]); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Collect metadata information from SQLMesh and send to Synq API",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := WithSQLMesh(func(baseUrl url.URL) error {
			logrus.Info("SQLMesh base URL:", baseUrl.String())

			output, err := sqlmesh.CollectMetadata(baseUrl)
			if err != nil {
				return err
			}

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
			os.Exit(1)
		}
	},
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

func init() {
	rootCmd.PersistentFlags().StringVar(&SynqApiToken, "synq-token", SynqApiToken, "Synq API token")
	rootCmd.PersistentFlags().StringVar(&SynqApiEndpoint, "synq-endpoint", SynqApiEndpoint, "Synq API endpoint URL")
	rootCmd.PersistentFlags().StringVar(&SQLMesh, "sqlmesh-cmd", SQLMesh, "SQLMesh launcher location")
	rootCmd.PersistentFlags().StringVar(&SQLMeshProjectDir, "sqlmesh-project-dir", SQLMeshProjectDir, "Location of SQLMesh project directory")
	rootCmd.PersistentFlags().BoolVar(&SQLMeshUiStart, "sqlmesh-ui-start", SQLMeshUiStart, "Launch and control SQLMesh UI process automatically")
	rootCmd.PersistentFlags().StringVar(&SQLMeshUiHost, "sqlmesh-ui-host", SQLMeshUiHost, "SQLMesh UI host")
	rootCmd.PersistentFlags().IntVar(&SQLMeshUiPort, "sqlmesh-ui-port", SQLMeshUiPort, "SQLMesh UI port")

	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(uploadCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
