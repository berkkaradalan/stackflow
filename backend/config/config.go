package config


type Config struct {
	Env     *Env
}

func Load() (*Config, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		Env:     env,
	}, nil
}