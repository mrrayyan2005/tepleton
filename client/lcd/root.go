package lcd

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tepleton/tepleton/libs/log"

	tmserver "github.com/tepleton/tepleton/rpc/lib/server"
	cmn "github.com/tepleton/tepleton/libs/common"

	client "github.com/tepleton/tepleton-sdk/client"
	"github.com/tepleton/tepleton-sdk/client/context"
	keys "github.com/tepleton/tepleton-sdk/client/keys"
	rpc "github.com/tepleton/tepleton-sdk/client/rpc"
	tx "github.com/tepleton/tepleton-sdk/client/tx"
	"github.com/tepleton/tepleton-sdk/wire"
	auth "github.com/tepleton/tepleton-sdk/x/auth/client/rest"
	bank "github.com/tepleton/tepleton-sdk/x/bank/client/rest"
	gov "github.com/tepleton/tepleton-sdk/x/gov/client/rest"
	ibc "github.com/tepleton/tepleton-sdk/x/ibc/client/rest"
	slashing "github.com/tepleton/tepleton-sdk/x/slashing/client/rest"
	stake "github.com/tepleton/tepleton-sdk/x/stake/client/rest"
)

// ServeCommand will generate a long-running rest server
// (aka Light Client Daemon) that exposes functionality similar
// to the cli, but over rest
func ServeCommand(cdc *wire.Codec) *cobra.Command {
	flagListenAddr := "laddr"
	flagCORS := "cors"
	flagMaxOpenConnections := "max-open"

	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) error {
			listenAddr := viper.GetString(flagListenAddr)
			handler := createHandler(cdc)
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
				With("module", "rest-server")
			maxOpen := viper.GetInt(flagMaxOpenConnections)
			listener, err := tmserver.StartHTTPServer(listenAddr, handler, logger, tmserver.Config{MaxOpenConnections: maxOpen})
			if err != nil {
				return err
			}
			logger.Info("REST server started")

			// Wait forever and cleanup
			cmn.TrapSignal(func() {
				err := listener.Close()
				logger.Error("error closing listener", "err", err)
			})
			return nil
		},
	}
	cmd.Flags().StringP(flagListenAddr, "a", "tcp://localhost:1317", "Address for server to listen on")
	cmd.Flags().String(flagCORS, "", "Set to domains that can make CORS requests (* for all)")
	cmd.Flags().StringP(client.FlagChainID, "c", "", "ID of chain we connect to")
	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	cmd.Flags().IntP(flagMaxOpenConnections, "o", 1000, "Maximum open connections")
	return cmd
}

func createHandler(cdc *wire.Codec) http.Handler {
	r := mux.NewRouter()

	kb, err := keys.GetKeyBase() //XXX
	if err != nil {
		panic(err)
	}

	ctx := context.NewCoreContextFromViper()

	// TODO make more functional? aka r = keys.RegisterRoutes(r)
	r.HandleFunc("/version", CLIVersionRequestHandler).Methods("GET")
	r.HandleFunc("/node_version", NodeVersionRequestHandler(ctx)).Methods("GET")
	keys.RegisterRoutes(r)
	rpc.RegisterRoutes(ctx, r)
	tx.RegisterRoutes(ctx, r, cdc)
	auth.RegisterRoutes(ctx, r, cdc, "acc")
	bank.RegisterRoutes(ctx, r, cdc, kb)
	ibc.RegisterRoutes(ctx, r, cdc, kb)
	stake.RegisterRoutes(ctx, r, cdc, kb)
	slashing.RegisterRoutes(ctx, r, cdc, kb)
	gov.RegisterRoutes(ctx, r, cdc)
	return r
}
