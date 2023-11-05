package duco

import "os"

func loadEnv() env {
	fnName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	if fnName == "" {
		return env{dev: true}
	}
	return env{
		fnName:  fnName,
		handler: os.Getenv("_HANDLER"),
		root:    os.Getenv("LAMBDA_TASK_ROOT"),
		api:     os.Getenv("AWS_LAMBDA_RUNTIME_API"),
	}
}

type env struct {
	dev     bool
	fnName  string
	handler string
	root    string
	api     string
}
