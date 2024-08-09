package limits

import (
	s "grp/usersession"
	"testing"
)

func TestLimits(t *testing.T) {
	SetupRateLimiters()

	a1, a2 := GetLimitersForPathAndUserType(0, "/experts/general", s.UserTypeUnverified)
	b1, b2 := GetLimitersForPathAndUserType(0, "/experts/general", s.UserTypeCaptchad)
	c1, c2 := GetLimitersForPathAndUserType(0, "/experts/general", s.UserTypeGithub)
	d1, d2 := GetLimitersForPathAndUserType(0, "/experts/general", s.UserTypeGoogle)
	e1, e2 := GetLimitersForPathAndUserType(0, "/experts/general", "whatever")

	f1, f2 := GetLimitersForPathAndUserType(0, "/recieved_images/a.jpg", s.UserTypeUnverified)
	g1, g2 := GetLimitersForPathAndUserType(0, "/recieved_images/a.jpg", s.UserTypeCaptchad)
	h1, h2 := GetLimitersForPathAndUserType(0, "/recieved_images/a.jpg", s.UserTypeGithub)
	i1, i2 := GetLimitersForPathAndUserType(0, "/recieved_images/a.jpg", s.UserTypeGoogle)

	if a1.GetNr() != Code498Limiter.GetNr() || a2.GetNr() != Code498Limiter.GetNr() {
		t.Fatalf("TestLimits Fail A")
	}
	if b1.GetNr() != ExpertIPLimiter.GetNr() || b2.GetNr() != ExpertAIPLimiter.GetNr() {
		t.Fatalf("TestLimits Fail B")
	}
	if c1.GetNr() != ExpertUSRLimiter.GetNr() || c2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail C")
	}
	if d1.GetNr() != ExpertUSRLimiter.GetNr() || d2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail D")
	}
	if e1.GetNr() != ExpertIPLimiter.GetNr() || e2.GetNr() != ExpertAIPLimiter.GetNr() {
		t.Fatalf("TestLimits Fail E")
	}

	if f1.GetNr() != ImageIPLimiter.GetNr() || f2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail F")
	}
	if g1.GetNr() != ImageIPLimiter.GetNr() || g2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail G")
	}
	if h1.GetNr() != ImageUSRLimiter.GetNr() || h2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail H")
	}
	if i1.GetNr() != ImageUSRLimiter.GetNr() || i2.GetNr() != UnLimiter.GetNr() {
		t.Fatalf("TestLimits Fail I")
	}

}
