package contracts

import (
	"errors"
	"fmt"
	"time"

	"github.com/holiman/uint256"
	"github.com/jaanek/jeth/eth"
	"github.com/jaanek/jeth/rpc"
	"github.com/jaanek/jeth/ui"
	"github.com/ledgerwatch/erigon/common"
)

func DeployRarity(term ui.Screen, ep rpc.Endpoint, fromAddr common.Address, bin []byte, value *uint256.Int, waitTime time.Duration, txSigner eth.TxSigner) (string, *eth.TxReceipt, error) {
	typeNames := []string{}
	values := []string{}
	return eth.Deploy(term, ep, fromAddr, bin, value, typeNames, values, waitTime, txSigner)
}

type MethodSpec struct {
	name    string
	inputs  []string
	outputs []string
}

var rarityMethods = []MethodSpec{{
	name:    "summon",
	inputs:  []string{"uint256"},
	outputs: []string{},
}, {
	name:    "adventure",
	inputs:  []string{"uint256"},
	outputs: []string{},
}, {
	name:    "level_up",
	inputs:  []string{"uint256"},
	outputs: []string{},
}, {
	name:    "adventurers_log",
	inputs:  []string{"uint256"},
	outputs: []string{"uint256"},
}, {
	name:    "balanceOf",
	inputs:  []string{"address"},
	outputs: []string{"uint256"},
}, {
	name:    "summoner",
	inputs:  []string{"uint256"},
	outputs: []string{"uint256", "uint256", "uint256", "uint256"},
}}

type rarityContract struct {
	name         string
	term         ui.Screen
	endpoint     rpc.Endpoint
	contractAddr common.Address
	fromAddr     *common.Address
	txSigner     eth.TxSigner
	methods      map[string]eth.Method
}

type Rarity interface {
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
		method, err := eth.NewMethod(term, ep, methodSpec.name, methodSpec.inputs, methodSpec.outputs)
		if err != nil {
			return nil, fmt.Errorf("error in rarity %s method: %w", methodSpec.name, err)
		}
		methods[methodSpec.name] = method
	}
	return &rarityContract{
		name:         "rarity",
		term:         term,
		endpoint:     ep,
		contractAddr: contractAddr,
		fromAddr:     fromAddr,
		txSigner:     txSigner,
		methods:      methods,
	}, nil
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
