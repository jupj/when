package main

import (
	"strings"
	"testing"
)

func TestZones(t *testing.T) {
	expectedZones := []string{
		"Africa/Dar_es_Salaam",
		"America/Toronto",
		"Asia/Hong_Kong",
		"Europe/Helsinki",
		"UTC",
	}

	got := map[string]bool{}
	for _, zs := range strings.Split(zoneData, "\n") {
		got[zs] = true
	}

	for _, ez := range expectedZones {
		if !got[ez] {
			t.Errorf("No %s zone", ez)
		}
	}
}
