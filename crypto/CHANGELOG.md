# Changelog

## 0.2.0 (May 18, 2017)

BREAKING CHANGES:

- [hd] The following functions no longer take a `coin string` as argument: `ComputeAddress`, `AddrFromPubKeyBytes`, `ComputeAddressForPrivKey`, `ComputeWIF`, `WIFFromPrivKeyBytes`
- Changes to `PrivKey`, `PubKey`, and `Signature` (denoted `Xxx` below):
  - interfaces are renamed `XxxInner`, and are not for use outside the package, though they must be exposed for sake of serialization.
  - `Xxx` is now a struct that wraps the corresponding `XxxInner` interface

FEATURES:

- `github.com/tepleton/go-keys -> github.com/tepleton/go-crypto/keys` - command and lib for generating and managing encrypted keys
- [hd] New function `WIFFromPrivKeyBytes(privKeyBytes []byte, compress bool) string`
- Changes to `PrivKey`, `PubKey`, and `Signature` (denoted `Xxx` below):
  - Expose a new method `Unwrap() XxxInner` on the `Xxx` struct which returns the corresponding `XxxInner` interface
  - Expose a new method `Wrap() Xxx` on the `XxxInner` interface which returns the corresponding `Xxx` struct

IMPROVEMENTS:

- Update to use new `tmlibs` repository

## 0.1.0 (April 14, 2017)

Initial release

