package cmd

type CMDHandler interface {
	Handle(arg string) error
}
