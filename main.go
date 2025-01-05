package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ja-howell/playlister/videoclient"
)

const apiKeyFilepath = "API_KEY"

func main() {
	err := run()
	if err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func run() error {
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API Key: %w", err)
	}
	client := videoclient.New(apiKey)
	resp, err := client.GetResponse()
	if err != nil {
		return fmt.Errorf("failed to get a response: %w", err)
	}

	config, err := newConfig("config.json")
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	// database, err := newDatabase("../db.json")
	// if err != nil {
	// 	return fmt.Errorf("failed to create database: %w", err)
	// }

	fmt.Println(resp)

	_ = resp
	_ = config
	// _ = database
	return nil
}

func getAPIKey() (string, error) {
	f, err := os.Open(apiKeyFilepath)
	if err != nil {
		return "", fmt.Errorf("failed to open API Key file: %w", err)
	}

	defer f.Close()

	apiKey, err := io.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("failed to read API Key: %w", err)
	}

	return string(apiKey), nil
}
