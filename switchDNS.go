package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func main() {
	log.Println("hello world")
	DNSs := []string{
		"114.114.114.114",
		"8.8.8.8",
		"8.8.4.4",
		"218.102.23.228",
		"211.136.192.6",
		"223.5.5.5",
		"168.126.63.1",
		"168.126.63.2",
		"168.95.1.1",
		"168.95.192.1",
		"203.80.96.9",
		"61.10.0.130",
		"61.10.1.130",
		"208.67.222.222",
		"208.67.220.220",
		"202.14.67.4",
		"203.80.96.10",
		"202.14.67.14",
		"198.153.194.1",
		"198.153.192.1",
		"112.106.53.22",
		"168.126.63.1",
		"168.95.192.1",
		"198.153.194.1",
		"210.2.4.8",
		"203.80.96.9",
		"220.67.240.221",
		"84.200.69.80",
		"81.218.119.11",
		"180.76.76.76",
		"119.29.29.29",
	}
	ch_avbl := make(chan bool)
	// ch_avbl==1说明通道可用
	for _, dns_server := range DNSs {
		go control(dns_server, ch_avbl)
		break
	}
	time.Sleep(300 * time.Second)

}

func control(dns_server string, ch_avbl chan bool) {
	Nintendo_Dl_Url := "ctest-dl-lp1.cdn.nintendo.net"
	Nintendo_Up_Url := "ctest-ul-lp1.cdn.nintendo.net"
	t1, ip1, err1 := resolve_speed(Nintendo_Dl_Url, dns_server)
	t2, ip2, err2 := resolve_speed(Nintendo_Up_Url, dns_server)
	if err1 == nil && err2 == nil {
		log.Printf("%15s%15s\n", dns_server, t1+t2)
		// _ = ip1 + ip2
		log.Printf("%15s%15s\n", ip1, ip2)
		transfer_speed(ip1, ip2, ch_avbl)
	} else {
		log.Printf("%15s%15s\n", dns_server, "error")
	}
}

func resolve_speed(url, dns_server string) (time.Duration, string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(url), dns.TypeA)
	m.RecursionDesired = true

	r, t, _ := c.Exchange(m, net.JoinHostPort(dns_server, "53"))
	if r == nil {
		return t, "", errors.New("dns error")
	} else if r.Rcode != dns.RcodeSuccess {
	}
	var ip string
	for _, ansa := range r.Answer {
		switch ansb := ansa.(type) {
		case *dns.A:
			ip = ansb.A.String()
		}
	}
	return t, ip, nil
}

func transfer_speed(dl_ip string, up_ip string, ch_avbl chan bool) time.Duration {
	req1, err := http.NewRequest(http.MethodGet, "http://"+dl_ip+"/30m", nil)
	if err != nil {
		log.Printf("*** REQ FAILED")
	}
	req1.Header.Set("user-agent", "Nintendo NX")
	req1.Header.Set("Accept-Encoding", "gzip, deflate")
	req1.Header.Set("Accept", "*/*")
	req1.Header.Set("Connection", "keep-alive")
	req1.Host = "ctest-dl-lp1.cdn.nintendo.net"

	body := []byte(strings.Repeat(" ", 1024*1024))
	req2, err := http.NewRequest(http.MethodPost, "http://"+up_ip+"/1m", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("*** REQ FAILED:\n%s", err.Error())
	}
	req2.Header.Set("User-Agent", "Nintendo NX")
	req2.Header.Set("Accept-Encoding", "gzip, deflate")
	req2.Header.Set("Accept", "*/*")
	req2.Header.Set("Connection", "keep-alive")
	req2.Header.Set("Content_Type", "application/x-www-form-urlencoded")
	req2.Host = "ctest-ul-lp1.cdn.nintendo.net"

	client := &http.Client{}
	start_time := time.Now()
	resp, err := client.Do(req2)
	if resp.StatusCode != 200 {
		log.Printf("*** RESP FAILED:\n%s", resp.Status)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("error")
	}
	bodySize := len(body)
	log.Println(resp)
	log.Println(bodySize)
	elapsed_time := time.Since(start_time) / time.Millisecond

	resp, err = client.Do(req2)
	if resp.StatusCode != 200 {
		log.Printf("*** RESP FAILED:\n%s", resp.Status)
	}
	return elapsed_time
}
