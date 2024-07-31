package model

const (
	EnvProd = "product"
	EnvDev  = "dev"
)

const (
	ConfigFileTypeYaml = "yaml"
	ConfigFileTypeToml = "toml"
	ConfigFileTypeJson = "json"
)

const (
	IdGenTypeIncrease = iota //Auto-increment ID
	IdGenTypeSnow            //snow
)
