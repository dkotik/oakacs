package botswat

import (
	"errors"
	"net/http"

	"github.com/dkotik/oakacs/oakhttp"
)

func New(
	extractor ResponseExtractor,
	withOptions ...Option,
) (oakhttp.Middleware, error) {
	if extractor == nil {
		return nil, errors.New("cannot use a <nil> response extractor")
	}
	// verifier, err := New(withOptions)
	// if err != nil {
	// 	return nil, err
	// }
	return func(h oakhttp.Handler) oakhttp.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			return errors.New("unimplemented")
		}
	}, nil
}