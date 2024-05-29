package sqlmesh

import (
	"github.com/sirupsen/logrus"
	"net/url"
	"os"
	"time"
)

func WaitForSqlMeshToStart(url url.URL) {

	logrus.Info("Waiting for sqlmesh to start")
	api := NewAPIClient(url)
	t := time.Now()
	for t.Add(30 * time.Second).After(time.Now()) {
		_, err := api.Health()
		if err == nil {
			return
		}
		logrus.WithError(err).Error("Failed to get health of SqlMesh api")
		time.Sleep(1 * time.Second)
	}
	logrus.Error("SqlMesh did not start in time")
	os.Exit(1)
}
