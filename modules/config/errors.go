package config

type JsonConfigNotExit struct{}
type CannotParseConfigFile struct{}

func (e *JsonConfigNotExit) Error() string {
	return "Json config file not exist"
}

func (e *CannotParseConfigFile) Error() string {
	return "Cannot parse config pares"
}
