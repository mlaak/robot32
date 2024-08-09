package ratelimiter

import (
	"testing"
	"time"
)

func TestRate(t *testing.T) {

	r := NewRate(67*time.Second, 3, 100)
	r.AddRequest(0, "john")
	r.AddRequest(0, "janar")
	r.AddRequest(0, "janar")
	r.AddRequest(0, "janar")

	if r.IsRequestLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 1 - request limit should not yet be broken")
	}

	r.AddRequest(0, "janar")
	if !r.IsRequestLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 2 - request limit should already be broken")
	}
	if r.IsBytesLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 3 - bytes limit should not yet be broken")
	}
	if r.IsRequestLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 4 - request limit should not yet be broken for John")
	}

	r.Addbytes(0, "john", 60)
	if r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 5 - request limit should not yet be broken for John")
	}

	r.Addbytes(0, "john", 60)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 6 - request limit should already be broken for John")
	}
	if r.IsBytesLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 7 - request limit should not yet be broken for Janar")
	}

	r.AddRequest(0, "janar")
	if !r.IsRequestLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 8 - request limit should now be broken for Janar")
	}

	now := time.Now()
	r.ResetIfTime(0, "john", now)
	if r.IsRequestLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 9 - limit should be reset for john")
	}
	if r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 10 - limit should be reset for john")
	}
	if r.IsRequestLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 11 - limit should be reset for john")
	}
	if !r.IsRequestLimitBroken(0, "janar") {
		t.Fatalf("TestRate Fail 12 - limit should not be reset for janar")
	}

	r.Addbytes(0, "john", 120)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 13 - bytes limit should be broken for john")
	}

	secondTime := now.Add(15 * time.Second)
	r.ResetIfTime(0, "john", secondTime)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 14 - bytes limit should still be broken for john (no reset)")
	}

	thirdTime := now.Add(66 * time.Second)
	r.ResetIfTime(0, "john", thirdTime)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 15 - bytes limit should still be broken for john (no reset)")
	}

	fourthTime := now.Add(71 * time.Second)
	r.ResetIfTime(0, "john", fourthTime)
	if r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 16 - bytes limit should no longer be broken for john (reset happened)")
	}

	r.Addbytes(0, "john", 60)
	if r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 17 - bytes limit should not be broken for john")
	}

	r.Addbytes(0, "john", 60)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 18 - bytes limit should now be broken for john")
	}

	fifthTime := fourthTime.Add(15 * time.Second)
	r.ResetIfTime(0, "john", fifthTime)
	if !r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 19 - bytes limit should still be broken for john (no reset)")
	}

	sixthTime := fourthTime.Add(71 * time.Second)
	r.ResetIfTime(0, "john", sixthTime)
	if r.IsBytesLimitBroken(0, "john") {
		t.Fatalf("TestRate Fail 20 - bytes limit should no longer be broken for john (reset happened)")
	}
}

func TestRateLimiter(t *testing.T) {
	//TODO: Write tests
}
