//nolint
package nonce

import (
	"fmt"

	wrsp "github.com/tepleton/wrsp/types"

	"github.com/tepleton/tepleton-sdk/errors"
)

var (
	errNoNonce      = fmt.Errorf("Tx doesn't contain nonce")
	errNotMember    = fmt.Errorf("Nonce contains non-permissioned member")
	errZeroSequence = fmt.Errorf("Sequence number cannot be zero")
	errNoSigners    = fmt.Errorf("There are no signers")
	errTxEmpty      = fmt.Errorf("The provided Tx is empty")

	unauthorized = wrsp.CodeType_Unauthorized
	badNonce     = wrsp.CodeType_BadNonce
	invalidInput = wrsp.CodeType_BaseInvalidInput
)

func ErrBadNonce(got, expected uint32) errors.TMError {
	return errors.WithCode(fmt.Errorf("Bad nonce sequence, got %d, expected %d", got, expected), badNonce)
}
func ErrNoNonce() errors.TMError {
	return errors.WithCode(errNoNonce, badNonce)
}
func ErrNotMember() errors.TMError {
	return errors.WithCode(errNotMember, unauthorized)
}
func ErrZeroSequence() errors.TMError {
	return errors.WithCode(errZeroSequence, invalidInput)
}
func ErrNoSigners() errors.TMError {
	return errors.WithCode(errNoSigners, invalidInput)
}
func ErrTxEmpty() errors.TMError {
	return errors.WithCode(errTxEmpty, invalidInput)
}
