package rarity

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

func Deploy(term ui.Screen, ep rpc.Endpoint, fromAddr common.Address, bin []byte, value *uint256.Int, waitTime time.Duration, txSigner eth.TxSigner) (string, *eth.TxReceipt, error) {
	return eth.Deploy(term, ep, fromAddr, bin, value, []string{}, []string{}, waitTime, txSigner)
}

// declared methods
var methods = []eth.MethodSpec{{
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
	Name:    "xp_required",
	Inputs:  []string{"uint256"},
	Outputs: []string{"uint256"},
}, {
	Name:    "summoner",
	Inputs:  []string{"uint256"},
	Outputs: []string{"uint256", "uint256", "uint256", "uint256"},
}}

// declared events
var events = []eth.EventSpec{{
	Name:      "summoned",
	TopicArgs: []string{"address"},
	DataArgs:  []string{"uint256", "uint256"},
}}

type EventSummoned struct {
	Owner    common.Address
	Class    *big.Int
	Summoner *big.Int
}

type contract struct {
	name         string
	term         ui.Screen
	endpoint     rpc.Endpoint
	contractAddr common.Address
	fromAddr     *common.Address
	txSigner     eth.TxSigner
	methods      map[string]eth.Method
	events       map[string]eth.Event
}

type Contract interface {
	GetEvent(eventName string, out interface{}, logs []types.Log) error
	Summon(class *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	Adventure(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	LevelUp(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error)
	AdventurersLog(summoner *uint256.Int) (*uint256.Int, error)
	BalanceOf(address common.Address) (*uint256.Int, error)
	XpRequired(level *uint256.Int) (*uint256.Int, error)
	Summoner(summoner *uint256.Int) (*uint256.Int, *uint256.Int, *uint256.Int, *uint256.Int, error)
}

func New(term ui.Screen, ep rpc.Endpoint, contractAddr common.Address, fromAddr *common.Address, txSigner eth.TxSigner) (Contract, error) {
	methodsMap := make(map[string]eth.Method, len(methods))
	for _, methodSpec := range methods {
		method, err := eth.NewMethod(term, ep, methodSpec.Name, methodSpec.Inputs, methodSpec.Outputs)
		if err != nil {
			return nil, fmt.Errorf("error in rarity %s method: %w", methodSpec.Name, err)
		}
		methodsMap[methodSpec.Name] = method
	}
	eventsMap := make(map[string]eth.Event, len(events))
	for _, eventSpec := range events {
		event, err := eth.NewEvent(eventSpec.Name, eventSpec.TopicArgs, eventSpec.DataArgs)
		if err != nil {
			return nil, fmt.Errorf("error in rarity %s event: %w", eventSpec.Name, err)
		}
		eventsMap[eventSpec.Name] = event
	}
	return &contract{
		name:         "rarity",
		term:         term,
		endpoint:     ep,
		contractAddr: contractAddr,
		fromAddr:     fromAddr,
		txSigner:     txSigner,
		methods:      methodsMap,
		events:       eventsMap,
	}, nil
}

func (c *contract) GetEvent(eventName string, out interface{}, logs []types.Log) error {
	event, ok := c.events[eventName]
	if !ok {
		return errors.New(fmt.Sprintf("Contract %s event %s not declared", c.name, eventName))
	}
	return event.ParseInto(out, logs)
}

func (c *contract) Summon(class *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["summon"].Send(*c.fromAddr, c.contractAddr, nil, []string{class.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["summon"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *contract) Adventure(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["adventure"].Send(*c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["adventure"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *contract) LevelUp(summoner *uint256.Int, waitTime time.Duration) (string, *eth.TxReceipt, error) {
	hash, receipt, err := c.methods["level_up"].Send(*c.fromAddr, c.contractAddr, nil, []string{summoner.ToBig().String()}, waitTime, c.txSigner)
	if err != nil {
		return hash, nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["level_up"].Name(), err)
	}
	return hash, receipt, nil
}

func (c *contract) AdventurersLog(summoner *uint256.Int) (*uint256.Int, error) {
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

func (c *contract) BalanceOf(address common.Address) (*uint256.Int, error) {
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

func (c *contract) XpRequired(level *uint256.Int) (*uint256.Int, error) {
	result, unpacked, err := c.methods["xp_required"].Call(c.fromAddr, c.contractAddr, nil, []string{level.ToBig().String()})
	if err != nil {
		return nil, fmt.Errorf("Contract %s method: %s, err: %w\n", c.name, c.methods["xp_required"].Name(), err)
	}
	if len(unpacked) == 0 {
		// it can happen if call sent to wrong or missing contract
		return nil, errors.New(fmt.Sprintf("Contract %s method: %s. No result: %s returned from node", c.name, c.methods["xp_required"].Name(), string(result)))
	}
	return unpacked[0].ToUint256()
}

func (c *contract) Summoner(summoner *uint256.Int) (*uint256.Int, *uint256.Int, *uint256.Int, *uint256.Int, error) {
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
