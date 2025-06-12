package sqlmesh

import (
	"fmt"
	"os"
	"time"

	sqlmeshv1 "buf.build/gen/go/getsynq/api/protocolbuffers/go/synq/ingest/sqlmesh/v1"
	"github.com/djherbis/times"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CollectExecutionLog reads a run log file produced by `sqlmesh run` or
// `sqlmesh audit` and populates the provided request with its contents and
// timestamps.
func CollectExecutionLog(output *sqlmeshv1.IngestExecutionRequest, filename string) error {

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s: %w", filename, err)
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
