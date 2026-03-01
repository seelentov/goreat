package static

type Folder string

const (
	FolderStatic    Folder = "static"
	FolderTemplates Folder = "templates"
)

type File struct {
	Name  string
	Path  string
	IsDir bool
	Ext   string
}
