package types

import (
	"fmt"
	"testing"
)

func TestStatus(t *testing.T) {
	var s Status

	s.Set(Status_IsLogin)
	fmt.Println("after set,", Status_IsLogin, "=", s.Is(Status_IsLogin))
	if !s.Is(Status_IsLogin) {
		t.FailNow()
	}

	s.Unset(Status_IsLogin)
	fmt.Println("after unset", Status_IsLogin, "=", s.Is(Status_IsLogin))
	if s.Is(Status_IsLogin) {
		t.FailNow()
	}

	s.Set(Status_IsFirstlogin)
	fmt.Println("after set,", Status_IsFirstlogin, "=", s.Is(Status_IsFirstlogin))
	if !s.Is(Status_IsFirstlogin) {
		t.FailNow()
	}

	s.Unset(Status_IsFirstlogin)
	fmt.Println("after unset", Status_IsFirstlogin, "=", s.Is(Status_IsFirstlogin))
	if s.Is(Status_IsFirstlogin) {
		t.FailNow()
	}
}
