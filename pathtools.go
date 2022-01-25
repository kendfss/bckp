package main

import (
	"errors"
	"fmt"
	_ "io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func FixExt(ext string) string {
	if m, err := regexp.MatchString(`(\.\w+)+$`, ext); m && err == nil {
		return ext
	} else {
		DATAERR.Warnf(`Couldn't recognize extension pattern (for:%s). Prepending ".": (match=%t, err=%w)`, ext, m, err)
	}
	return "." + ext
}

func splitExt(path string) (string, string) {
	parts := strings.Split(path, ".")

	if len(parts) >= 2 {
		return strings.Join(parts[:len(parts)-1], "."), parts[len(parts)-1]
	}
	return path, ""
}

func isDir(path string) bool {
	stat, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
	}
	return stat.Mode().IsDir()
}

func Files(root string) []string {
	paths := make([]string, 0)
	for _, path := range Listdir(root) {
		stat, err := os.Lstat(path)
		if err != nil {
			log.Fatal(err)
		}
		switch mode := stat.Mode(); {
		case mode.IsRegular():
			paths = append(paths, path)
		case mode.IsDir():
			paths = append(paths, Files(path)...)
		}
	}
	return paths
}

func Listdir(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	rack := make([]string, 0)
	for _, file := range files {
		rack = append(rack, filepath.Join(path, file.Name()))
	}
	return rack
}

func Merge(receiver *[]string, giver []string) {
	for _, str := range giver {
		*receiver = append(*receiver, str)
	}
}

// Check if an argument points to a real path on the system
func Exists(path string) bool {
	_, err := os.Lstat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func subtract(path, root string) string {
	max := len(path)
	if strings.Index(path, root) == 0 {
		max = len(root)
	}
	return path[:max]
}

func AssureTree(path string) error {
	return os.MkdirAll(filepath.Dir(path), os.ModeDir|os.ModePerm)
}
