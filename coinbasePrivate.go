package main

import (
	"fmt"
	"github.com/preichenberger/go-coinbasepro"
	"github.com/shopspring/decimal"
	"log"
	"strings"
)

// Config structure
// Using Coinbase Sandbox for Development because broke
//var (
//	SandboxApiConfig coinbasepro.ClientConfig = coinbasepro.ClientConfig{
//		BaseURL:    "https://api-public.sandbox.pro.coinbase.com",
//		Key:        "API_KEY",
//		Passphrase: "API_PASSPHRASE",
//		Secret:     "API_SECRET"
//	}
//)

func NewCBProClient(config *coinbasepro.ClientConfig) *coinbasepro.Client {
	client := coinbasepro.NewClient()
	client.UpdateConfig(config)
	return client
}

func getFormattedBalance(c *coinbasepro.Client) string {
	var res strings.Builder
	res.WriteString("-- Your Balance --\n")
	accounts, err := c.GetAccounts()
	if err != nil {
		log.Fatal(err)
	}
	for _, a := range accounts {
		balanceDecimal, err := decimal.NewFromString(a.Balance)
		if err != nil {
			log.Fatal(err)
		}
		res.WriteString(fmt.Sprintf("%-4s: %v\n", a.Currency, balanceDecimal.Round(5)))
	}
	return res.String()
}