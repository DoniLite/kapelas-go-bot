package conf

var env *Env

func init() {
	env = &Env{}
	env.Load()
}

func GetEnv() *Env {
	return env
}