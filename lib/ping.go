package lib

import (
	"time"

	"github.com/prometheus-community/pro-bing"
)

func Ping(ip string) (time.Duration, error) {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		return 0, err
	}
	pinger.Count = 3
	err = pinger.Run()
	if err != nil {
		return 0, err
	}
	stats := pinger.Statistics()

	return stats.AvgRtt.Round(time.Microsecond), nil
} 
