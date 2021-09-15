package contracts

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/eth"
	"github.com/jaanek/jeth/rpc"
	"github.com/jaanek/jeth/ui"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/core/types"
)

func DeployRarity(term ui.Screen, ep rpc.Endpoint, fromAddr common.Address, bin []byte, value *uint256.Int, waitTime time.Duration, txSigner eth.TxSigner) (string, *eth.TxReceipt, error) {
	typeNames := []string{}
	values := []string{}
	return eth.Deploy(term, ep, fromAddr, bin, value, typeNames, values, waitTime, txSigner)
}

var rarityMethods = []eth.MethodSpec{{
	Name:    "summon",
	Inputs:  []string{"uint256"},
	Outputs: []string{},
}, {
	Name:    "adventure",
	Inputs:  []string{"uint256"},
	Outputs: []string{},
}, {
	Name:    "level_up",
	Inputs:  []string{"uint256"},
	Outputs: []string{},
}, {
	Name:    "adventurers_log",
	Inputs:  []string{"uint256"},
	Outputs: []string{"uint256"},
}, {
	Name:    "balanceOf",
	Inputs:  []string{"address"},
	Outputs: []string{"uint256"},
}, {
	Name:    "summoner",
	Inputs:  []string{"uint256"},
	Outputs: []string{"uint256", "uint256", "uint256", "uint256"},
}}

var rarityEvents = []eth.EventSpec{{
	Name:      "summoned",
	TopicArgs: []string{"address"},
	DataArgs:  []string{"uint256", "uint256"},
}}

type RarityEventSummoned struct {
	Owner    common.Address
	Class    *big.Int
	Summoner *big.Int
}

type rarityContract struct {
	name         string
	term         ui.Screen
	endpoint     rpc.Endpoint
	contractAddr common.Address
	fromAddr     *common.Address
	txSigner     eth.TxSigner
	methods      map[string]eth.Method
	events       map[string]eth.Event
}

type Rarity interface {
	GetEvent(eventName string, out interface{}, logs []types.Log) error
	Summon(class *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	Adventure(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	LevelUp(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	AdventurersLog(summoner *uint256.Int) (*uint256.Int, error)
	BalanceOf(address common.Address) (*uint256.Int, error)
	Summoner(summoner *uint256.Int) (*uint256.Int, *uint256.Int, *uint256.Int, *uint256.Int, error)
}

func NewRarity(term ui.Screen, ep rpc.Endpoint, contractAddr common.Address, fromAddr *common.Address, txSigner eth.TxSigner) (Rarity, error) {
	methods := make(map[string]eth.Method, len(rarityMethods))
	for _, methodSpec := range rarityMethods {
		method, err := eth.NewMethod(term, ep, methodSpec.Name, methodSpec.Inputs, methodSpec.Outputs)
		if err != nil {
			return nil, fmt.Errorf("error in rarity %s method: %w", methodSpec.Name, err)
		}
		methods[methodSpec.Name] = method
	}
	events := make(map[string]eth.Event, len(rarityEvents))
	for _, eventSpec := range rarityEvents {
		event, err := eth.NewEvent(eventSpec.Name, eventSpec.TopicArgs, eventSpec.DataArgs)
		if err != nil {
			return nil, fmt.Errorf("error in rarity %s event: %w", eventSpec.Name, err)
		}
		events[eventSpec.Name] = event
	}
	return &rarityContract{
		name:         "rarity",
		term:         term,
		endpoint:     ep,
		contractAddr: contractAddr,
		fromAddr:     fromAddr,
		txSigner:     txSigner,
		methods:      methods,
		events:       events,
	}, nil
}

func (c *rarityContract) GetEvent(eventName string, out interface{}, logs []types.Log) error {
	event, ok := c.events[eventName]
	if !ok {
		return errors.New(fmt.Sprintf("Contract %s event %s not declared", c.name, eventName))
	}
	return event.ParseInto(out, logs)
}

func (c *rarityContract) Summon(class *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["summon"].Send(*c.fromAddr, c.contractAddr, nil, []string{class.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["summon"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *rarityContract) Adventure(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["adventure"].Send(*c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["adventure"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *rarityContract) LevelUp(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["level_up"].Send(*c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["level_up"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *rarityContract) AdventurersLog(summoner *uint256.Int) (*uint256.Int, error) {
	result, unpacked, err := c.methods["adventurers_log"].Call(c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()})
	if err != nil {
		return nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["adventurers_log"].Name(), err)
	}
	if len(unpacked) == 0 {
		// it can happen if call sent to wrong or missing contract
		return nil, errors.New(fmt.Sprintf("Contract %s method: %s. No result: %s returned from node", c.name, c.methods["adventurers_log"].Name(), string(result)))
	}
	return unpacked[0].ToUint256()
}

func (c *rarityContract) BalanceOf(address common.Address) (*uint256.Int, error) {
	result, unpacked, err := c.methods["balanceOf"].Call(c.fromAddr, c.contractAddr, nil, []string{address.Hex()})
	if err != nil {
		return nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["balanceOf"].Name(), err)
	}
	if len(unpacked) == 0 {
		// it can happen if call sent to wrong or missing contract
		return nil, errors.New(fmt.Sprintf("Contract %s method: %s. No result: %s returned from node", c.name, c.methods["balanceOf"].Name(), string(result)))
	}
	return unpacked[0].ToUint256()
}

func (c *rarityContract) Summoner(summoner *uint256.Int) (*uint256.Int, *uint256.Int, *uint256.Int, *uint256.Int, error) {
	result, unpacked, err := c.methods["summoner"].Call(c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["summoner"].Name(), err)
	}
	if len(unpacked) == 0 {
		// it can happen if call sent to wrong or missing contract
		return nil, nil, nil, nil, errors.New(fmt.Sprintf("Contract %s method: %s. No result: %s returned from node", c.name, c.methods["summoner"].Name(), string(result)))
	}
	val0, err := unpacked[0].ToUint256()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	val1, err := unpacked[1].ToUint256()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	val2, err := unpacked[2].ToUint256()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	val3, err := unpacked[3].ToUint256()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return val0, val1, val2, val3, nil
}
