package checker

import (
	"fmt"

	"github.com/mikeder/gopinger/internal/minestat"
)

// MinecraftServerChecker checks the status of a Minecraft server.
type MinecraftServerChecker struct {
	host     string
	port     string
	interval int
}

// NewMinecraft returns a pointer to an instance of MinecraftServerChecker.
func NewMinecraft(host, port string, interval int) *MinecraftServerChecker {
	return &MinecraftServerChecker{host: host, port: port, interval: interval}
}

// PerformCheck runs a check against a Minecraft server.
func (mc *MinecraftServerChecker) PerformCheck() error {
	minestat.Init(mc.host, mc.port)
	if minestat.Online {
		fmt.Printf("Server is online running version %s with %s out of %s players.\n", minestat.Version, minestat.Current_players, minestat.Max_players)
		fmt.Printf("Message of the day: %s\n", minestat.Motd)
		fmt.Printf("Latency: %s\n", minestat.Latency)
	} else {
		fmt.Println("Server is offline!")
	}
	return nil
}
