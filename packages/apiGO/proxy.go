package apiGO

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func (Proxy *Proxys) GetProxys(uselist bool, list []string) {
	Proxy.Proxys = []string{}
	if uselist {
		Proxy.Proxys = append(Proxy.Proxys, list...)
	} else {
		file, err := os.Open("proxys.txt")
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				Proxy.Proxys = append(Proxy.Proxys, scanner.Text())
			}
		}
	}
}

func (Proxy *Proxys) CompRand() string {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(10 * time.Millisecond)
	return Proxy.Proxys[rand.Intn(len(Proxy.Proxys))]
}

func (Proxy *Proxys) Setup() {
	Proxy.Used = make(map[string]bool)
	for _, proxy := range Proxy.Proxys {
		Proxy.Used[proxy] = false
	}
}

func (Proxy *Proxys) RandProxy() string {
	for _, proxy := range Proxy.Proxys {
		if !Proxy.Used[proxy] {
			Proxy.Used[proxy] = true
			return proxy
		}
	}

	Proxy.Setup()

	return ""
}

func (Bearers *MCbearers) GenSocketConns(Proxy ReqConfig) (pro []Proxys) {
	var Accs [][]Info
	var incr int
	var use int
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(ProxyByte)

	for _, Acc := range Bearers.Details {
		if len(Accs) == 0 {
			Accs = append(Accs, []Info{
				Acc,
			})
		} else {
			if incr == 3 {
				incr = 0
				use++
				Accs = append(Accs, []Info{})
			}
			Accs[use] = append(Accs[use], Acc)
		}
		incr++
	}

	var wg sync.WaitGroup
	for _, Accs := range Accs {
		wg.Add(1)
		go func(Accs []Info) {
			var user, pass, ip, port string
			auth := strings.Split(Proxy.Proxys.RandProxy(), ":")
			ip, port = auth[0], auth[1]
			if len(auth) > 2 {
				user, pass = auth[2], auth[3]
			}
			req, err := proxy.SOCKS5("tcp", fmt.Sprintf("%v:%v", ip, port), &proxy.Auth{
				User:     user,
				Password: pass,
			}, proxy.Direct)
			if err == nil {
				conn, err := req.Dial("tcp", "api.minecraftservices.com:443")
				if err == nil {
					pro = append(pro, Proxys{
						Accounts: Accs,
						Conn:     tls.Client(conn, &tls.Config{RootCAs: roots, InsecureSkipVerify: true, ServerName: "api.minecraftservices.com"}),
					})
				}
			}
			wg.Done()
		}(Accs)
	}

	wg.Wait()
	return
}
