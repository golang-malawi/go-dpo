// Package dpo provides functionality for interacting with DPO Group's payment gateway from Go applications.
// Currently the module only supports performing payments through DPOs verify token workflow.
//
// # Usage: User Agent
//
// You are recommended to set the user agent for the client to some string that identifies your application.
//
//	clientToken := os.Getenv("DPO_TOKEN")
//	client := dpo.NewClient(clientToken, true)
//	client.SetUserAgent("Example User Agent")
//
// # Usage: Error Handling
//
// The dpo package exposes errors that are thrown from DPO API.
package dpo
