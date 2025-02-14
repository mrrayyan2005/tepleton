package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	ton "github.com/tepleton/tepleton-sdk/cmd/ton/app"
	sdk "github.com/tepleton/tepleton-sdk/types"
	"github.com/tepleton/tepleton-sdk/x/auth"
	"github.com/spf13/cobra"
	"github.com/tepleton/tepleton/crypto"
)

func init() {
	rootCmd.AddCommand(txCmd)
	rootCmd.AddCommand(pubkeyCmd)
	rootCmd.AddCommand(addrCmd)
	rootCmd.AddCommand(hackCmd)
	rootCmd.AddCommand(rawBytesCmd)
}

var rootCmd = &cobra.Command{
	Use:          "tondebug",
	Short:        "Gaia debug tool",
	SilenceUsage: true,
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "Decode a ton tx from hex or base64",
	RunE:  runTxCmd,
}

var pubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Decode a pubkey from hex, base64, or bech32",
	RunE:  runPubKeyCmd,
}

var addrCmd = &cobra.Command{
	Use:   "addr",
	Short: "Convert an address between hex and bech32",
	RunE:  runAddrCmd,
}

var hackCmd = &cobra.Command{
	Use:   "hack",
	Short: "Boilerplate to Hack on an existing state by scripting some Go...",
	RunE:  runHackCmd,
}

var rawBytesCmd = &cobra.Command{
	Use:   "raw-bytes",
	Short: "Convert raw bytes output (eg. [10 21 13 255]) to hex",
	RunE:  runRawBytesCmd,
}

func runRawBytesCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}
	stringBytes := args[0]
	stringBytes = strings.Trim(stringBytes, "[")
	stringBytes = strings.Trim(stringBytes, "]")
	spl := strings.Split(stringBytes, " ")

	byteArray := []byte{}
	for _, s := range spl {
		b, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		byteArray = append(byteArray, byte(b))
	}
	fmt.Printf("%X\n", byteArray)
	return nil
}

func runPubKeyCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	pubkeyString := args[0]
	var pubKeyI crypto.PubKey

	// try hex, then base64, then bech32
	pubkeyBytes, err := hex.DecodeString(pubkeyString)
	if err != nil {
		var err2 error
		pubkeyBytes, err2 = base64.StdEncoding.DecodeString(pubkeyString)
		if err2 != nil {
			var err3 error
			pubKeyI, err3 = sdk.GetAccPubKeyBech32(pubkeyString)
			if err3 != nil {
				var err4 error
				pubKeyI, err4 = sdk.GetValPubKeyBech32(pubkeyString)

				if err4 != nil {
					return fmt.Errorf(`Expected hex, base64, or bech32. Got errors:
			hex: %v,
			base64: %v
			bech32 acc: %v
			bech32 val: %v
			`, err, err2, err3, err4)

				}
			}

		}
	}

	var pubKey crypto.PubKeyEd25519
	if pubKeyI == nil {
		copy(pubKey[:], pubkeyBytes)
	} else {
		pubKey = pubKeyI.(crypto.PubKeyEd25519)
		pubkeyBytes = pubKey[:]
	}

	cdc := ton.MakeCodec()
	pubKeyJSONBytes, err := cdc.MarshalJSON(pubKey)
	if err != nil {
		return err
	}
	accPub, err := sdk.Bech32ifyAccPub(pubKey)
	if err != nil {
		return err
	}
	valPub, err := sdk.Bech32ifyValPub(pubKey)
	if err != nil {
		return err
	}
	fmt.Println("Address:", pubKey.Address())
	fmt.Printf("Hex: %X\n", pubkeyBytes)
	fmt.Println("JSON (base64):", string(pubKeyJSONBytes))
	fmt.Println("Bech32 Acc:", accPub)
	fmt.Println("Bech32 Val:", valPub)
	return nil
}

func runAddrCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	addrString := args[0]
	var addr sdk.Address

	// try hex, then bech32
	var err error
	addr, err = hex.DecodeString(addrString)
	if err != nil {
		var err2 error
		addr, err2 = sdk.GetAccAddressBech32(addrString)
		if err2 != nil {
			var err3 error
			addr, err3 = sdk.GetValAddressBech32(addrString)

			if err3 != nil {
				return fmt.Errorf(`Expected hex or bech32. Got errors:
			hex: %v,
			bech32 acc: %v
			bech32 val: %v
			`, err, err2, err3)

			}
		}
	}

	accAddr, err := sdk.Bech32ifyAcc(addr)
	if err != nil {
		return err
	}
	valAddr, err := sdk.Bech32ifyVal(addr)
	if err != nil {
		return err
	}
	fmt.Println("Address:", addr)
	fmt.Println("Bech32 Acc:", accAddr)
	fmt.Println("Bech32 Val:", valAddr)
	return nil
}

func runTxCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Expected single arg")
	}

	txString := args[0]

	// try hex, then base64
	txBytes, err := hex.DecodeString(txString)
	if err != nil {
		var err2 error
		txBytes, err2 = base64.StdEncoding.DecodeString(txString)
		if err2 != nil {
			return fmt.Errorf(`Expected hex or base64. Got errors:
			hex: %v,
			base64: %v
			`, err, err2)
		}
	}

	var tx = auth.StdTx{}
	cdc := ton.MakeCodec()

	err = cdc.UnmarshalBinary(txBytes, &tx)
	if err != nil {
		return err
	}

	bz, err := cdc.MarshalJSON(tx)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(buf.String())
	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
