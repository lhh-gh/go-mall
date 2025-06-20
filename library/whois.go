package library

import (
	"context"
	"encoding/json"
	"github/lhh-gh/go-mall/comon/logger"
	"github/lhh-gh/go-mall/comon/util/httptool"
)

// 对接 ipwhois.io 的Lib
// Documentation: https://ipwhois.io/documentation

type WhoisLib struct {
	ctx context.Context
}

func NewWhoisLib(ctx context.Context) *WhoisLib {
	return &WhoisLib{ctx: ctx}
}

type WhoisIpDetail struct {
	Ip            string  `json:"ip"`
	Success       bool    `json:"success"`
	Type          string  `json:"type"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	City          string  `json:"city"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEu          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`
}

func (whois *WhoisLib) GetHostIpDetail() (*WhoisIpDetail, error) {
	log := logger.New(whois.ctx)

	httpStatusCode, respBody, err := httptool.Get(
		whois.ctx, "https://www.baidu.com",
		httptool.WithHeaders(map[string]string{
			"User-Agent": "curl/7.77.0",
		}),
	)
	if err != nil {
		log.Error("whois request error", "err", err, "httpStatusCode", httpStatusCode)
		return nil, err
	}
	reply := new(WhoisIpDetail)
	json.Unmarshal(respBody, reply)

	return reply, nil
}
