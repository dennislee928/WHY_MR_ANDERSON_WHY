package network

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

// PacketCapture handles packet capture using libpcap
// 封包捕獲器，使用 libpcap 進行實時封包捕獲
type PacketCapture struct {
	handle       *pcap.Handle
	interfaceName string
	snapshotLen  int32
	promiscuous  bool
	timeout      time.Duration
	logger       *logrus.Logger
	mu           sync.RWMutex
	active       bool
	packetSource *gopacket.PacketSource
}

// NewPacketCapture creates a new packet capture instance
func NewPacketCapture(interfaceName string, snapshotLen int32, promiscuous bool, timeout time.Duration, logger *logrus.Logger) *PacketCapture {
	if logger == nil {
		logger = logrus.New()
	}

	return &PacketCapture{
		interfaceName: interfaceName,
		snapshotLen:   snapshotLen,
		promiscuous:   promiscuous,
		timeout:       timeout,
		logger:        logger,
	}
}

// Start starts packet capture
func (pc *PacketCapture) Start() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if pc.active {
		return fmt.Errorf("capture already active")
	}

	// 打開網路介面
	handle, err := pcap.OpenLive(
		pc.interfaceName,
		pc.snapshotLen,
		pc.promiscuous,
		pc.timeout,
	)
	if err != nil {
		pc.logger.Errorf("Failed to open interface %s: %v", pc.interfaceName, err)
		return fmt.Errorf("failed to open interface: %w", err)
	}

	pc.handle = handle
	pc.packetSource = gopacket.NewPacketSource(handle, handle.LinkType())
	pc.active = true

	pc.logger.Infof("Packet capture started on %s (snapshot: %d, promiscuous: %v)",
		pc.interfaceName, pc.snapshotLen, pc.promiscuous)

	return nil
}

// Stop stops packet capture
func (pc *PacketCapture) Stop() error {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	if !pc.active {
		return nil
	}

	if pc.handle != nil {
		pc.handle.Close()
		pc.handle = nil
	}

	pc.active = false
	pc.packetSource = nil

	pc.logger.Infof("Packet capture stopped on %s", pc.interfaceName)
	return nil
}

// SetBPFFilter sets a BPF filter for packet capture
func (pc *PacketCapture) SetBPFFilter(filter string) error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if pc.handle == nil {
		return fmt.Errorf("capture not started")
	}

	if err := pc.handle.SetBPFFilter(filter); err != nil {
		pc.logger.Errorf("Failed to set BPF filter '%s': %v", filter, err)
		return fmt.Errorf("failed to set BPF filter: %w", err)
	}

	pc.logger.Infof("BPF filter set: %s", filter)
	return nil
}

// CapturePackets captures packets and processes them
func (pc *PacketCapture) CapturePackets(handler func(*PacketInfo) error) error {
	pc.mu.RLock()
	if !pc.active || pc.packetSource == nil {
		pc.mu.RUnlock()
		return fmt.Errorf("capture not active")
	}
	packetSource := pc.packetSource
	pc.mu.RUnlock()

	for packet := range packetSource.Packets() {
		pc.mu.RLock()
		if !pc.active {
			pc.mu.RUnlock()
			break
		}
		pc.mu.RUnlock()

		info := pc.parsePacket(packet)
		if info != nil {
			if err := handler(info); err != nil {
				pc.logger.Errorf("Packet handler error: %v", err)
			}
		}
	}

	return nil
}

// parsePacket parses a packet into PacketInfo
func (pc *PacketCapture) parsePacket(packet gopacket.Packet) *PacketInfo {
	info := &PacketInfo{
		Timestamp: packet.Metadata().Timestamp,
		Length:    packet.Metadata().Length,
	}

	// 解析網路層
	if networkLayer := packet.NetworkLayer(); networkLayer != nil {
		switch networkLayer.LayerType() {
		case layers.LayerTypeIPv4:
			ipv4 := networkLayer.(*layers.IPv4)
			info.SourceIP = ipv4.SrcIP.String()
			info.DestIP = ipv4.DstIP.String()
			info.Protocol = ipv4.Protocol.String()

		case layers.LayerTypeIPv6:
			ipv6 := networkLayer.(*layers.IPv6)
			info.SourceIP = ipv6.SrcIP.String()
			info.DestIP = ipv6.DstIP.String()
			info.Protocol = ipv6.NextHeader.String()
		}
	}

	// 解析傳輸層
	if transportLayer := packet.TransportLayer(); transportLayer != nil {
		switch transportLayer.LayerType() {
		case layers.LayerTypeTCP:
			tcp := transportLayer.(*layers.TCP)
			info.SourcePort = int(tcp.SrcPort)
			info.DestPort = int(tcp.DstPort)
			info.Protocol = "TCP"

			// TCP flags
			if tcp.SYN {
				info.Flags = append(info.Flags, "SYN")
			}
			if tcp.ACK {
				info.Flags = append(info.Flags, "ACK")
			}
			if tcp.FIN {
				info.Flags = append(info.Flags, "FIN")
			}
			if tcp.RST {
				info.Flags = append(info.Flags, "RST")
			}
			if tcp.PSH {
				info.Flags = append(info.Flags, "PSH")
			}

		case layers.LayerTypeUDP:
			udp := transportLayer.(*layers.UDP)
			info.SourcePort = int(udp.SrcPort)
			info.DestPort = int(udp.DstPort)
			info.Protocol = "UDP"

		case layers.LayerTypeICMPv4:
			info.Protocol = "ICMP"

		case layers.LayerTypeICMPv6:
			info.Protocol = "ICMPv6"
		}
	}

	// 解析應用層
	if appLayer := packet.ApplicationLayer(); appLayer != nil {
		payload := appLayer.Payload()
		if len(payload) > 0 {
			// 只保存前 100 字節的 payload
			maxLen := 100
			if len(payload) < maxLen {
				maxLen = len(payload)
			}
			info.PayloadSample = payload[:maxLen]
		}
	}

	return info
}

// PacketInfo contains parsed packet information
type PacketInfo struct {
	Timestamp     time.Time
	Length        int
	SourceIP      string
	DestIP        string
	SourcePort    int
	DestPort      int
	Protocol      string
	Flags         []string
	PayloadSample []byte
}

// GetInterfaceInfo returns information about a network interface
func GetInterfaceInfo(interfaceName string) (*InterfaceInfo, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find devices: %w", err)
	}

	for _, device := range devices {
		if device.Name == interfaceName {
			return &InterfaceInfo{
				Name:        device.Name,
				Description: device.Description,
				Addresses:   formatAddresses(device.Addresses),
				Flags:       device.Flags,
			}, nil
		}
	}

	return nil, fmt.Errorf("interface not found: %s", interfaceName)
}

// ListNetworkInterfaces lists all available network interfaces
func ListNetworkInterfaces() ([]*InterfaceInfo, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, fmt.Errorf("failed to find devices: %w", err)
	}

	var interfaces []*InterfaceInfo
	for _, device := range devices {
		interfaces = append(interfaces, &InterfaceInfo{
			Name:        device.Name,
			Description: device.Description,
			Addresses:   formatAddresses(device.Addresses),
			Flags:       device.Flags,
		})
	}

	return interfaces, nil
}

// InterfaceInfo contains information about a network interface
type InterfaceInfo struct {
	Name        string
	Description string
	Addresses   []string
	Flags       uint32
}

// formatAddresses formats interface addresses
func formatAddresses(addresses []pcap.InterfaceAddress) []string {
	var addrs []string
	for _, addr := range addresses {
		if addr.IP != nil {
			addrs = append(addrs, addr.IP.String())
		}
	}
	return addrs
}

// GetPacketStatistics returns packet capture statistics
func (pc *PacketCapture) GetPacketStatistics() (*pcap.Stats, error) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if pc.handle == nil {
		return nil, fmt.Errorf("capture not started")
	}

	stats, err := pc.handle.Stats()
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return stats, nil
}

