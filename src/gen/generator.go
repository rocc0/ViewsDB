/*
 * Copyright (c) 2018.
 */

package gen

import (
	"math/rand"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//SeqLength of the string that will be returned
type SeqLength struct {
	Length int
}

//Generate returns a random seq of symbols
func (g SeqLength) Generate() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, g.Length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
