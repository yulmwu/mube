## Concept

**Mube** is a project that originated from the idea of a 'lightweight container orchestration platform', aiming to simplify container orchestration within a specific domain (Domain-Specific).

To achieve this, it is based on the Kubernetes architecture, but is designed to simplify Kubernetes's complex configurations and networking.
It retains the fundamental roles of the Kubernetes Control Plane and Worker Nodes while removing or simplifying the associated auxiliary components.

The architecture of Mube consists of the following main components (items in *italics* are currently in the planning stage and may not yet be implemented or could be removed/changed in the future):

**Control Plane**

* Mube API Server
* *Controller Manager*
* *Scheduler*
* *Etcd* (*may be removed or replaced in the future*)

**Worker Node**

* Mube Node Agent (mubelet)
* *Mube Proxy* (*mube-proxy*)

Mube aims to support OCI(Open Container Initiative) compliant container runtimes, and is planned to use Containerd by default.
Compatibility with container runtimes will be maintained through Kubernetes's CRI(Container Runtime Interface).

In addition, it plans to implement standard interfaces such as CNI (Container Network Interface) and CSI (Container Storage Interface) to simplify networking and storage concerns.

For more details, please refer to the architecture diagram below.

![Architecture](assets/architecture.png)

> [!NOTE]
>
> This architecture represents a planned design and may be subject to change at any time. It will be updated accordingly whenever changes occur. For now, it closely follows the Kubernetes architectural model.

## Bootstrap (MVP)

The repository includes two runnable components:

- `Mube API Server` (`cmd/mube-apiserver`)
- `mubelet` (`cmd/mubelet`)

### Design Change (Static Node Registry)

For MVP simplicity, node network information is hardcoded in API Server YAML config.

- API Server loads `configs/apiserver.yaml`.
- Only nodes listed in this file can send heartbeat.
- API Server keeps authoritative node `ip` and `port` from this file.

```yaml
listenAddress: ":8080"
nodes:
  - name: "mube-node-1"
    ip: "127.0.0.1"
    port: 10250
```

### API Endpoints (MVP)

- `GET /healthz`
- `GET /api/v1/nodes` (includes `state`: `Unknown|Ready|NotReady`)
- `POST /api/v1/nodes/heartbeat`

### mubelet Endpoints (MVP)

- `GET /healthz`
- `GET /status`

### Running the API Server and mubelet:

```bash
cp .env.apiserver.example .env.apiserver
cp .env.mubelet.example .env.mubelet
make tidy
make run-apiserver
```

```bash
make run-mubelet
```

Check node registration:

```bash
curl -s http://127.0.0.1:8080/api/v1/nodes | jq
```

### Environment Variables

**API Server**:

- `MUBE_APISERVER_CONFIG` (default `configs/apiserver.yaml`)
- `MUBE_APISERVER_SHUTDOWN_TIMEOUT` (default `5s`)
- `MUBE_NODE_NOTREADY_TIMEOUT` (default `30s`)

State behavior:

- `Unknown`: node is registered but no heartbeat yet within timeout window
- `Ready`: heartbeat received within timeout
- `NotReady`: no heartbeat beyond timeout (including never-heartbeated nodes after timeout)

**mubelet**:

- `MUBELET_NODE_NAME` (default hostname)
- `MUBELET_RUNTIME` (default `containerd`)
- `MUBELET_VERSION` (default `0.1.0`)
- `MUBELET_CAPACITY` (default `10`)
- `MUBE_API_SERVER` (default `http://127.0.0.1:8080`)
- `MUBELET_HEARTBEAT_INTERVAL` (default `10s`)
- `MUBELET_HEALTH_LISTEN` (default `:10250`)
- `MUBELET_SHUTDOWN_TIMEOUT` (default `5s`)
