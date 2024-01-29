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
	if v := ca.Get("gi"); v != nil {
		t.Errorf("\t%s %s", failed, WRAnswer)
	}
}

