package config

import (
	"decentralodge/tool"
	"encoding/json"
	"errors"
	"github.com/libp2p/go-libp2p-core/crypto"
	"log"
	mrand "math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var cpath string
var kpath string

func init() {
	cpath, _ = os.Getwd()
	kpath = cpath + "\\values\\priv_key"
	cpath += "\\values\\config.json"
}

type Config struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	RandomSeed int64  `json:"randomSeed"`
	// Private Key
	PrvKey crypto.PrivKey `json:"prvKey"`
	// Bootstrap Node
	BootstrapNode string `json:"bootstrapNode"`
}

func (c *Config) Save() bool {

	pk, _ := crypto.MarshalPrivateKey(c.PrvKey)
	err := tool.WriteFile(pk, kpath)
	if err != nil {
		return false
	}
	c.PrvKey = nil
	b, _ := json.Marshal(*c)
	err = tool.WriteFile(b, cpath)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (c *Config) Load() *Config {
	data, err := tool.LoadFile(cpath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	pk, err := tool.LoadFile(kpath)
	c.PrvKey, _ = crypto.UnmarshalPrivateKey(pk)
	return c
}

func (c *Config) AddrString() string {
	// "/ip4/127.0.0.1/tcp/2000"
	return "/ip4/" + c.IP + "/tcp/" + c.Port
}

// New Get a Configuration
func New(username, pwd, ipAddr string, rs int64, bn string) (*Config, error) {
	strs := strings.Split(ipAddr, ":")
	if strs[0] == "" {
		strs[0] = "127.0.0.1"
	}
	// 获取私钥
	r := mrand.New(mrand.NewSource(rs))
	prvKey, _, err := crypto.GenerateRSAKeyPair(2048, r)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if b := ipFormatCheck(strs); b {
		return &Config{
			Username:      username,
			Password:      pwd,
			IP:            strs[0],
			Port:          strs[1],
			RandomSeed:    rs,
			PrvKey:        prvKey,
			BootstrapNode: bn,
		}, nil
	} else {
		return nil, errors.New("the format of ipAddr is wrong")
	}
}

func ipFormatCheck(ipAddr []string) bool {
	compile, _ := regexp.Compile("((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)")
	if b := compile.MatchString(ipAddr[0]); b {
		port, err := strconv.ParseInt(ipAddr[1], 10, 64)
		if err != nil || port < 0 || port > 65535 {
			return false
		} else {
			return true
		}
	} else {
		return false
	}
}
