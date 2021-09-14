package main

import (
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/eth"
	"github.com/jaanek/jeth/rpc"
	"github.com/jaanek/jeth/ui"
	"github.com/jaanek/rarity-bot/build/contracts"
	"github.com/ledgerwatch/erigon/common"
	"github.com/urfave/cli"
)

var summonerIds = []uint64{
	310839, 313281, 313710, 314929, 309672,
	315966, 317148, 318066, 318732, 306092,
	319076, 815955, 816237, 816357, 816414,
	816451, 816561, 816629, 816692, 816744,
	816786, 816819,
}

func RunCommand(term ui.Screen, ctx *cli.Context) error {
	// txSigner := NewKeystoreTxSigner(term, "/home/jaanek/.keystore-dev")
	// fromAddr := common.HexToAddress("0x997a72d25791c4E2B717c094B924Fd5FFA825AFa")
	// endpoint := rpc.NewEndpoint("http://localhost:8545")
	txSigner := NewKeystoreTxSigner(term, "/home/jaanek/.keystore")
	fromAddr := common.HexToAddress("0x09B1f3837E7B2F758Fab836B4e85b7AD0cA32f86")
	endpoint := rpc.NewEndpoint("https://rpc.ftm.tools")
	err := txSigner.AskPasswordFor(fromAddr)
	if err != nil {
		return err
	}
	// value := new(uint256.Int).SetUint64(1)
	// sendValue(term, endpoint, fromAddr, fromAddr, value, txSigner)
	// _, receipt, err := deployRarity(term, endpoint, fromAddr, txSigner)
	// if err != nil {
	// 	return err
	// }
	// rarityAddr := common.HexToAddress(receipt.ContractAddress)
	rarityAddr := common.HexToAddress("0xce761D788DF608BD21bdd59d6f4B54b2e27F25Bb")
	rarity, err := contracts.NewRarity(term, endpoint, rarityAddr, &fromAddr, txSigner)
	if err != nil {
		return err
	}
	listSummonersInfo(term, rarity)
	// adventure(term, rarity)
	return nil
}

func deployRarity(term ui.Screen, endpoint rpc.Endpoint, from common.Address, txSigner eth.TxSigner) (string, *eth.TxReceipt, error) {
	bin, err := eth.ReadHexFile("/opt/eth/rarity-bot/rarity.bin")
	if err != nil {
		return "", nil, err
	}
	hash, receipt, err := contracts.DeployRarity(term, endpoint, from, bin, nil, -1, txSigner)
	if hash != "" {
		term.Print(fmt.Sprintf("[Sent] Tx hash: %s", hash))
	}
	if err != nil {
		term.Errorf("Error while deploying rarity contract: err: %w\n", err)
	} else {
		term.Print(fmt.Sprintf("Received receipt: %+v", receipt))
	}
	return hash, receipt, nil
}

func sendValue(term ui.Screen, endpoint rpc.Endpoint, from common.Address, to common.Address, value *uint256.Int, txSigner eth.TxSigner) (string, *eth.TxReceipt, error) {
	hash, receipt, err := eth.Send(term, endpoint, from, to, value, []byte{}, -1, txSigner)
	if hash != "" {
		term.Print(fmt.Sprintf("[Sent] Tx hash: %s", hash))
	}
	if err != nil {
		term.Errorf("Error while sending value: err: %w\n", err)
	} else {
		term.Print(fmt.Sprintf("[sent value] Received receipt: %+v", receipt))
	}
	return hash, receipt, nil
}

func adventure(term ui.Screen, rarity contracts.Rarity) {
	for _, summoner := range summonerIds {
		summonerId, _ := uint256.FromBig(new(big.Int).SetUint64(summoner))
		hash, receipt, err := rarity.Adventure(summonerId, -1)
		if hash != "" {
			term.Print(fmt.Sprintf("[Sent] Tx hash: %s", hash))
		}
		if err != nil {
			term.Errorf("Error while calling adventure: err: %w\n", err)
			continue
		}
		term.Print(fmt.Sprintf("Received receipt: %+v", receipt))
	}
}

func listSummonersInfo(term ui.Screen, rarity contracts.Rarity) {
	for _, summoner := range summonerIds {
		summonerId, _ := uint256.FromBig(new(big.Int).SetUint64(summoner))
		xp, adventureLog, class, level, err := rarity.Summoner(summonerId)
		if err != nil {
			term.Errorf("Error while calling summoner: err: %w\n", err)
			continue
		}
		term.Print(fmt.Sprintf("Summoner: %v, %v, %v, %v\n", xp, adventureLog, class, level))
	}
}
