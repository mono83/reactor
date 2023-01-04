package inflate

// Template is a function able to product Reactor
type Template func() (Reactor, error)

// NewTemplate constructs reactor template function configured with given
// values or value providers.
func NewTemplate(providers ...interface{}) Template {
	return func() (Reactor, error) {
		reactor := New()
		if err := reactor.Put(providers...); err != nil {
			return nil, err
		}
		return reactor, nil
	}
}
