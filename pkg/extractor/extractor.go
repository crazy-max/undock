package extractor

// Handler is an extractor interface
type Handler interface {
	Extract() error
	Type() string
}

// Client represents an active extractor object
type Client struct {
	Handler
}
