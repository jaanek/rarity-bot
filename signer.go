package main

import (
	"fmt"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/eth"
	"github.com/jaanek/jeth/ui"
	"github.com/jaanek/jethwallet/keystore"
	"github.com/jaanek/jethwallet/wallet"
	"github.com/ledgerwatch/erigon/common"
)

func NewKeystoreTxSigner(term ui.Screen, keystorePath string) KeystoreTxSigner {
	return &signer{
		term:         term,
		keystorePath: keystorePath,
		passwords:    make(map[common.Address][]byte),
	}
}

type KeystoreTxSigner interface {
	eth.TxSigner
	AskPasswordFor(addr common.Address) error
	SetPasswordFor(addr common.Address, pass string) error
}

type signer struct {
	term         ui.Screen
	keystorePath string
	passwords    map[common.Address][]byte
}

func (s *signer) GetSignedRawTx(chainID uint256.Int, nonce uint64, from common.Address, to *common.Address, value *uint256.Int, input []byte, gasLimit uint64, gasPrice, gasTip, gasFeeCap *uint256.Int) ([]byte, error) {
	// tx, err := wallet.NewTx(chainID, nonce, to, value, input, gasLimit, gasPrice, gasTip, gasFeeCap)
	tx, err := wallet.NewTx(chainID, nonce, to, value, input, gasLimit, gasPrice, nil, gasFeeCap) // force legacy tx
	if err != nil {
		return nil, fmt.Errorf("Error while creating a tx: %w", err)
	}
	pass, ok := s.passwords[from]
	if !ok {
		return nil, fmt.Errorf("No password for %v provided", from)
	}
	signed, err := keystore.SignTxWithPassphrase(s.term, s.keystorePath, from, tx, string(pass))
	if err != nil {
		return nil, err
	}
	return wallet.EncodeTx(signed)
}

func (s *signer) AskPasswordFor(addr common.Address) error {
	s.term.Print(fmt.Sprintf("*** Enter passphrase (not echoed) account: %v ...", addr))
	pass, err := s.term.ReadPassword()
	if err != nil {
		return err
	}
	s.passwords[addr] = pass
	return nil
}

func (s *signer) SetPasswordFor(addr common.Address, pass string) error {
	s.passwords[addr] = []byte(pass)
	return nil
}
