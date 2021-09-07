package finder

type Finder interface {
	Find() ([]string, error)
}
