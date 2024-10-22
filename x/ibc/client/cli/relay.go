package cli

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tepleton/tepleton/libs/log"

	"github.com/tepleton/tepleton-sdk/client/context"
	sdk "github.com/tepleton/tepleton-sdk/types"
	wire "github.com/tepleton/tepleton-sdk/wire"
	"github.com/tepleton/tepleton-sdk/x/auth"
	authcmd "github.com/tepleton/tepleton-sdk/x/auth/client/cli"
	"github.com/tepleton/tepleton-sdk/x/ibc"
)

// flags
const (
	FlagFromChainID   = "from-chain-id"
	FlagFromChainNode = "from-chain-node"
	FlagToChainID     = "to-chain-id"
	FlagToChainNode   = "to-chain-node"
)

type relayCommander struct {
	cdc       *wire.Codec
	address   sdk.Address
	decoder   auth.AccountDecoder
	mainStore string
	ibcStore  string
	accStore  string

	logger log.Logger
}

// IBC relay command
func IBCRelayCmd(cdc *wire.Codec) *cobra.Command {
	cmdr := relayCommander{
		cdc:       cdc,
		decoder:   authcmd.GetAccountDecoder(cdc),
		ibcStore:  "ibc",
		mainStore: "main",
		accStore:  "acc",

		logger: log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	}

	cmd := &cobra.Command{
		Use: "relay",
		Run: cmdr.runIBCRelay,
	}

	cmd.Flags().String(FlagFromChainID, "", "Chain ID for ibc node to check outgoing packets")
	cmd.Flags().String(FlagFromChainNode, "tcp://localhost:26657", "<host>:<port> to tepleton rpc interface for this chain")
	cmd.Flags().String(FlagToChainID, "", "Chain ID for ibc node to broadcast incoming packets")
	cmd.Flags().String(FlagToChainNode, "tcp://localhost:36657", "<host>:<port> to tepleton rpc interface for this chain")

	cmd.MarkFlagRequired(FlagFromChainID)
	cmd.MarkFlagRequired(FlagFromChainNode)
	cmd.MarkFlagRequired(FlagToChainID)
	cmd.MarkFlagRequired(FlagToChainNode)

	viper.BindPFlag(FlagFromChainID, cmd.Flags().Lookup(FlagFromChainID))
	viper.BindPFlag(FlagFromChainNode, cmd.Flags().Lookup(FlagFromChainNode))
	viper.BindPFlag(FlagToChainID, cmd.Flags().Lookup(FlagToChainID))
	viper.BindPFlag(FlagToChainNode, cmd.Flags().Lookup(FlagToChainNode))

	return cmd
}

// nolint: unparam
func (c relayCommander) runIBCRelay(cmd *cobra.Command, args []string) {
	fromChainID := viper.GetString(FlagFromChainID)
	fromChainNode := viper.GetString(FlagFromChainNode)
	toChainID := viper.GetString(FlagToChainID)
	toChainNode := viper.GetString(FlagToChainNode)
	address, err := context.NewCoreContextFromViper().GetFromAddress()
	if err != nil {
		panic(err)
	}
	c.address = address

	c.loop(fromChainID, fromChainNode, toChainID, toChainNode)
}

func (c relayCommander) loop(fromChainID, fromChainNode, toChainID,
	toChainNode string) {

	ctx := context.NewCoreContextFromViper()
	// get password
	passphrase, err := ctx.GetPassphraseFromStdin(ctx.FromAddressName)
	if err != nil {
		panic(err)
	}

	ingressKey := ibc.IngressSequenceKey(fromChainID)

OUTER:
	for {
		time.Sleep(5 * time.Second)

		processedbz, err := query(toChainNode, ingressKey, c.ibcStore)
		if err != nil {
			panic(err)
		}

		var processed int64
		if processedbz == nil {
			processed = 0
		} else if err = c.cdc.UnmarshalBinary(processedbz, &processed); err != nil {
			panic(err)
		}

		lengthKey := ibc.EgressLengthKey(toChainID)
		egressLengthbz, err := query(fromChainNode, lengthKey, c.ibcStore)
		if err != nil {
			c.logger.Error("error querying outgoing packet list length", "err", err)
			continue OUTER //TODO replace with continue (I think it should just to the correct place where OUTER is now)
		}
		var egressLength int64
		if egressLengthbz == nil {
			egressLength = 0
		} else if err = c.cdc.UnmarshalBinary(egressLengthbz, &egressLength); err != nil {
			panic(err)
		}
		if egressLength > processed {
			c.logger.Info("Detected IBC packet", "number", egressLength-1)
		}

		seq := c.getSequence(toChainNode)

		for i := processed; i < egressLength; i++ {
			egressbz, err := query(fromChainNode, ibc.EgressKey(toChainID, i), c.ibcStore)
			if err != nil {
				c.logger.Error("error querying egress packet", "err", err)
				continue OUTER // TODO replace to break, will break first loop then send back to the beginning (aka OUTER)
			}

			err = c.broadcastTx(seq, toChainNode, c.refine(egressbz, i, passphrase))
			seq++
			if err != nil {
				c.logger.Error("error broadcasting ingress packet", "err", err)
				continue OUTER // TODO replace to break, will break first loop then send back to the beginning (aka OUTER)
			}

			c.logger.Info("Relayed IBC packet", "number", i)
		}
	}
}

func query(node string, key []byte, storeName string) (res []byte, err error) {
	return context.NewCoreContextFromViper().WithNodeURI(node).QueryStore(key, storeName)
}

func (c relayCommander) broadcastTx(seq int64, node string, tx []byte) error {
	_, err := context.NewCoreContextFromViper().WithNodeURI(node).WithSequence(seq + 1).BroadcastTx(tx)
	return err
}

func (c relayCommander) getSequence(node string) int64 {
	res, err := query(node, c.address, c.accStore)
	if err != nil {
		panic(err)
	}
	if nil != res {
		account, err := c.decoder(res)
		if err != nil {
			panic(err)
		}

		return account.GetSequence()
	}

	return 0
}

func (c relayCommander) refine(bz []byte, sequence int64, passphrase string) []byte {
	var packet ibc.IBCPacket
	if err := c.cdc.UnmarshalBinary(bz, &packet); err != nil {
		panic(err)
	}

	msg := ibc.IBCReceiveMsg{
		IBCPacket: packet,
		Relayer:   c.address,
		Sequence:  sequence,
	}

	ctx := context.NewCoreContextFromViper().WithSequence(sequence)
	res, err := ctx.SignAndBuild(ctx.FromAddressName, passphrase, []sdk.Msg{msg}, c.cdc)
	if err != nil {
		panic(err)
	}
	return res
}
