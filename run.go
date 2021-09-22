package main

import (
	"fmt"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/rpc"
	"github.com/jaanek/jeth/ui"
	"github.com/jaanek/rarity-bot/rarity"
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
	txSigner := NewKeystoreTxSigner(term, "/home/jaanek/.keystore")
	fromAddr := common.HexToAddress("0x09B1f3837E7B2F758Fab836B4e85b7AD0cA32f86")
	endpoint := rpc.NewEndpoint("https://rpc.ftm.tools")
	err := txSigner.AskPasswordFor(fromAddr)
	if err != nil {
		return err
	}
	rarityAddr := common.HexToAddress("0xce761D788DF608BD21bdd59d6f4B54b2e27F25Bb")
	rarity, err := rarity.New(term, endpoint, rarityAddr, &fromAddr, txSigner)
	if err != nil {
		return err
	}
	adventure(term, rarity)
	// levelUp(term, rarity)
	// summoners(term, rarity)
	// xpRequired(term, rarity, 4)
	return nil
}

func adventure(term ui.Screen, rarity rarity.Contract) {
	for _, summoner := range summonerIds {
		summonerId := new(uint256.Int).SetUint64(summoner)
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

func levelUp(term ui.Screen, rarity rarity.Contract) {
	for _, summoner := range summonerIds {
		summonerId := new(uint256.Int).SetUint64(summoner)
		hash, receipt, err := rarity.LevelUp(summonerId, -1)
		if hash != "" {
			term.Print(fmt.Sprintf("[Sent] Tx hash: %s", hash))
		}
		if err != nil {
			term.Errorf("Error while calling level_up: err: %w\n", err)
			continue
		}
		term.Print(fmt.Sprintf("Received receipt: %+v", receipt))
		break
	}
}

func summoners(term ui.Screen, rarity rarity.Contract) {
	for _, id := range summonerIds {
		summoner(term, rarity, id)
	}
}

func summoner(term ui.Screen, rarity rarity.Contract, id uint64) {
	summonerId := new(uint256.Int).SetUint64(id)
	xp, adventureLog, class, level, err := rarity.Summoner(summonerId)
	e18 := new(uint256.Int).SetUint64(1e18)
	if err != nil {
		term.Errorf("Error while calling summoner: err: %w\n", err)
		return
	}
	xpRequired, err := rarity.XpRequired(level)
	if err != nil {
		term.Errorf("Error while calling xp required: err: %w\n", err)
		return
	}
	term.Print(fmt.Sprintf("Summoner %v. xp: %v, log: %v, class: %v, level: %v, xp required for level_up: %v\n", id, xp.Div(xp, e18), adventureLog, class, level, xpRequired.Div(xpRequired, e18)))
}

func xpRequired(term ui.Screen, rarity rarity.Contract, level uint64) {
	lvl := new(uint256.Int).SetUint64(level)
	xp, err := rarity.XpRequired(lvl)
	e18 := new(uint256.Int).SetUint64(1e18)
	if err != nil {
		term.Errorf("Error while calling summoner: err: %w\n", err)
		return
	}
	term.Print(fmt.Sprintf("Xp required. Level %v, xp: %v\n", level, xp.Div(xp, e18)))
}
