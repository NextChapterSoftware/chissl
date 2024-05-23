package settings

import (
	"reflect"
	"testing"
)

func TestRemoteDecode(t *testing.T) {
	//test table
	for i, test := range []struct {
		Input   string
		Output  Remote
		Encoded string
	}{
		{
			"8080->80 ",
			Remote{
				UserAddress: "8080->80",
				RemotePort:  "80",
				RemoteHost:  "127.0.0.1",
				LocalHost:   "0.0.0.0",
				LocalPort:   "8080",
				Reverse:     true,
			},
			"0.0.0.0:8080->127.0.0.1:80",
		},
		{
			"80:localhost->8080:remotehost",
			Remote{
				UserAddress: "80:localhost->8080:remotehost",
				LocalHost:   "localhost",
				LocalPort:   "80",
				RemoteHost:  "remotehost",
				RemotePort:  "8080",
				Reverse:     true,
			},
			"localhost:80->remotehost:8080",
		},
		{
			" 80:10.1.2.3->8080:localhost ",
			Remote{
				UserAddress: "80:10.1.2.3->8080:localhost",
				LocalHost:   "10.1.2.3",
				LocalPort:   "80",
				RemoteHost:  "localhost",
				RemotePort:  "8080",
				Reverse:     true,
			},
			"10.1.2.3:80->localhost:8080",
		},
		{
			" 8080->80:localhost",
			Remote{
				UserAddress: "8080->80:localhost",
				LocalHost:   "0.0.0.0",
				LocalPort:   "8080",
				RemoteHost:  "localhost",
				RemotePort:  "80",
				Reverse:     true,
			},
			"0.0.0.0:8080->localhost:80",
		},
		{
			" 8080:127.0.0.1->80",
			Remote{
				UserAddress: "8080:127.0.0.1->80",
				LocalHost:   "127.0.0.1",
				LocalPort:   "8080",
				RemoteHost:  "127.0.0.1",
				RemotePort:  "80",
				Reverse:     true,
			},
			"127.0.0.1:8080->127.0.0.1:80",
		},
	} {
		//expected defaults
		expected := test.Output
		if expected.LocalHost == "" {
			expected.LocalHost = "0.0.0.0"
		}

		//compare
		got, err := DecodeRemote(test.Input)
		if err != nil {
			t.Fatalf("decode #%d '%s' failed: %s", i+1, test.Input, err)
		}
		if !reflect.DeepEqual(got, &expected) {
			t.Fatalf("decode #%d '%s' expected\n  %#v\ngot\n  %#v", i+1, test.Input, expected, got)
		}
		if e := got.Encode(); test.Encoded != e {
			t.Fatalf("encode #%d '%s' expected\n  %#v\ngot\n  %#v", i+1, test.Input, test.Encoded, e)
		}
	}
}
