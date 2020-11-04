package fh_common

import (
    "os"
    log "logex-gls"
)

func CheckOrSetEnv(name string, default_val string) {
	if len(os.Getenv(name)) < 1 {
		os.Setenv(name, default_val)
		log.Info("Env:", name, "is missing, will use default: ", os.Getenv(name))
	} else {
		log.Info("Env:", name, " set to ", os.Getenv(name))
	}
}

func CheckOrSetEnvPassword(name string, default_val string) {
	if len(os.Getenv(name)) < 1 {
		os.Setenv(name, default_val)
		log.Info("Env (password):", name, "is missing, will use default")
	} else {
		log.Info("Env (password):", name, " set to <***password***>")
	}
} 
