package ebtree_v2

import (
	"fmt"
	"github.com/ethereum/go-ethereum"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fmt.Println(ethereum.GetValueFromDefaultPath("thread", "pretasknum"))
}
