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
	IdGenTypeIncrease = iota //自增id
	IdGenTypeSnow            //雪花
)
