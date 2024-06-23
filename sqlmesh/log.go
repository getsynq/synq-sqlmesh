package sqlmesh

import (
	sqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"github.com/djherbis/times"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"time"
)

func CollectAuditLog(output *sqlmeshv1.IngestExecutionRequest, filename string) error {

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	output.StdOut = fileContent

	startedAt := time.Now().UTC()
	finishedAt := startedAt

	if t, err := times.Stat(filename); err == nil {
		startedAt = t.ModTime()
		finishedAt = t.ModTime()
		if t.HasBirthTime() {
			startedAt = t.BirthTime()
		}
		if t.HasChangeTime() {
			finishedAt = t.ChangeTime()
		}
	}

	output.StartedAt = timestamppb.New(startedAt)
	output.FinishedAt = timestamppb.New(finishedAt)

	return nil
}
