package errx

type Error map[string][]string

func New() Error {
	return make(Error)
}

func (e Error) Add(key string, err string) {
	e[key] = append(e[key], err)
}

func (e Error) Get(key string) []string {
	return e[key]
}
