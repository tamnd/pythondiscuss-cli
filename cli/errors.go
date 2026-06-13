package cli

import (
	"errors"

	"github.com/tamnd/pythondiscuss-cli/pythondiscuss"
)

func isNotFound(err error) bool {
	return errors.Is(err, pythondiscuss.ErrNotFound)
}
