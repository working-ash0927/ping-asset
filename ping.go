package main

import (
	"fmt"
	"net"
	"time"
	"flag"
	"github.com/go-ping/ping"
)

func getunixtime() int64 {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}
	now := time.Now()
	t := now.In(loc)
	a := t.UnixNano() / 1000000
	return a
}

func getIPAddr() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("# Failed to get IP address:", err)
		return "failed to get IP of Agent"
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}

	return "# failed to get IP of Agent"
}

func sendPing(addr string) {
	source := getIPAddr()  // icmp를 호출하는 주체 
	target := addr         // icmp 패킷을 받는 대상
	nowtime := getunixtime()
	
	// Pinger 생성
	pinger, err := ping.NewPinger(target)  // 대상서버 IP 입력
	if err != nil {
		fmt.Printf("ping_test{source=\"%s\", target=\"%s\"} %d %d\n", source, target, 0, nowtime)
		fmt.Printf("# failed to make new pinger: %s\n", err.Error())
		return
	}

	// PING 설정
	pinger.Count = 1             // 1회 PING 전송
	pinger.SetPrivileged(true)   // ICMP 대신 UDP 사용
	pinger.Timeout = 100 * time.Millisecond // 타임아웃 0.1초로 설정

	// 결과 처리 함수 등록
	pinger.OnRecv = func(pkt *ping.Packet) {
		fmt.Printf("# received: %d bytes, time: %v\n", pkt.Nbytes, pkt.Rtt)
	}

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("# %d packets transmitted, %d received, %.2f%% packet loss, average RTT: %v\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss, stats.AvgRtt)

		if stats.PacketsRecv == 0 {
			fmt.Printf("ping_test{source=\"%s\", target=\"%s\"} %d %d\n", source, target, 0, nowtime)
			fmt.Println("# ping test failed or timed out")
		} else{
			fmt.Printf("ping_test{source=\"%s\", target=\"%s\"} %d %d\n", source, target, 1, nowtime)
			fmt.Println("# ping test succeed")
		}
	}

	// PING 실행
	fmt.Printf("# ping test start: %s\n", pinger.Addr())
	err = pinger.Run()
	if err != nil {
		fmt.Println("# ping did not run")
		return
	}
}

func main() {
	var addr string

	// command 상의 변수 할당
	flag.StringVar( &addr, "addr", "", "AWS region") // Ping을 보내려는 서버의 IP(즉, Sensu backend)
	flag.Parse()
	sendPing(addr)
}
