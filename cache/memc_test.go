package cache

import (
	"testing"
	"time"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)
const WRAnswer = " <WRONG ANSWER>"

func TestCacheFirst(t *testing.T) {
	ca := New()
	ca.Set("mama", "papa", 10*time.Minute)

	go func() {
		ca.Set("gi", "ig", 1*time.Second)
	}()

	if v := ca.Get("mama").String(); v != "papa" {
		t.Errorf(WRAnswer)
	}
	<-time.After(2 * time.Second)
	if v := ca.Get("gi"); v != (Item{}) {
		t.Errorf("\t%s %s", failed, WRAnswer)
	}
}

func TestCacheDelete(t *testing.T) {
	ca := New()
	ca.Set("float", 28.232342, 1*time.Hour)
	ca.Set("int", 56565656, 1*time.Hour)
	ca.Set("byte", []byte("Hi tommy"), 1*time.Hour)

	var d1, err = ca.GetD("float").Float64()
	if  err != nil || d1 != 28.232342 {
		t.Errorf("\t%s %v", failed, err)
	}
	var d2, err2 = ca.Get("byte").Byte()
	if err2 != nil || string(d2) != "Hi tommy"{
		t.Errorf("\t%s %v", failed, err)
	}
	if ca.Count() != 2 {
		t.Errorf("\t%s Count is not correct", failed)
	}
}
