package situation
import (
    "testing"
)

//Not much to test here at the moment, but lets leave the skeleton here
func TestRequestContext(t *testing.T){
	context := NewRequestContext()
	context.SetIporid("123:ab.c")
	
	s := context.GetIporid() 
	if(s!="123:ab.c"){
		t.Fatalf("TestRequestContext Fail 1 - expected to get back 123:ab.c but go "+s)
	}
}