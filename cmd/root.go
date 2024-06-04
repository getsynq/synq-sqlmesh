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
	Short: "Small utility to collect SqlMesh metadata information and upload it to Synq",
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
	Short: "Collect metadata information from SqlMesh and store to the file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := WithSqlMesh(func(baseUrl url.URL) error {
			logrus.Info("SqlMesh base URL:", baseUrl.String())

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
	Short: "Collect metadata information from SqlMesh and send to Synq API",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		err := WithSqlMesh(func(baseUrl url.URL) error {
			logrus.Info("SqlMesh base URL:", baseUrl.String())

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

func WithSqlMesh(f func(baseUrl url.URL) error) error {
	baseUrl := url.URL{
		Host:   fmt.Sprintf("%s:%d", SqlMeshUiHost, SqlMeshUiPort),
		Scheme: "http",
	}
	if SqlMeshUiStart {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()
		sqlMeshProcess, err := process.ExecuteCommand(ctx, SqlMesh, []string{"ui", "--host", SqlMeshUiHost, "--port", fmt.Sprintf("%d", SqlMeshUiPort)}, process.WithDir(SqlMeshProjectDir))
		if err != nil {
			return err
		}

		sqlmesh.WaitForSqlMeshToStart(baseUrl)

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
var SqlMesh string = "sqlmesh"
var SqlMeshProjectDir string = "."
var SqlMeshUiStart bool = true
var SqlMeshUiHost string = "localhost"
var SqlMeshUiPort int = 8080

func init() {
	rootCmd.PersistentFlags().StringVar(&SynqApiToken, "synq-token", SynqApiToken, "Synq API token")
	rootCmd.PersistentFlags().StringVar(&SynqApiEndpoint, "synq-endpoint", SynqApiEndpoint, "Synq API endpoint URL")
	rootCmd.PersistentFlags().StringVar(&SqlMesh, "sqlmesh-cmd", SqlMesh, "SqlMesh launcher location")
	rootCmd.PersistentFlags().StringVar(&SqlMeshProjectDir, "sqlmesh-project-dir", SqlMeshProjectDir, "Location of SqlMesh project directory")
	rootCmd.PersistentFlags().BoolVar(&SqlMeshUiStart, "sqlmesh-ui-start", SqlMeshUiStart, "Launch and control SqlMesh UI process automatically")
	rootCmd.PersistentFlags().StringVar(&SqlMeshUiHost, "sqlmesh-ui-host", SqlMeshUiHost, "SqlMesh UI host")
	rootCmd.PersistentFlags().IntVar(&SqlMeshUiPort, "sqlmesh-ui-port", SqlMeshUiPort, "SqlMesh UI port")

	rootCmd.AddCommand(collectCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(uploadCmd)

}

func Execute() error {
	return rootCmd.Execute()
}
