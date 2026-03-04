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

var registry = map[string]Check{
	"blockdownload": checkBlockDownload,
	"txindex":       checkTxIndex,
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

func doRPC(ctx context.Context, client *http.Client, cfg Config, method string) (json.RawMessage, error) {
	logger := slog.Default()
	body, err := json.Marshal(rpcRequest{JSONRPC: "1.0", ID: "hc", Method: method, Params: []any{}})
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	logger.Debug("rpc response", "method", method, "body", string(data))

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
