## Connected Chains Healthcheck

A minimal Go binary/container that exposes a `/readyz` HTTP endpoint for blockchain node readiness probes. Designed to be used alongside any compatible node (Bitcoin, Dogecoin, Litecoin, etc.) in a Kubernetes pod.

### How it works

On each `/readyz` request, the sidecar runs a configurable set of checks against the node's JSON-RPC API. It returns 200 OK when all checks pass, or an error with a message on the first failing check.

### Checks

| Check | RPC Method | Description |
|---|---|---|
| `blockdownload` | `getblockchaininfo` | Passes when `initialblockdownload` is `false`. Default for all chains. |
| `txindex` | `getindexinfo` | Passes when `txindex.synced` is `true`. |
| `connectioncount` | `getconnectioncount` | Passes when the node has at least `MIN_CONNECTIONS` peers. |

### Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `NODE_URL` | yes | — | RPC endpoint, e.g. `http://localhost:8332` |
| `NODE_USER` | no | — | RPC auth username |
| `NODE_PASS` | no | — | RPC auth password |
| `CHECKS` | no | `blockdownload` | Comma-separated list of checks to run (`blockdownload, txindex, connectioncount`) |
| `MIN_CONNECTIONS` | no | `8` | Minimum peer connections required for the `connectioncount` check |


### Running Locally

```bash
go build -o node-healthcheck .
NODE_URL=http://localhost:8332 NODE_USER=user NODE_PASS=secret CHECKS=blockdownload,txindex ./node-healthcheck
```

```bash
# Readiness
curl http://localhost:8080/readyz
```

Example Output:
```json
{"time":"2026-03-04T08:43:48.697Z","level":"INFO","msg":"starting","addr":":8080","checks":["blockdownload"],"node":"http://localhost:8332"}
{"time":"2026-03-04T08:44:01.123Z","level":"DEBUG","msg":"rpc response","method":"getblockchaininfo","body":"{\"result\":{\"initialblockdownload\":false,...}}"}
```

### Kubernetes Sidecar Setup

```yaml
containers:
  - name: bitcoin
    ...
  - name: healthcheck
    image: node-healthcheck:latest
    env:
      - name: NODE_URL
        value: "http://bitcoin:8332"
      - name: NODE_USER
        value: "admin"
      - name: NODE_PASS
        valueFrom:
          secretKeyRef:
            name: node-rpc-secret
            key: password
      - name: CHECKS
        value: blockdownload,txindex
    ports:
      - containerPort: 8080

readinessProbe:
  httpGet:
    path: /readyz
    port: 8080
  initialDelaySeconds: 60
  periodSeconds: 30
  timeoutSeconds: 20
  failureThreshold: 10
```

