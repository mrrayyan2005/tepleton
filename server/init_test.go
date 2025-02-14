package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tepleton/tepleton/libs/log"

	"github.com/tepleton/tepleton-sdk/server/mock"
	"github.com/tepleton/tepleton-sdk/wire"
	tcmd "github.com/tepleton/tepleton/cmd/tepleton/commands"
)

// TODO update
func TestInitCmd(t *testing.T) {
	defer setupViper(t)()

	logger := log.NewNopLogger()
	cfg, err := tcmd.ParseConfig()
	require.Nil(t, err)
	ctx := NewContext(cfg, logger)
	cdc := wire.NewCodec()
	appInit := AppInit{
		AppGenState: mock.AppGenState,
		AppGenTx:    mock.AppGenTx,
	}
	cmd := InitCmd(ctx, cdc, appInit)
	err = cmd.RunE(nil, nil)
	require.NoError(t, err)
}

func TestGenTxCmd(t *testing.T) {
	// TODO
}

func TestTestnetFilesCmd(t *testing.T) {
	// TODO
}

func TestSimpleAppGenTx(t *testing.T) {
	// TODO
}

func TestSimpleAppGenState(t *testing.T) {
	// TODO
}
