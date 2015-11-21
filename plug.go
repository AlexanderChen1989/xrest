package xrest

type Plugger interface {
	Plug(Handler) Handler
}
