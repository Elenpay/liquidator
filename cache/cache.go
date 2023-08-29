package cache

import (

	//Log
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	//Nodeguard
	"github.com/Elenpay/liquidator/nodeguard"
	"github.com/allegro/bigcache/v3"
)

type Cache interface {

	//Set a slice of LiquidityRules from nodeguard grpc response
	SetLiquidityRules(string, []nodeguard.LiquidityRule) error
	//Get a slice of nodeguard.LiquidityRules from cache returns map where the key is the channel id
	GetLiquidityRules(string) (map[uint64][]nodeguard.LiquidityRule, error)
}

type BigCache struct {
	cache *bigcache.BigCache
}

// Create a new BigCache
func NewCache() (Cache, error) {

	log.Debug("creating new bigcache instance")

	cache, _ := bigcache.New(context.Background(), bigcache.Config{
		Shards:      1,
		LifeWindow:  0,
		CleanWindow: 0,
		Verbose:     false,
		OnRemove: func(key string, entry []byte) {
			log.Debugf("removing key from cache: %s", key)
		},
		OnRemoveWithMetadata: func(key string, entry []byte, keyMetadata bigcache.Metadata) {
			log.Debugf("removing key from cache: %s", key)

		},
		OnRemoveWithReason: func(key string, entry []byte, reason bigcache.RemoveReason) {
			log.Debugf("removing key from cache: %s", key)
		},
		Logger: log.StandardLogger(),
	})

	return &BigCache{
		cache: cache,
	}, nil

}

// Implement BigCache methods according to the Cache Interface
func (c *BigCache) SetLiquidityRules(nodePubkey string, rules []nodeguard.LiquidityRule) error {

	log.Debugf("setting liquidity rules in cache: %+v", rules)

	//Check if rules are empty
	if len(rules) == 0 {
		log.Debug("no liquidity rules to set in cache")
		return nil
	}

	//Convert rules to bytes
	rulesBytes, err := c.MarshalLiquidityRules(rules)
	if err != nil {
		log.Errorf("error marshalling liquidity rules: %s", err)
		return err
	}

	err = c.cache.Set(nodePubkey, rulesBytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *BigCache) GetLiquidityRules(nodePubkey string) (map[uint64][]nodeguard.LiquidityRule, error) {

	entry, err := c.cache.Get(nodePubkey)
	if err != nil {
		return nil, err
	}

	//Convert bytes to rules
	rules, err := c.UnmarshalLiquidityRules(entry)
	if err != nil {
		return nil, err
	}

	//Convert rules to map
	rulesMap := make(map[uint64][]nodeguard.LiquidityRule)
	for _, rule := range rules {
		rulesMap[rule.ChannelId] = append(rulesMap[rule.ChannelId], rule)
	}

	return rulesMap, nil
}

// Marshals a slice of nodeguard.LiquidityRules to bytes
func (c *BigCache) MarshalLiquidityRules(rules []nodeguard.LiquidityRule) ([]byte, error) {

	rulesBytes, err := json.Marshal(rules)
	if err != nil {
		log.Errorf("error marshalling liquidity rules: %s", err)
		return nil, err
	}

	return rulesBytes, nil
}

// Unmarshals a slice of nodeguard.LiquidityRules from bytes
func (c *BigCache) UnmarshalLiquidityRules(rulesBytes []byte) ([]nodeguard.LiquidityRule, error) {

	var rules []nodeguard.LiquidityRule
	err := json.Unmarshal(rulesBytes, &rules)
	if err != nil {
		log.Errorf("error unmarshalling liquidity rules: %s", err)
		return nil, err
	}

	return rules, nil
}
