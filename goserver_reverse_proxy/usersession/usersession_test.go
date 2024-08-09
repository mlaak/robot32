package usersession

import (
	"testing"
)

func TestGetUser(t *testing.T) {

	r_ression_id := "123abc"
	ip := "123.123.223.224"

	Getwd := func() (string, error) {
		return "/server/goserver_reverse_proxy", nil
	}

	fileReaderA := func(fnam string) ([]byte, error) {
		if fnam != "/server/working_data/sessions/123abc.txt" {
			t.Fatalf("Error, looking for session file in wrong place. is" + fnam + " but should be /server/working_data/sessions/123abc.txt")
		}
		return []byte("github,g123,m@m.com"), nil
	}

	fileReaderB := func(fnam string) ([]byte, error) {
		return []byte("google,g123,m@m.com"), nil
	}

	fileReaderC := func(fnam string) ([]byte, error) {
		return []byte("ip,ipus123.123.223.224,m@m.com"), nil
	}
	fileReaderD := func(fnam string) ([]byte, error) {
		return []byte("ip,ipus223.123.223.224,m@m.com"), nil
	}
	fileReaderE := func(fnam string) ([]byte, error) {
		return []byte("whatever,ipus223.123.223.224,m@m.com"), nil
	}

	_, usertype, userid := GetUserByCookieIP(0, r_ression_id, ip, Getwd, fileReaderA)
	if usertype != UserTypeGithub {
		t.Fatalf("TestGetUser Error A1 wrong user type")
	}
	if userid != "g123" {
		t.Fatalf("TestGetUser Error A2 wrong user id")
	}

	_, usertype, userid = GetUserByCookieIP(0, r_ression_id, ip, Getwd, fileReaderB)
	if usertype != UserTypeGoogle {
		t.Fatalf("TestGetUser Error B1 wrong user type")
	}
	if userid != "g123" {
		t.Fatalf("TestGetUser Error B2 wrong user id")
	}

	_, usertype, userid = GetUserByCookieIP(0, r_ression_id, ip, Getwd, fileReaderC)
	if usertype != UserTypeCaptchad {
		t.Fatalf("TestGetUser Error C1 wrong user type")
	}
	if userid != "ipus123.123.223.224" {
		t.Fatalf("TestGetUser Error C2 wrong user id")
	}

	_, usertype, userid = GetUserByCookieIP(0, r_ression_id, ip, Getwd, fileReaderD)
	if usertype != UserTypeUnverified { //user capchad but wrong IP
		t.Fatalf("TestGetUser Error D1 wrong user type")
	}
	if userid != "123.123.223.224" {
		t.Fatalf("TestGetUser Error D2 wrong user id")
	}

	_, usertype, userid = GetUserByCookieIP(0, r_ression_id, ip, Getwd, fileReaderE)
	if usertype != UserTypeUnverified {
		t.Fatalf("TestGetUser Error E1 wrong user type")
	}
	if userid != "123.123.223.224" {
		t.Fatalf("TestGetUser Error E2 wrong user id")
	}

}
