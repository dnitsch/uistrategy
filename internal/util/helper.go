package util

import (
	"log"
	"os"
	"path/filepath"
)

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("unable to get the user home dir")
	}
	return home
}

func WriteDataDir(datadir string) {
	os.MkdirAll(datadir, 0755)
}

func InitDirDeps() error {
	basePath := filepath.Join(".", ".uistrategy", "captures")
	return os.MkdirAll(basePath, os.ModePerm)
}

// CleanExit signals 0 exit code and should clean up any current process
func CleanExit() {
	os.Exit(0)
}

func Exit(err error) {
	log.Fatal(err)
}

func Str(s string) *string {
	return &s
}

func Int(v int) *int {
	return &v
}
