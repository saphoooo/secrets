package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	storedSecrets := make(map[string]string)
	storedSecrets["secret1"] = "secret_value"
	storedSecrets["secret2"] = ""

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		log.Fatalf("Failed to read: %v", scanner.Err())
	}

	text := secrets{}
	err := json.Unmarshal(scanner.Bytes(), &text)
	if err != nil {
		log.Println("failed to unmasharl the object")
	}
	data := make(map[string]output)
	for _, key := range text.Secrets {
		if value, ok := storedSecrets[key]; ok {
			if value != "" {
				data[key] = output{value, "null"}
			} else {
				data[key] = output{value, "empty value for the requested secret"}
			}
		} else {
			data[key] = output{value, "could not fetch the secret"}
		}
	}
	jsonString, _ := json.Marshal(data)
	fmt.Println(string(jsonString))
}

type secrets struct {
	Version string   `json:"version,omitempty"`
	Secrets []string `json:"secrets"`
}

type output struct {
	Value string `json:"value"`
	Err   string `json:"error"`
}
