package config

type AdaptersConfig struct {
	IsDeterministic bool
}

func MakeAdaptersConfig(isDet bool) AdaptersConfig {
	return AdaptersConfig{IsDeterministic: isDet}
}
