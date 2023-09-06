package my_jwt

import (
	"fmt"
	"testing"
	"time"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(1244567654345)
	if err != nil {
		fmt.Printf("GenToken err: %v", err)
	}
	fmt.Printf("token: %v", token)
	time.Sleep(3 * time.Second)
	claims, err := ParseToken(token)
	if err != nil {
		fmt.Printf("ParseToken err: %v", err)
	}
	fmt.Printf("claims: %v", claims)
}
