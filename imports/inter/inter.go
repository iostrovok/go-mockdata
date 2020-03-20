package inter

type IImp interface {
	Add(file, pkg, path, alias string)
	UseFile(file, alias string) error
	UsePath(path, alias string) error
	List() []string
}
