package bridges

import (
	"bytes"
	"encoding/json"
	"fmt"
	hue_control_http "hue-control/internal/http"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

const path = "http://%s/api"

type pairingRequestDto struct {
	Devicetype string `json:"devicetype"`
}

func fromDiscoveredBridge(discoveredBridge DisvoveredBridge) ([]byte, error) {
	pairingRequestDto := pairingRequestDto{
		Devicetype: "light-control#remshams",
	}
	log.Debugf("Pairing request dto: %v", pairingRequestDto)
	return json.Marshal(pairingRequestDto)
}

type pairResultDto struct {
	Username string `json:"username"`
}

type pairSuccessResponseDto struct {
	Success pairResultDto `json:"success"`
}

func parseResponse(body []byte) (*pairSuccessResponseDto, error) {
	var pairSuccessResponseDto []pairSuccessResponseDto
	if strings.Contains(string(body), "error") {
		return nil, fmt.Errorf("Failed to pair")
	}
	err := json.Unmarshal(body, &pairSuccessResponseDto)
	if err != nil {
		log.Error("Failed to pair")
		return nil, err
	}
	return &pairSuccessResponseDto[0], nil
}

type BridgeHttpAdapter struct {
	ip net.IP
}

func InitBridgesHttpAdapter(ip net.IP) BridgeHttpAdapter {
	return BridgeHttpAdapter{
		ip: ip,
	}
}

func (adapter BridgeHttpAdapter) Pair(discoveredBridge DisvoveredBridge) (*Bridge, error) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			pairSuccessResponseDto, err := adapter.doPair(discoveredBridge)
			if err == nil {
				log.Debugf("Paired bridge: %v", pairSuccessResponseDto)
				return &Bridge{
					ip:     adapter.ip,
					apiKey: pairSuccessResponseDto.Success.Username,
				}, nil
			} else {
				log.Errorf("Failed to pair bridge: %v", err)
			}
		}
	}
}

func (adapter BridgeHttpAdapter) doPair(discoveredBridge DisvoveredBridge) (*pairSuccessResponseDto, error) {
	pairingRequestDto, err := fromDiscoveredBridge(InitDiscoverdBridge(adapter, discoveredBridge.Id, discoveredBridge.Ip))
	req, client, cancel, err := hue_control_http.RequestWithTimeout(
		http.MethodPost,
		fmt.Sprintf(path, adapter.ip),
		bytes.NewReader(pairingRequestDto),
		nil,
	)
	defer cancel()
	res, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send pair request")
		return nil, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error("Failed to read pair response")
		return nil, err
	}
	pairSuccessResponseDto, err := parseResponse(body)
	if err != nil {
		return nil, err
	}
	return pairSuccessResponseDto, nil
}
