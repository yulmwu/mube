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
