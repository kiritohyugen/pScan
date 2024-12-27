package scan_test

import (
	"fmt"
	"kiritohyugen/cobra/pScan/scan"
	"net"
	"strconv"
	"testing"
)

func TestStateString(t *testing.T) {

	ps := scan.PortState{}

	if ps.Open.String() != "closed" {
		t.Errorf("Expected %q, got %q instead\n", "closed", ps.Open.String())
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Expected %q, got %q instead\n", "open", ps.Open.String())
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{
			name:        "OpenPort",
			expectState: "open",
		},
		{
			name:        "ClosedPort",
			expectState: "closed",
		},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	ports := []int{}

	for _, tc := range testCases {
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}

		defer ln.Close()

		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}

		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		ports = append(ports, port)
		fmt.Println("Added port:", port) // Print statement to display the added port

		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}
	res := scan.Run(hl, ports)
	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if res[0].NotFound {
		t.Errorf("Expected host %q to be found\n", host)
	}

	//why
	if len(res[0].PortStates) != 2 {
		t.Fatalf("Expected 2 port states, got %d instead\n", len(res[0].PortStates))
	}
}

func TestRunHostNotFound(t *testing.T) {

	// 	Create an instance of the scan.HostsList and add the host 389.389.389.389 to it.
	//  Name resolution on this host should fail unless you have it on your DNS:

	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("Expected 1 results, got %d instead\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host %q, got %q instead\n", host, res[0].Host)
	}
	if !res[0].NotFound {
		t.Errorf("Expected host %q NOT to be found\n", host)
	}
	if len(res[0].PortStates) != 0 {
		t.Fatalf("Expected 0 port states, got %d instead\n", len(res[0].PortStates))
	}
}
