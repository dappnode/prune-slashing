package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
)

type SlashingProtection struct {
	Metadata struct {
		InterchangeFormatVersion string `json:"interchange_format_version"`
		GenesisValidatorsRoot    string `json:"genesis_validators_root"`
	} `json:"metadata"`
	Data []PubkeyData `json:"data"`
}

type PubkeyData struct {
	Pubkey       string `json:"pubkey"`
	SignedBlocks []struct {
		Slot        string `json:"slot"`
		SigningRoot string `json:"signing_root"`
	} `json:"signed_blocks"`
	SignedAttestations []struct {
		SourceEpoch string `json:"source_epoch"`
		TargetEpoch string `json:"target_epoch"`
		SigningRoot string `json:"signing_root"`
	} `json:"signed_attestations"`
}

func main() {
	// make sure it is called with 3 args
	if len(os.Args) != 5 {
		fmt.Println("Usage: slashing-prune --source-path <source-path> --target-path <target-path>")
		os.Exit(1)
	}

	sourcePath := ""
	targetPath := ""
	// make sure --target-path and --source-path are set
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--source-path" {
			sourcePath = os.Args[i+1]
		}
		if os.Args[i] == "--target-path" {
			targetPath = os.Args[i+1]
		}
	}
	if sourcePath == "" {
		fmt.Println("source path not provided")
		os.Exit(1)
	}
	if targetPath == "" {
		fmt.Println("target path not provided")
		os.Exit(1)
	}

	// make sure --source-path exists
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		fmt.Printf("source path %s does not exist", sourcePath)
		os.Exit(1)
	}

	fmt.Printf("pruning slashing: %s\n", sourcePath)

	// prune --source-path into --target-path
	err := pruneSlashing(sourcePath, targetPath)
	if err != nil {
		panic(fmt.Sprintf("pruneSlashing(%q) failed with %v", sourcePath, err))
	}

	// make sure the target path file exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("target path %s does not exist", targetPath))
	}

	fmt.Printf("pruning slashing complete: %s\n", targetPath)
}

// Prune the slashing protection of a given source path into a target path
func pruneSlashing(sourcePath string, targetPath string) error {
	// Read the file
	file, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Unmarshal the file into a SlashingProtection struct
	var slashingProtection SlashingProtection
	err = json.Unmarshal(file, &slashingProtection)
	if err != nil {
		return err
	}

	// Create another struct to hold the pruned data
	var prunedSlashingProtection SlashingProtection
	prunedSlashingProtection.Metadata = slashingProtection.Metadata

	for _, pubkeyData := range slashingProtection.Data {

		var prunedData PubkeyData
		prunedData.Pubkey = pubkeyData.Pubkey

		if len(pubkeyData.SignedBlocks) > 10 {
			sort.Slice(pubkeyData.SignedBlocks, func(i, j int) bool {
				return pubkeyData.SignedBlocks[i].Slot < pubkeyData.SignedBlocks[j].Slot
			})
			prunedData.SignedBlocks = pubkeyData.SignedBlocks[len(pubkeyData.SignedBlocks)-10:]
		} else {
			prunedData.SignedBlocks = pubkeyData.SignedBlocks
		}

		if len(pubkeyData.SignedAttestations) > 10 {
			sort.Slice(pubkeyData.SignedAttestations, func(i, j int) bool {
				return pubkeyData.SignedAttestations[i].SourceEpoch < pubkeyData.SignedAttestations[j].SourceEpoch
			})
			prunedData.SignedAttestations = pubkeyData.SignedAttestations[len(pubkeyData.SignedAttestations)-10:]
		} else {
			prunedData.SignedAttestations = pubkeyData.SignedAttestations
		}

		prunedSlashingProtection.Data = append(prunedSlashingProtection.Data, prunedData)
	}

	// Marshal the struct into a json file
	json, err := json.Marshal(prunedSlashingProtection)
	if err != nil {
		return err
	}

	// Write the json file
	err = ioutil.WriteFile(targetPath, json, 0644)
	if err != nil {
		return err
	}

	return nil
}
