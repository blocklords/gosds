// The sdk/db is the internal key-value database to track the subscription status
// of each smartcontract.
package db

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/blocklords/gosds/env"
	"github.com/blocklords/gosds/static"
	"github.com/blocklords/gosds/topic"
	"github.com/cockroachdb/pebble"
)

type KVM struct {
	topicFilter *topic.TopicFilter
	db          *pebble.DB
}

// the name of the envrionment variable that keeps the path for the database.
const dbPathName = "LOCAL_DB_PATH"

func OpenKVM(topicFilter *topic.TopicFilter) (*KVM, error) {
	if !env.Exists(dbPathName) {
		return nil, fmt.Errorf("missing '%s' envrionment variable", dbPathName)
	}

	db, err := pebble.Open(dbPathName, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	return &KVM{db: db, topicFilter: topicFilter}, nil
}

// Closes the key-value database.
// Its intended to be used with defer.
func (kvm *KVM) Close() {
	kvm.db.Close()
}

// Database keeps the topic filter of the subscribed data.
// Return it.
func (kvm *KVM) TopicFilter() *topic.TopicFilter { return kvm.topicFilter }

// Block Timestamp of the smartcontract on the client side.
func (kvm *KVM) KeyBlockTimestamp(key *static.SmartcontractKey) []byte {
	topicString := kvm.topicFilter.ToString()
	keyString := string(*key)

	return []byte(fmt.Sprintf("%s_%s_sub_block_timestamp", topicString, keyString))
}

// Topic string of the smartcontract on the client side.
func (kvm *KVM) KeyTopicString(key *static.SmartcontractKey) []byte {
	topicString := kvm.topicFilter.ToString()
	keyString := string(*key)

	return []byte(fmt.Sprintf("%s_%s_sub_block_timestamp", topicString, keyString))
}

// Returns the cached block timestamp of the smartcontract.
// If the smartcontract doesn't exist in the database, then it returns 0
func (kvm *KVM) GetBlockTimestamp(key *static.SmartcontractKey) uint64 {
	dbKey := kvm.KeyBlockTimestamp(key)

	bytes, closer, err := kvm.db.Get(dbKey)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	timestamp := binary.BigEndian.Uint64(bytes)

	if err := closer.Close(); err != nil {
		log.Fatal(err)
		return 0
	}

	return timestamp
}

// Returns the cached topics string of the smartcontract.
// If the smartcontract doesn't exist in the database, then it returns the empty string.
func (kvm *KVM) GetTopicString(key *static.SmartcontractKey) string {
	dbKey := kvm.KeyTopicString(key)

	bytes, closer, err := kvm.db.Get(dbKey)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	topicString := string(bytes)

	if err := closer.Close(); err != nil {
		log.Fatal(err)
		return ""
	}

	return topicString
}

// Sets the block timestamp for the given smartcontract.
// If it fails to set it, then returns an error.
func (kvm *KVM) SetBlockTimestamp(key *static.SmartcontractKey, blockTimestamp uint64) error {
	// prepare the value
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, blockTimestamp)

	// prepare the key
	dbKey := kvm.KeyBlockTimestamp(key)

	// store the data
	err := kvm.db.Set(dbKey, bytes, pebble.Sync)
	return err
}

// Sets the block timestamp for the given smartcontract.
// If it fails to set it, then returns an error.
func (kvm *KVM) SetTopicString(key *static.SmartcontractKey, topicString string) error {
	dbKey := kvm.KeyBlockTimestamp(key)

	// store the data
	err := kvm.db.Set(dbKey, []byte(topicString), pebble.Sync)
	return err
}
