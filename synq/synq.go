package synq

import (
	"encoding/json"
	"github.com/getsynq/synq-sqlmesh/sqlmesh"
	"os"
)

func DumpMetadata(output *sqlmesh.SqlMeshMetadata, filename string) error {

	asJson, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, asJson, 0644)
}
