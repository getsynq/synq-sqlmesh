# synq-sqlmesh

synq-sqlmesh is a client application to integrate locally running SQLMesh to SYNQ

## How to use

If you have Golang installed, you can build the binary yourself, otherwise download appropriate binary from the [releases screen](https://github.com/getsynq/synq-sqlmesh/releases) (darwin == macOS).

`synq-sqlmesh` uses `web` module of `sqlmesh` to collect metadata. It was tested with versions `>= 0.96.x`. If you do not have `web` module installed do

```bash
pip install "sqlmesh[web]"
```

All commands assume `sqlmesh` command is available in the `PATH`. If that is not the case, `--sqlmesh-cmd` could be used to point synq-sqlmesh to proper location.

### Dump metadata for inspection

```bash
cd sqlmesh-project
synq-sqlmesh collect meta.json
```

### Automatic upload to SYNQ

Run the following code after you've executed `sqlmesh run`, `sqlmesh audit`, and `sqlmesh test`. For example, you can add it to your Airflow code if you use Airflow for orchestrating.

```bash
export SYNQ_TOKEN=<token>
synq-sqlmesh upload
```

The `token` value is obtained from the SYNQ UI when you click 'create' under SQLMesh integration

### Upload execution log

```bash

# run normal sqlmesh

sqlmesh audit | tee audit.log
synq-sqlmesh upload_audit audit.log

sqlmesh run | tee run.log
synq-sqlmesh upload_audit run.log

```

### Advanced usage

```bash
❯ synq-sqlmesh --help
Small utility to collect SQLMesh metadata information and upload it to SYNQ

Usage:
  synq-sqlmesh [command]

Available Commands:
  collect      Collect metadata information from SQLMesh and store to the file
  completion   Generate the autocompletion script for the specified shell
  help         Help about any command
  upload       Collect metadata information from SQLMesh and send to SYNQ API
  upload_audit Sends to SYNQ output of `audit` command
  upload_run   Sends to SYNQ output of `run` command
  version      Print the version number of synq-sqlmesh

Flags:
  -h, --help                                          help for synq-sqlmesh
      --sqlmesh-cmd string                            SQLMesh launcher location (default "sqlmesh")
      --sqlmesh-collect-file-content                  If content of the project files should be collected
      --sqlmesh-collect-file-content-exclude string   File patterns to exclude content (default "*.log")
      --sqlmesh-collect-file-content-include string   File patterns to include content (default "external_models.yaml,models/**/*.sql,models/**/*.py,audits/**/*.sql,tests/**/*.yaml")
      --sqlmesh-project-dir string                    Location of SQLMesh project directory (default ".")
      --sqlmesh-ui-host string                        SQLMesh UI host (default "localhost")
      --sqlmesh-ui-port int                           SQLMesh UI port (default 8080)
      --sqlmesh-ui-start                              Launch and control SQLMesh UI process automatically (default true)
      --synq-endpoint string                          SYNQ API endpoint URL (default "https://developer.synq.io/")
      --synq-token string                             SYNQ API token

Use "synq-sqlmesh [command] --help" for more information about a command.
```

## Building from source code

Make sure at least golang 1.21 is installed and do:

```bash
go generate
go build -o synq-sqlmesh
./synq-sqlmesh version
```

## Troubleshooting

If you encounter issues using `synq-sqlmesh`, check the following common problems:

**1. `sqlmesh` or `sqlmesh[web]` not installed**

- Ensure you have installed SQLMesh with the web module:
  ```bash
  pip install "sqlmesh[web]"
  ```
- Make sure the `sqlmesh` command is available in your `PATH`, or use the `--sqlmesh-cmd` flag to specify its location.

**2. `SYNQ_TOKEN` not set or invalid**

- The `SYNQ_TOKEN` environment variable is required for uploading data to SYNQ.
- If you see an error like `SYNQ_TOKEN environment variable is not set`, obtain a token from the SYNQ UI and set it:
  ```bash
  export SYNQ_TOKEN=<your_token>
  ```

**3. SQLMesh UI fails to start or connect**

- The tool tries to start the SQLMesh UI by default. If it fails to connect, you may see errors like `SQLMesh did not start in time`.
- Check that you can run `sqlmesh ui` manually and access it at the configured host/port (default: `localhost:8080`).
- You can disable automatic UI startup with `--sqlmesh-ui-start=false` if you want to manage the process yourself.

**4. Network or API errors**

- If uploading fails, ensure you have network connectivity and the SYNQ API endpoint is reachable (`https://developer.synq.io/` by default).
- If you are behind a proxy or firewall, ensure it allows outbound connections to the SYNQ API.

**5. File content not collected as expected**

- By default, only certain file patterns are included. Use `--sqlmesh-collect-file-content` and adjust `--sqlmesh-collect-file-content-include`/`--sqlmesh-collect-file-content-exclude` as needed.
- If you see errors about file patterns, check your glob syntax.

**6. General errors**

- Most errors are printed to the console. If you see a message like `Failed to get meta information` or `Failed to upload execution log`, check the logs for more details.
- Run with increased verbosity (if supported) or check the logs for more context.

**7. Version compatibility**

- This tool was tested with SQLMesh versions `>= 0.96.x`. If you use a different version, compatibility is not guaranteed.
