// Generated by: main
// TypeWriter: wrapper
// Directive: +gen on TxInner

package sdk

import (
	"github.com/tepleton/go-wire/data"
)

// Auto-generated adapters for happily unmarshaling interfaces
// Apache License 2.0
// Copyright (c) 2017 Ethan Frey (ethan.frey@tepleton.com)

type Tx struct {
	TxInner "json:\"unwrap\""
}

var TxMapper = data.NewMapper(Tx{})

func (h Tx) MarshalJSON() ([]byte, error) {
	return TxMapper.ToJSON(h.TxInner)
}

func (h *Tx) UnmarshalJSON(data []byte) (err error) {
	parsed, err := TxMapper.FromJSON(data)
	if err == nil && parsed != nil {
		h.TxInner = parsed.(TxInner)
	}
	return err
}

// Unwrap recovers the concrete interface safely (regardless of levels of embeds)
func (h Tx) Unwrap() TxInner {
	hi := h.TxInner
	for wrap, ok := hi.(Tx); ok; wrap, ok = hi.(Tx) {
		hi = wrap.TxInner
	}
	return hi
}

func (h Tx) Empty() bool {
	return h.TxInner == nil
}

/*** below are bindings for each implementation ***/