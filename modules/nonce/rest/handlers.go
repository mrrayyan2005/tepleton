package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	sdk "github.com/tepleton/tepleton-sdk"
	"github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/commands"
	"github.com/tepleton/tepleton-sdk/client/commands/query"
	"github.com/tepleton/tepleton-sdk/errors"
	"github.com/tepleton/tepleton-sdk/modules/coin"
	"github.com/tepleton/tepleton-sdk/modules/nonce"
	"github.com/tepleton/tepleton-sdk/stack"
	wire "github.com/tepleton/go-wire"
	"github.com/tepleton/tmlibs/common"
)

// doQueryNonce is the HTTP handlerfunc to query a nonce
func doQueryNonce(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	signature := args["signature"]
	actor, err := commands.ParseActor(signature)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	var h int
	qHeight := r.URL.Query().Get("height")
	if qHeight != "" {
		h, err = strconv.Atoi(qHeight)
		if err != nil {
			common.WriteError(w, err)
			return
		}
	}

	actor = coin.ChainAddr(actor)
	key := nonce.GetSeqKey([]sdk.Actor{actor})
	key = stack.PrefixedKey(nonce.NameNonce, key)

	prove := !viper.GetBool(commands.FlagTrustNode)

	// query sequence number
	data, height, err := query.Get(key, h, prove)
	if client.IsNoDataErr(err) {
		err = fmt.Errorf("nonce empty for address: %q", signature)
		common.WriteError(w, err)
		return
	} else if err != nil {
		common.WriteError(w, err)
		return
	}

	// unmarshal sequence number
	var seq uint32
	err = wire.ReadBinaryBytes(data, &seq)
	if err != nil {
		msg := fmt.Sprintf("Error reading sequence for %X", key)
		common.WriteError(w, errors.ErrInternal(msg))
		return
	}

	// write result to client
	if err := query.FoutputProof(w, seq, height); err != nil {
		common.WriteError(w, err)
	}
}

// mux.Router registrars

// RegisterQueryNonce is a mux.Router handler that exposes GET
// method access on route /query/nonce/{signature} to query nonces
func RegisterQueryNonce(r *mux.Router) error {
	r.HandleFunc("/query/nonce/{signature}", doQueryNonce).Methods("GET")
	return nil
}

// End of mux.Router registrars
