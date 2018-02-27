package main

import "testing"


func TestRandSeq(t *testing.T) {
	r := RandSeq(3)
	if len(r) < 3 || r == "" {
		t.Fail()
	}
	r = RandSeq(8)
	if len(r) < 8 || r == "" {
		t.Fail()
	}
	r = RandSeq(11)
	if len(r) < 11 || r == "" {
		t.Fail()
	}
	r = RandSeq(20)
	if len(r) < 20 || r == "" {
		t.Fail()
	}
}