package main

import (
	"io/ioutil"
	"testing"
)

func TestPruneSlashing(t *testing.T) {
	sourcePath := "./slashing_files/slashing_protection.json"
	targetPath := "./slashing_files/slashing_protection_prune.json"
	err := pruneSlashing(sourcePath, targetPath)
	if err != nil {
		t.Fatalf("pruneSlashing(%q) failed with %v", sourcePath, err)
	}

	// compare ./slashing_files/slashing_protection_prune.json with ./slashing_files/expected_slashing_protectionjson
	// if they are not the same, then the pruneSlashing function is not working correctly
	expectedPath := "./slashing_files/expected_slashing_protection_prune.json"
	expectedFile, err := ioutil.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("expected file %s does not exist", expectedPath)
	}
	targetFile, err := ioutil.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("target file %s does not exist", targetPath)
	}
	// compare the two files byte by byte
	if len(expectedFile) != len(targetFile) {
		t.Fatalf("expected file %s and target file %s are not the same length", expectedPath, targetPath)
	}
	for i := 0; i < len(expectedFile); i++ {
		if expectedFile[i] != targetFile[i] {
			t.Fatalf("expected file %s and target file %s are not the same", expectedPath, targetPath)
		}
	}

}
