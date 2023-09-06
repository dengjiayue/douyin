package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	got := Init("../../../configs/gateway.yaml")
	fmt.Printf("got=%#v\n", got)
}
