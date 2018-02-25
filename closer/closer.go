package closer

import "github.com/hashicorp/go-multierror"

type Closer interface {
	Close() error
}

func CloseMany(closers ...Closer) error {
	var result *multierror.Error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result.ErrorOrNil()
}
