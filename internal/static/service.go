package static

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrParseName = errors.New("failed parse name")

func GetAll(folder Folder) ([]File, error) {
	m := make([]File, 0)

	err := filepath.WalkDir(string(folder), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		f := strings.Split(d.Name(), ".")
		if len(f) == 0 {
			return ErrParseName
		}

		ext := "dir"
		if !d.IsDir() {
			ext = f[1]
		}

		m = append(m, File{
			Name:  f[0],
			Path:  path,
			IsDir: d.IsDir(),
			Ext:   ext,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

func Get(folder Folder, dir []string, filename string) (string, error) {
	d := filepath.Join(string(folder), filepath.Join(dir...))
	if err := ensureFolderExists(d); err != nil {
		return "", err
	}

	b, err := os.ReadFile(filepath.Join(d, filename))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func CreateOrUpdate(folder Folder, dir []string, filename, content string) error {
	d := filepath.Join(string(folder), filepath.Join(dir...))
	if err := ensureFolderExists(d); err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(d, filename), []byte(content), os.ModePerm)
}

func ensureFolderExists(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
