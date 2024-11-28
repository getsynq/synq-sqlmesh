# synq-sqlmesh
synq-sqlmesh is a client application to integrate locally running SQLMesh to Synq


## How to use

If you have Golang installed, you can build the binary yourself, otherwise download appropriate binary from the [releases screen](https://github.com/getsynq/synq-sqlmesh/releases) (darwin == macOS).

`synq-sqlmesh` uses `web` module of `sqlmesh` to collect metadata. It was tested with versions ` >= 0.96.x`. If you do not have `web` module installed do

```bash
pip install "sqlmesh[web]"
```

All commands assume `sqlmesh` command is available in the `PATH`. If that is not the case, `--sqlmesh-cmd` could be used to point synq-sqlmesh to proper location.

### Dump metadata for inspection

```bash
cd sqlmesh-project
synq-sqlmesh collect meta.json
```

### Automatic upload to Synq
Run the following code after you've executed `sqlmesh run`, `sqlmesh audit`, and `sqlmesh test`. For example, you can add it to your Airflow code if you use Airflow for orchestrating.

```bash
export SYNQ_TOKEN=<token>
synq-sqlmesh upload
```

The `token` value is obtained from the SYNQ UI when you click 'create' under SQLMesh integration

### Advanced usage

```bash
❯ synq-sqlmesh --help
Small utility to collect SQLMesh metadata information and upload it to Synq

Usage:
  synq-sqlmesh [command]

Available Commands:
  collect     Collect metadata information from SQLMesh and store to the file
  upload      Collect metadata information from SQLMesh and send to Synq API

Flags:
  -h, --help                         help for synq-sqlmesh
      --sqlmesh-cmd string           SQLMesh launcher location (default "sqlmesh")
      --sqlmesh-project-dir string   Location of SQLMesh project directory (default ".")
      --sqlmesh-ui-host string       SQLMesh UI host (default "localhost")
      --sqlmesh-ui-port int          SQLMesh UI port (default 8080)
      --sqlmesh-ui-start             Launch and control SQLMesh UI process automatically (default true)
      --synq-endpoint string         Synq API endpoint URL (default "https://developer.synq.io/")
      --synq-token string            Synq API token

Use "synq-sqlmesh [command] --help" for more information about a command.
```


## Building from source code

Make sure at least golang 1.21 is installed and do:

```bash
go generate
go build -o synq-sqlmesh
./synq-sqlmesh version
```
