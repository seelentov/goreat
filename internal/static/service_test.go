package static_test

import (
	"goreat/internal/static"
	"os"
	"path/filepath"
	"testing"
)

var jsFile = static.File{
	Name:  "test",
	Path:  "static/js/test.js",
	IsDir: false,
	Ext:   "js",
}

var cssFile = static.File{
	Name:  "test",
	Path:  "static/css/test.css",
	IsDir: false,
	Ext:   "css",
}

var htmlFile = static.File{
	Name:  "test",
	Path:  "templates/test.html",
	IsDir: false,
	Ext:   "html",
}

var staticFiles = []static.File{jsFile, cssFile}
var templatesFiles = []static.File{htmlFile}

var testFiles = append(staticFiles, templatesFiles...)

var jsFolder = filepath.Join(string(static.FolderStatic), "js")
var cssFolder = filepath.Join(string(static.FolderStatic), "css")
var htmlFolder = filepath.Join(string(static.FolderTemplates))

var testFolders = []string{jsFolder, cssFolder, htmlFolder}

var jsPath = filepath.Join(jsFolder, "test.js")
var cssPath = filepath.Join(cssFolder, "test.css")
var htmlPath = filepath.Join(htmlFolder, "test.html")

func setupStaticServiceTest(t *testing.T) {
	t.Helper()

	for _, folder := range testFolders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			t.Fatalf("failed to create folder: %v", err)
		}
	}

	if err := os.WriteFile(jsPath, []byte("console.log('test')"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := os.WriteFile(cssPath, []byte("* {\nbackground-color: red;\n}"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	if err := os.WriteFile(htmlPath, []byte("<html>\n</html>"), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(string(static.FolderStatic))
		_ = os.Remove(string(static.FolderTemplates))
	})
}

func TestStaticService_GetAll(t *testing.T) {
	setupStaticServiceTest(t)

	fs, err := static.GetAll(static.FolderStatic)
	if err != nil {
		t.Error(err)
	}

	for _, f := range staticFiles {
		found := false
		for _, file := range fs {
			if file.Name == f.Name && file.Ext == f.Ext {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("file %s not found", f.Name+"."+f.Ext)
		}
	}

	fs, err = static.GetAll(static.FolderTemplates)
	if err != nil {
		t.Error(err)
	}

	for _, f := range templatesFiles {
		found := false
		for _, file := range fs {
			if file.Name == f.Name && file.Ext == f.Ext {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("file %s not found", f.Name+"."+f.Ext)
		}
	}
}

func TestStaticService_Get(t *testing.T) {
	setupStaticServiceTest(t)

	c, err := static.Get(static.FolderStatic, []string{"js"}, "test.js")
	if err != nil {
		t.Error(err)
	}

	if c != "console.log('test')" {
		t.Errorf("content not received")
	}
}

func TestStaticService_CreateOrUpdate(t *testing.T) {
	setupStaticServiceTest(t)

	err := static.CreateOrUpdate(static.FolderStatic, []string{"js"}, "test.js", "console.log('456')")
	if err != nil {
		t.Error(err)
	}

	c, err := static.Get(static.FolderStatic, []string{"js"}, "test.js")
	if err != nil {
		t.Error(err)
	}

	if c != "console.log('456')" {
		t.Errorf("content not changed")
	}

	if err = static.CreateOrUpdate(static.FolderStatic, []string{"js"}, "test_new.js", "console.log('456')"); err != nil {
		t.Error(err)
	}

	if c != "console.log('456')" {
		t.Errorf("content not changed")
	}

	t.Cleanup(func() {
		_ = os.Remove(filepath.Join("static", "js", "test.js"))
	})
}
