package build

import _ "embed"

//go:embed version.txt
var Version string

//go:embed time.txt
var Time string
