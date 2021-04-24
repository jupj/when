package main

import "testing"

func TestZones(t *testing.T) {
	zs, err := readTZNames()
	if err != nil {
		t.Fatal(err)
	}

	const expectedLen = 607
	if len(zs) != expectedLen {
		t.Errorf("got %d zones, expected %d", len(zs), expectedLen)
	}

	for _, zone := range zs {
		if zone == "Europe/Helsinki" {
			return
		}
	}
	t.Errorf("No Europe/Helsinki zone")
}
