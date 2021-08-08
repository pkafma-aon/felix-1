package util

import (
	"archive/zip"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Zippie struct {
	Base    string
	dir     string
	zipPath string
	zipFile *os.File
}

func NewZippie(dir string) *Zippie {
	base := filepath.Base(dir)
	tempFile, err := os.CreateTemp("", base+"_*.zip")
	if err != nil {
		log.Println(err)
	}
	ins := &Zippie{
		Base:    base,
		dir:     dir,
		zipPath: tempFile.Name(),
		zipFile: tempFile,
	}
	return ins
}

func (z *Zippie) Close() {
	if z.zipFile != nil {
		z.zipFile.Close()
	}
	err := os.Remove(z.zipPath)
	if err != nil {
		log.Println(err)
	}
}

func (z *Zippie) Zip() (path string, f *os.File, err error) {
	w := zip.NewWriter(z.zipFile)
	defer w.Close()

	err = filepath.Walk(z.dir, func(path string, info fs.FileInfo, err error) error {
		//fmt.Printf("Crawling: %#v\n", path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Ensure that `path` is not absolute; it should not start with "/".
		// This snippet happens to work because I don't use
		// absolute paths, but ensure your real-world code
		// transforms path into a zip-root relative path.
		thisPath, err := filepath.Rel(z.dir, path)
		if err != nil {
			return err
		}
		f, err := w.Create(thisPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	})

	return z.zipPath, z.zipFile, err
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
