package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kendfss/namespacer"
)

var (
	settings = new(Settings)

	dstArg, dst                         string
	archSwitch, unzipSwitch, nestSwitch bool
)

func Bool(ptr *bool, name string, val bool, msg string) (helpText string) {
	helpText = fmt.Sprintf("%s [default: %t]", msg, val)
	flag.BoolVar(ptr, name, val, helpText)
	return
}

func String(ptr *string, name, val, msg string) (helpText string) {
	helpText = fmt.Sprintf("%s [default: %s]", msg, val)
	flag.StringVar(ptr, name, val, helpText)
	return
}

func init() {
	// if err := settings.Initialize(); err != nil {
	// 	panic(err)
	// }

	String(&dstArg, "d", ".", "place you want to keep your archive/backup")
	Bool(&archSwitch, "a", false, "Create archive instead of directory")
	Bool(&unzipSwitch, "u", false, "Unzip arguments")
	Bool(&nestSwitch, "n", false, "Unzip each argument into its own directory")
}

func main() {
	flag.Parse()

	var err error
	dstArg, err = filepath.Abs(dstArg)
	must(err)
	if isDir(dstArg) && archSwitch {
		dstArg = filepath.Join(dstArg, "bckp.zip")
		dstArg, err = namespacer.SpacedName(dstArg)
		must(err)
	}
	// fmt.Println(dstArg)
	// return
	if unzipSwitch {
		must(AssureTree(dstArg))
		// fmt.Printf("Extracting\n\t%#v\nTo\n\t%s", flag.Args(), dstArg)
		for _, arg := range flag.Args() {
			dst = dstArg
			if nestSwitch {
				parts := strings.Split(filepath.Join(dst, filepath.Base(arg)), ".")
				if len(parts) > 1 {
					dst = strings.Join(parts[:len(parts)-1], ".")
				}
			}
			fmt.Printf("Extracting\n\t%#v\nTo\n\t%s\n", arg, dst)
			unzipSource(arg, dstArg)
		}
	} else if archSwitch {
		err = zipSources(dstArg, flag.Args()...)
		must(err)
	} else {
		panic("not implemented for directory-writing yet")
	}
}

func zipSources(target string, sources ...string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	for _, source := range sources {
		f := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// 3. Create a local file header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// set compression
			header.Method = zip.Deflate

			// 4. Set relative path of a file as the header name
			header.Name, err = filepath.Rel(filepath.Dir(source), path)
			if err != nil {
				return err
			}
			if info.IsDir() {
				header.Name += "/"
			}

			// 5. Create writer for the file header and save content of the file
			headerWriter, err := writer.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(headerWriter, f)
			return err
		}

		if err := filepath.Walk(source, f); err != nil {
			return err
		}
	}
	return nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
