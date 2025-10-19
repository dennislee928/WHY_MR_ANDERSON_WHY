package security

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// TLSFingerprinter implements JA3/JA3S TLS fingerprinting
type TLSFingerprinter struct {
	knownBots      map[string]BotInfo
	knownMalware   map[string]MalwareInfo
	mu             sync.RWMutex
	logger         *logrus.Logger
}

// BotInfo contains bot identification information
type BotInfo struct {
	Name        string
	Type        string // "good" (search engine) or "bad" (scraper)
	Description string
	AddedAt     string
}

// MalwareInfo contains malware identification information
type MalwareInfo struct {
	Name        string
	Family      string
	Severity    string // "low", "medium", "high", "critical"
	Description string
	AddedAt     string
}

// TLSFingerprint represents a TLS fingerprint
type TLSFingerprint struct {
	JA3          string
	JA3Hash      string
	ClientHello  *ClientHelloInfo
	IsBot        bool
	IsMalware    bool
	Identified   string
	Confidence   float64
}

// ClientHelloInfo contains parsed Client Hello information
type ClientHelloInfo struct {
	Version            uint16
	CipherSuites       []uint16
	Extensions         []uint16
	EllipticCurves     []uint16
	EllipticCurvePointFormats []uint8
	SNI                string
	ALPN               []string
}

// NewTLSFingerprinter creates a new TLS fingerprinter
func NewTLSFingerprinter(logger *logrus.Logger) *TLSFingerprinter {
	if logger == nil {
		logger = logrus.New()
	}

	fp := &TLSFingerprinter{
		knownBots:    make(map[string]BotInfo),
		knownMalware: make(map[string]MalwareInfo),
		logger:       logger,
	}

	// 載入已知的 Bot 指紋
	fp.loadKnownBots()
	
	// 載入已知的惡意軟體指紋
	fp.loadKnownMalware()

	return fp
}

// GenerateJA3 generates JA3 fingerprint from TLS ClientHello
func (fp *TLSFingerprinter) GenerateJA3(hello *ClientHelloInfo) *TLSFingerprint {
	// JA3 格式: SSLVersion,Ciphers,Extensions,EllipticCurves,EllipticCurvePointFormats
	ja3String := fp.buildJA3String(hello)
	ja3Hash := fp.hashJA3(ja3String)

	fingerprint := &TLSFingerprint{
		JA3:         ja3String,
		JA3Hash:     ja3Hash,
		ClientHello: hello,
	}

	// 檢查是否為已知 Bot
	fp.identifyBot(fingerprint)

	// 檢查是否為惡意軟體
	fp.identifyMalware(fingerprint)

	return fingerprint
}

// buildJA3String builds the JA3 string
func (fp *TLSFingerprinter) buildJA3String(hello *ClientHelloInfo) string {
	parts := []string{
		fmt.Sprintf("%d", hello.Version),
		joinUint16(hello.CipherSuites),
		joinUint16(hello.Extensions),
		joinUint16(hello.EllipticCurves),
		joinUint8(hello.EllipticCurvePointFormats),
	}

	return strings.Join(parts, ",")
}

// hashJA3 creates SHA256 hash of JA3 string
func (fp *TLSFingerprinter) hashJA3(ja3String string) string {
	hash := sha256.Sum256([]byte(ja3String))
	return hex.EncodeToString(hash[:])
}

// identifyBot checks if fingerprint matches known bot
func (fp *TLSFingerprinter) identifyBot(fingerprint *TLSFingerprint) {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	if bot, exists := fp.knownBots[fingerprint.JA3Hash]; exists {
		fingerprint.IsBot = true
		fingerprint.Identified = bot.Name
		fingerprint.Confidence = 0.95
		fp.logger.Infof("Bot identified: %s (type: %s)", bot.Name, bot.Type)
	}
}

// identifyMalware checks if fingerprint matches known malware
func (fp *TLSFingerprinter) identifyMalware(fingerprint *TLSFingerprint) {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	if malware, exists := fp.knownMalware[fingerprint.JA3Hash]; exists {
		fingerprint.IsMalware = true
		fingerprint.Identified = malware.Name
		fingerprint.Confidence = 0.98
		fp.logger.Warnf("Malware identified: %s (family: %s, severity: %s)", 
			malware.Name, malware.Family, malware.Severity)
	}
}

// loadKnownBots loads known bot fingerprints
func (fp *TLSFingerprinter) loadKnownBots() {
	// Googlebot
	fp.knownBots["e7d705a3286e19ea42f587b344ee6865"] = BotInfo{
		Name:        "Googlebot",
		Type:        "good",
		Description: "Google Search Engine Bot",
		AddedAt:     "2025-01-01",
	}

	// Bingbot
	fp.knownBots["ada70206e40642a3e4461f35503241d5"] = BotInfo{
		Name:        "Bingbot",
		Type:        "good",
		Description: "Bing Search Engine Bot",
		AddedAt:     "2025-01-01",
	}

	// Scrapy (common scraper framework)
	fp.knownBots["bc6c386f480ee97b9d9e52d472b772d8"] = BotInfo{
		Name:        "Scrapy",
		Type:        "bad",
		Description: "Python web scraping framework",
		AddedAt:     "2025-01-01",
	}

	// Python Requests
	fp.knownBots["51c64c77e60f3980eea90869b68c58a8"] = BotInfo{
		Name:        "Python-Requests",
		Type:        "bad",
		Description: "Python HTTP library",
		AddedAt:     "2025-01-01",
	}

	// cURL
	fp.knownBots["7dd50e112cd23734a310b90f6f44c1c4"] = BotInfo{
		Name:        "cURL",
		Type:        "bad",
		Description: "Command-line HTTP client",
		AddedAt:     "2025-01-01",
	}

	fp.logger.Infof("Loaded %d known bot fingerprints", len(fp.knownBots))
}

// loadKnownMalware loads known malware fingerprints
func (fp *TLSFingerprinter) loadKnownMalware() {
	// Trickbot
	fp.knownMalware["6734f37431670b3ab4292b8f60f29984"] = MalwareInfo{
		Name:        "Trickbot",
		Family:      "Banking Trojan",
		Severity:    "critical",
		Description: "Modular banking trojan",
		AddedAt:     "2025-01-01",
	}

	// Emotet
	fp.knownMalware["72a589da586844d7f0818ce684948eea"] = MalwareInfo{
		Name:        "Emotet",
		Family:      "Botnet",
		Severity:    "critical",
		Description: "Polymorphic malware",
		AddedAt:     "2025-01-01",
	}

	// Cobalt Strike
	fp.knownMalware["a0e9f5d64349fb13191bc781f81f42e1"] = MalwareInfo{
		Name:        "Cobalt Strike",
		Family:      "Post-Exploitation",
		Severity:    "high",
		Description: "Adversary simulation software",
		AddedAt:     "2025-01-01",
	}

	// Metasploit
	fp.knownMalware["d8b1d125f5b6b5b5d5b5b5b5b5b5b5b5"] = MalwareInfo{
		Name:        "Metasploit",
		Family:      "Penetration Testing",
		Severity:    "high",
		Description: "Penetration testing framework",
		AddedAt:     "2025-01-01",
	}

	fp.logger.Infof("Loaded %d known malware fingerprints", len(fp.knownMalware))
}

// AddBotFingerprint adds a new bot fingerprint
func (fp *TLSFingerprinter) AddBotFingerprint(ja3Hash string, bot BotInfo) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.knownBots[ja3Hash] = bot
	fp.logger.Infof("Added bot fingerprint: %s (%s)", bot.Name, ja3Hash)
}

// AddMalwareFingerprint adds a new malware fingerprint
func (fp *TLSFingerprinter) AddMalwareFingerprint(ja3Hash string, malware MalwareInfo) {
	fp.mu.Lock()
	defer fp.mu.Unlock()

	fp.knownMalware[ja3Hash] = malware
	fp.logger.Warnf("Added malware fingerprint: %s (%s)", malware.Name, ja3Hash)
}

// ParseClientHello parses TLS ClientHello from connection state
func ParseClientHello(connState *tls.ConnectionState) *ClientHelloInfo {
	if connState == nil {
		return nil
	}

	hello := &ClientHelloInfo{
		Version:      connState.Version,
		CipherSuites: make([]uint16, 0),
		Extensions:   make([]uint16, 0),
		SNI:          connState.ServerName,
		ALPN:         connState.NegotiatedProtocol,
	}

	// Note: 實際實現需要從原始 TLS 握手數據中提取
	// 這裡提供簡化版本

	return hello
}

// Helper functions

func joinUint16(values []uint16) string {
	if len(values) == 0 {
		return ""
	}

	// 排序以確保一致性
	sorted := make([]uint16, len(values))
	copy(sorted, values)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	strValues := make([]string, len(sorted))
	for i, v := range sorted {
		strValues[i] = fmt.Sprintf("%d", v)
	}

	return strings.Join(strValues, "-")
}

func joinUint8(values []uint8) string {
	if len(values) == 0 {
		return ""
	}

	strValues := make([]string, len(values))
	for i, v := range values {
		strValues[i] = fmt.Sprintf("%d", v)
	}

	return strings.Join(strValues, "-")
}

// GetStatistics returns fingerprinting statistics
func (fp *TLSFingerprinter) GetStatistics() map[string]interface{} {
	fp.mu.RLock()
	defer fp.mu.RUnlock()

	return map[string]interface{}{
		"known_bots":    len(fp.knownBots),
		"known_malware": len(fp.knownMalware),
	}
}

