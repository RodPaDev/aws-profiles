package utils

import (
	"errors"
	"os"
	"os/user"
	"runtime"
)

func GetHomePath() (string, error) {
	if runtime.GOOS == "windows" {
		if home := os.Getenv("USERPROFILE"); home != "" {
			return home, nil
		}
	} else {
		if home := os.Getenv("HOME"); home != "" {
			return home, nil
		}
	}

	usr, err := user.Current()
	if err == nil && usr.HomeDir != "" {
		return usr.HomeDir, nil
	}

	return "", errors.New("Home Directory could not be located.")
}
