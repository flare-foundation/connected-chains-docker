package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Check func(ctx context.Context, client *http.Client, cfg Config) error

const jsonRPCVersion = "2.0"

var registry = map[string]Check{
	// BTC/Dogecoin
	"blockdownload":   checkBlockDownload,
	"txindex":         checkTxIndex,
	"connectioncount": checkConnectionCount,

	// Ripple
	"nodesynced":   checkNodeServerState,
	"peercount":    checkPeerCount,
	"serverstatus": checkServerStatus,
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func doRPC(ctx context.Context, client *http.Client, cfg Config, method string) (json.RawMessage, error) {
	logger := slog.Default()
	body, err := json.Marshal(rpcRequest{JSONRPC: jsonRPCVersion, ID: "hc", Method: method, Params: []any{}})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.NodeURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if cfg.NodeUser != "" {
		req.SetBasicAuth(cfg.NodeUser, cfg.NodePass)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("node unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("node returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if cfg.Debug {
		logger.Debug("rpc response", "method", method, "body", string(data))
	}

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  any             `json:"error"`
	}

	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, fmt.Errorf("invalid JSON from node: %w", err)
	}

	if envelope.Error != nil {
		return nil, fmt.Errorf("node RPC error: %v", envelope.Error)
	}

	return envelope.Result, nil
}

func checkBlockDownload(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "getblockchaininfo")
	if err != nil {
		return err
	}

	var info struct {
		InitialBlockDownload bool `json:"initialblockdownload"`
	}

	if err := json.Unmarshal(result, &info); err != nil {
		return fmt.Errorf("parse getblockchaininfo: %w", err)
	}

	if info.InitialBlockDownload {
		return fmt.Errorf("node is still syncing (initialblockdownload=true)")
	}

	return nil
}

func checkTxIndex(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "getindexinfo")
	if err != nil {
		return err
	}

	var info struct {
		TxIndex struct {
			Synced bool `json:"synced"`
		} `json:"txindex"`
	}

	if err := json.Unmarshal(result, &info); err != nil {
		return fmt.Errorf("parse getindexinfo: %w", err)
	}

	if !info.TxIndex.Synced {
		return fmt.Errorf("txindex is not synced")
	}

	return nil
}

func checkConnectionCount(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "getconnectioncount")
	if err != nil {
		return err
	}

	var count int
	if err := json.Unmarshal(result, &count); err != nil {
		return fmt.Errorf("parse getconnectioncount: %w", err)
	}
	if count < cfg.MinConnections {
		return fmt.Errorf("not enough connections: %d < %d", count, cfg.MinConnections)
	}

	return nil
}

func checkServerStatus(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "ping")
	if err != nil {
		return err
	}

	var info struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(result, &info); err != nil {
		return fmt.Errorf("parse ping response: %w", err)
	}
	if info.Status != "success" {
		return fmt.Errorf("unexpected status: %q", info.Status)
	}

	return nil
}

func checkNodeServerState(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "server_info")
	if err != nil {
		return err
	}

	var info struct {
		State struct {
			ServerState string `json:"server_state"`
		} `json:"info"`
	}

	if err := json.Unmarshal(result, &info); err != nil {
		return fmt.Errorf("parse server_state response: %w", err)
	}

	validServerStates := map[string]bool{
		"full":       true,
		"validating": true,
		"proposing":  true,
	}

	if !validServerStates[info.State.ServerState] {
		return fmt.Errorf("unexpected server_state: %q", info.State.ServerState)
	}

	return nil
}

func checkPeerCount(ctx context.Context, client *http.Client, cfg Config) error {
	result, err := doRPC(ctx, client, cfg, "server_info")
	if err != nil {
		return err
	}

	var info struct {
		State struct {
			Peers int `json:"peers"`
		} `json:"info"`
	}

	if err := json.Unmarshal(result, &info); err != nil {
		return fmt.Errorf("parse server_state response: %w", err)
	}

	if info.State.Peers < cfg.MinConnections {
		return fmt.Errorf("not enough peers: %d < %d", info.State.Peers, cfg.MinConnections)
	}

	return nil
}
