package configs

import _ "embed"

//go:embed config.toml
var ConfigFileData []byte
