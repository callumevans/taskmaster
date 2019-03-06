package main

import (
	"github.com/speps/go-hashids"
	"math/rand"
	"time"
)

var hash *hashids.HashID
const hashAlphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	hashData := hashids.NewData()
	hashData.Alphabet = hashAlphabet

	hash, _ = hashids.NewWithData(hashData)
}

func GenerateId() string {
	id, _ := hash.Encode([]int{int(time.Now().UnixNano()), rand.Int()})
	return id
}
