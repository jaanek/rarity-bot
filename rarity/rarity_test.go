package rarity

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/eth"
	"github.com/jaanek/jeth/rpc"
	"github.com/jaanek/jeth/ui"
	"github.com/jaanek/jethwallet/keystore"
	"github.com/ledgerwatch/erigon/common"
)

var (
	term     ui.Screen
	txSigner keystore.TxSigner
	fromAddr common.Address
	endpoint rpc.Endpoint
	rarity   Contract
)

func TestMain(m *testing.M) {
	term = ui.NewTerminal(false)
	err := setUp()
	if err != nil {
		term.Errorf("Failed to setup. Err: %v\n", err)
		os.Exit(-1)
		return
	}
	retCode := m.Run()
	tearDown(term)
	os.Exit(retCode)
}

func setUp() error {
	txSigner = keystore.NewTxSigner(term, "/home/jaanek/.keystore-dev")
	fromAddr = common.HexToAddress("0x997a72d25791c4E2B717c094B924Fd5FFA825AFa")
	endpoint = rpc.NewEndpoint("http://localhost:8545")
	err := txSigner.SetPasswordFor(fromAddr, "")
	if err != nil {
		return err
	}
	// dev environment does not support eip-1559 tx-es
	txSigner.ForceLegaxyTx()

	// send initial eth to test accounts
	value := new(uint256.Int).SetUint64(1)
	_, receipt, err := eth.SendValue(term, endpoint, fromAddr, fromAddr, value, -1, txSigner)
	if err != nil {
		return err
	}

	// deploy contracts under test
	bin, err := eth.ReadHexFile("../build/rarity.bin")
	if err != nil {
		return err
	}
	_, receipt, err = Deploy(term, endpoint, fromAddr, bin, nil, -1, txSigner)
	if err != nil {
		return err
	}
	rarityAddr := common.HexToAddress(receipt.ContractAddress)
	rarity, err = New(term, endpoint, rarityAddr, &fromAddr, txSigner)
	if err != nil {
		return err
	}
	return nil
}

func tearDown(term ui.Screen) {
}

func TestSummon(t *testing.T) {
	classId, _ := uint256.FromBig(new(big.Int).SetUint64(11))
	_, receipt, err := rarity.Summon(classId, -1)
	if err != nil {
		t.Error(err)
		return
	}
	event := new(EventSummoned)
	err = rarity.GetEvent("summoned", event, receipt.Logs)
	if err != nil {
		t.Error(err)
		return
	}
	term.Print(fmt.Sprintf("summoned event: %v\n", event))
}

func TestSummoner(t *testing.T) {
	summonerId, _ := uint256.FromBig(new(big.Int).SetUint64(0))
	xp, adventureLog, class, level, err := rarity.Summoner(summonerId)
	if err != nil {
		t.Errorf("Error while calling summoner: err: %w\n", err)
	}
	term.Print(fmt.Sprintf("Summoner: %v, %v, %v, %v\n", xp, adventureLog, class, level))
}
