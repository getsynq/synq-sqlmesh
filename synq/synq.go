package synq

import (
	ingestsqlmeshv1grpc "buf.build/gen/go/getsynq/api/grpc/go/synq/ingest/sqlmesh/v1/sqlmeshv1grpc"
	ingestsqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/url"
	"os"
	"time"
)

type GitContextDump struct {
	CloneUrl  string `json:"clone_url"`
	Branch    string `json:"branch"`
	CommitSha string `json:"commit_sha"`
}

type IngestMetadataRequestDump struct {
	ApiMeta           json.RawMessage            `json:"api_meta"`
	Models            json.RawMessage            `json:"models"`
	ModelDetails      map[string]json.RawMessage `json:"model_details"`
	ModelLineage      map[string]json.RawMessage `json:"model_lineage"`
	Files             json.RawMessage            `json:"files"`
	Environments      json.RawMessage            `json:"environments"`
	FileContent       map[string]json.RawMessage `json:"file_content"`
	UploaderVersion   string                     `json:"uploader_version"`
	UploaderBuildTime string                     `json:"uploader_build_time"`
	StateAt           time.Time                  `json:"state_at"`
	GitContext        *GitContextDump            `json:"git_context"`
}

func DumpMetadata(output *ingestsqlmeshv1.IngestMetadataRequest, filename string) error {
	outputRaw := IngestMetadataRequestDump{
		ApiMeta:           output.ApiMeta,
		Models:            output.Models,
		ModelDetails:      lo.MapValues(output.ModelDetails, func(v []byte, k string) json.RawMessage { return v }),
		ModelLineage:      lo.MapValues(output.ModelLineage, func(v []byte, k string) json.RawMessage { return v }),
		Files:             output.Files,
		Environments:      output.Environments,
		FileContent:       lo.MapValues(output.FileContent, func(v []byte, k string) json.RawMessage { return v }),
		UploaderVersion:   output.UploaderVersion,
		UploaderBuildTime: output.UploaderBuildTime,
		StateAt:           output.StateAt.AsTime(),
	}

	if output.GitContext != nil {
		outputRaw.GitContext = &GitContextDump{
			CloneUrl:  output.GitContext.CloneUrl,
			Branch:    output.GitContext.Branch,
			CommitSha: output.GitContext.CommitSha,
		}
	}

	asJson, err := json.MarshalIndent(outputRaw, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, asJson, 0644)
}

func UploadMetadata(ctx context.Context, output *ingestsqlmeshv1.IngestMetadataRequest, endpoint string, token string) error {
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	oauthTokenSource, err := LongLivedTokenSource(token, parsedEndpoint)
	if err != nil {
		return err
	}
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(oauthTokenSource),
		grpc.WithAuthority(parsedEndpoint.Host),
	}

	conn, err := grpc.DialContext(ctx, grpcEndpoint(parsedEndpoint), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	sqlMeshServiceClient := ingestsqlmeshv1grpc.NewSqlMeshServiceClient(conn)
	resp, err := sqlMeshServiceClient.IngestMetadata(ctx, output)
	if err != nil {
		return err
	}
	logrus.Infof("Metadata uploaded successfully: %s", resp.String())
	return nil
}

func UploadExecutionLog(ctx context.Context, output *ingestsqlmeshv1.IngestExecutionRequest, endpoint string, token string) error {
	parsedEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	oauthTokenSource, err := LongLivedTokenSource(token, parsedEndpoint)
	if err != nil {
		return err
	}
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(oauthTokenSource),
		grpc.WithAuthority(parsedEndpoint.Host),
	}

	conn, err := grpc.DialContext(ctx, grpcEndpoint(parsedEndpoint), opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	sqlMeshServiceClient := ingestsqlmeshv1grpc.NewSqlMeshServiceClient(conn)
	resp, err := sqlMeshServiceClient.IngestExecution(ctx, output)
	if err != nil {
		return err
	}
	logrus.Infof("Logs uploaded successfully: %s", resp.String())
	return nil
}

func grpcEndpoint(endpoint *url.URL) string {
	port := endpoint.Port()
	if port == "" {
		port = "443"
	}
	return fmt.Sprintf("%s:%s", endpoint.Hostname(), port)
}
