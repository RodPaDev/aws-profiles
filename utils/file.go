package utils

import (
	"errors"
	"io"
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

func CopyFile(path string, destination string) error {
	srcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
