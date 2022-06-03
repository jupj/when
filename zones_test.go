package main

import (
	"strings"
	"testing"
)

func TestZones(t *testing.T) {
	zs := strings.Split(zoneData, "\n")

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
