package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"strconv"
	"strings"
	"sort"
	"time"
)

var PopularCoins = map[string]bool{
	"BTC": true,
	"ETH": true,
	"XRP": false,
	"LTC": true,
	"BCH": true,
	"USDC": true,
	"EOS": true,
	"JPY": true,
}

type CBRate struct {
	base string
	rateMap map[string]float64
	fetchTime time.Time
}

func (c CBRate) String() string {
	var res strings.Builder
	res.WriteString("\033[38;5;242m" + c.base + "\033[m Rate\n")
	res.WriteString(c.fetchTime.Format(time.StampMilli) + "\n")
	sortedKeys := sortedStringKeys(c.rateMap)
	for _, v := range sortedKeys {
		res.WriteString(fmt.Sprintf("%-4s: \033[38;5;36m%15.3f\033[m\n", v, c.rateMap[v]))
	}
	res.WriteString("\n")
	return res.String()
}

func sortedStringKeys(m map[string]float64) []string {
	keyArr := make([]string, 0)
	for k:= range m {
		keyArr = append(keyArr, k)
	}
	sort.Strings(keyArr)
	return keyArr
}

func NewCBRate(coin string) *CBRate {
	c := CBRate{coin, map[string]float64{}, time.Now()}
	c.RefreshRate()
	return &c
}

func (c *CBRate) RefreshRate() {
	resp, err := http.Get("https://api.coinbase.com/v2/exchange-rates?currency=" + c.base)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	c.fetchTime = time.Now()
	var currBodyJson interface{}
	err = json.Unmarshal(body, &currBodyJson)
	if err != nil {
		log.Fatalln(err)
	}
	msg := currBodyJson.(map[string]interface{})
	c.rateMap = parseToRate(msg)
}

func parseToRate(msg map[string]interface{}) map[string]float64 {
	data := msg["data"].(map[string]interface{})
	base := data["currency"].(string)
	rawPairs := data["rates"].(map[string]interface{})
	pairs := make(map[string]float64)
	for k, v := range rawPairs {
		if include, ok := PopularCoins[k]; include && ok && k != base {
			floatVal, err := strconv.ParseFloat(v.(string), 64)
			if err != nil {
				log.Fatalln(err)
			}
			pairs[k] = floatVal
		}
	}
	return pairs
}
