package apiGO

import (
	"crypto/tls"
	"time"
)

func PingMC() float64 {
	var pingTimes float64
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	defer conn.Close()
	for i := 0; i < 10; i++ {
		recv := make([]byte, 4096)
		time1 := time.Now()
		conn.Write([]byte("PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(recv)
		pingTimes += float64(time.Since(time1).Milliseconds())
	}
	return float64(pingTimes/10000) * 5000
}
