package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cbroglie/mustache"
)

func build(inputDir, outputDir, layout string) error {
	os.RemoveAll(outputDir)

	err := filepath.Walk(inputDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if path == inputDir {
				return nil
			}

			newPath := path[len(inputDir):]

			if info, statErr := os.Stat(path); statErr == nil && info.IsDir() {
				os.MkdirAll(filepath.Join(outputDir, newPath), 0777)
				return nil
			}

			if info, statErr := os.Stat(path); statErr == nil && !info.IsDir() {
				outputFilePath := filepath.Join(outputDir, newPath)
				dir, _ := filepath.Split(outputFilePath)
				os.MkdirAll(dir, 0777)

				if filepath.Ext(path) == ".mustache" {
					newOutFilePath := strings.TrimSuffix(outputFilePath, ".mustache")
					newOutFilePath += ".html"
					return compileMustache(path, newOutFilePath, layout)
				} else {
					return copyFileContents(
						path,
						outputFilePath,
					)
				}
			}

			return nil

		})
	if err != nil {
		log.Println(err)
	}

	return nil
}

func compileMustache(in, out, layout string) error {
	outString, err := mustache.RenderFileInLayout(in, layout, nil)
	if err != nil {
		return err
	}

	return os.WriteFile(out, []byte(outString), 0777)
}
