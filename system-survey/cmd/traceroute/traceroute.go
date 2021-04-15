package traceroute

import (
	"github.com/aeden/traceroute"
	"time"
)

func Trace(in Input) (Result, error) {
	result := Result{StartTime: time.Now().Unix()}

	opts := &traceroute.TracerouteOptions{}
	opts.SetPort(in.DestinationPort)
	opts.SetMaxHops(in.MaxHops)
	opts.SetRetries(in.Retries)
	opts.SetTimeoutMs(in.MaxLatency)

	res, err := traceroute.Traceroute(in.DestinationIP, opts)
	if err != nil {
		return Result{}, err
	}

	for _, h := range res.Hops {
		result.Hops = append(result.Hops, hop{
			Success:     h.Success,
			Address:     h.AddressString(),
			Host:        h.Host,
			N:           h.N,
			ElapsedTime: h.ElapsedTime,
			TTL:         h.TTL,
		})
	}

	return result, nil
}
