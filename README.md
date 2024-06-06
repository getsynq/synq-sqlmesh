# synq-sqlmesh
synq-sqlmesh is a client application to integrate locally running SqlMesh to Synq


## How to use

If you have Golang installed, you can build the binary yourself, otherwise download appropriate binary from the [releases screen](https://github.com/getsynq/synq-sqlmesh/releases) (darwin == macOS).


Assuming `sqlmesh` command available in the `PATH`:

```bash
cd sqlmesh-project
synq-sqlmesh collect meta.json
```

### Automatic upload to Synq

```bash
export SYNQ_TOKEN=<token>
synq-sqlmesh upload
```


### Advanced usage

```bash
❯ synq-sqlmesh --help
Small utility to collect SqlMesh metadata information and upload it to Synq

Usage:
  synq-sqlmesh [command]

Available Commands:
  collect     Collect metadata information from SqlMesh and store to the file
  upload      Collect metadata information from SqlMesh and send to Synq API

Flags:
  -h, --help                         help for synq-sqlmesh
      --sqlmesh-cmd string           SqlMesh launcher location (default "sqlmesh")
      --sqlmesh-project-dir string   Location of SqlMesh project directory (default ".")
      --sqlmesh-ui-host string       SqlMesh UI host (default "localhost")
      --sqlmesh-ui-port int          SqlMesh UI port (default 8080)
      --sqlmesh-ui-start             Launch and control SqlMesh UI process automatically (default true)
      --synq-endpoint string         Synq API endpoint URL (default "https://developer.synq.io/")
      --synq-token string            Synq API token

Use "synq-sqlmesh [command] --help" for more information about a command.
```
