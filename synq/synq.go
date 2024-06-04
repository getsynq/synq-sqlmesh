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
)

type ingestMetadataRequest struct {
	ApiMeta      json.RawMessage            `json:"api_meta"`
	Models       json.RawMessage            `json:"models"`
	ModelDetails map[string]json.RawMessage `json:"model_details"`
	ModelLineage map[string]json.RawMessage `json:"model_lineage"`
	Files        json.RawMessage            `json:"files"`
	Environments json.RawMessage            `json:"environments"`
}

func DumpMetadata(output *ingestsqlmeshv1.IngestMetadataRequest, filename string) error {
	outputRaw := ingestMetadataRequest{
		ApiMeta:      output.ApiMeta,
		Models:       output.Models,
		ModelDetails: lo.MapValues(output.ModelDetails, func(v []byte, k string) json.RawMessage { return v }),
		ModelLineage: lo.MapValues(output.ModelLineage, func(v []byte, k string) json.RawMessage { return v }),
		Files:        output.Files,
		Environments: output.Environments,
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
		panic(err)
	}
	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(oauthTokenSource),
		grpc.WithAuthority(parsedEndpoint.Host),
	}

	conn, err := grpc.DialContext(ctx, grpcEndpoint(parsedEndpoint), opts...)
	if err != nil {
		panic(err)
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

func grpcEndpoint(endpoint *url.URL) string {
	port := endpoint.Port()
	if port == "" {
		port = "443"
	}
	return fmt.Sprintf("%s:%s", endpoint.Hostname(), port)
}
