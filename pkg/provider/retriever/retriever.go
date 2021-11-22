package retriever

type Retriever interface {
	Connect() error
	List() ([]string, error)
	Download(string, string) error
}
