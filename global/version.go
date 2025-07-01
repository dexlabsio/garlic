package global

import "os"

func version() string {
	if val, ok := os.LookupEnv("XCI_APP_VERSION"); ok {
		return val
	}

	return "undefined"
}

var Version = version()
