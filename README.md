

<p align="center" styles="background-color: white">
<img src="./assets/logo.svg" width="300"><br>
  by <a href="https://dcode.tech">Dcode</a>
</p>

# KubeTunnel: Develop locally while being connected to Kubernetes.

![Large GIF of Kubetunnel example](https://aaa) <br><br><br>
Website: [https://www.dcode.tech/](https://www.dcode.tech)  
Slack: [Discuss](https://we-dcode.slack.com/archives/C047WAUR41M)

**With KubeTunnel:**

* You run one service locally using your favorite IDE
* You run the rest of your microservices in Kubernetes, not limited to resources and compute power.

**This gives developers:**

* A fast local dev loop, with no waiting for a container build / push / deploy
* Ability to use their favorite local tools (IDE, debugger, etc.)
* Ability to run large-scale applications that can't run locally

---

## How It Works
When you select a service and a local process to tunnel, Kubetunnel launches a pod on your namespace, and changes the service to move all traffic to this pod. This pod forwards all traffic to your local proccess through a secured tunnel. The CLI then forwards all your Kubernetes services in your namespace to localhost. 
This creates a two way connection between your local process and the cluster meaning:

* All other pods now connect to your local process instead of the original pod in the cluster.
* You are able to connect to other services in the namespace from your local process.

You can read more about it [here](docs/Architecture.md).
<p align="center">
  <img src="./images/how_it_works.svg" alt="How It Works"/>
</p>

---

Kubetunnel CLI can be installed through multiple channels.

### Linux/Mac

There are 2 methods to install the CLI on Linux/Mac.

1. Install Kubetunnel CLI from [Kubetunnel tap](https://github.com/kubetunnel/homebrew-tap) with [Homebrew](https://brew.sh) by running:

```bash
brew tap we-dcode/tap
brew install kubetunnel
```

2. Download the KubeTunnel CLI with the latest binary in our [releases page](https://github.com/we-dcode/kubetunnel/releases/latest). 

### Windows

Download the latest CLI binary in our [releases page](https://github.com/we-dcode/kubetunnel/releases/latest)


# Getting started with Kubetunnel CLI

Things to consider before you start:

* For this quickstart guide, your Kubernetes cluster is assumed to be already up and running. Before you proceed with the KubeTunnel installation, make sure you check the supported versions.
* Make sure your user has `cluster-admin` permissions for the initial installation of the operator component. For the tunnels themselves, this is not needed.
* The operator needs network access to each tunnel. If your namespaces deny ingress and egress traffic, please create NetworkPolicies to enable traffic between them as explained [here](docs/Network.md).
* You will need local administrator privileges to create each tunnel as the KubeTunnel client modifies the local hosts file to include the cluster services.

Once you installed the KubeTunnel CLI, you can verify it's working by running:

```bash
kubetunnel --help
```

For each of the following commands, you can run --help for more options.


1. To install the operator and CRD, run the following command:

```bash
kubetunnel install 
```
At this point, the KubeTunnel Kubernetes Operator is successfully installed. Once the KubeTunnel Operator pod is running, you are able to start tunneling processes to your cluster. 

2. For each service you want to tunnel, run the following command:

```bash
  sudo -E kubetunnel create-tunnel -p '8080:80' svc_name
```

This command waits for the local process to be available with the process you want to tunnel. When the process is up, it is tunneled to cluster and the application service is switched to forward traffic to it. If your local process becomes unavailable, the service is switched back to the original pod.

##  Autocomplete with Kubetunnel CLI

KubeTunnel supports completion for multiple shells.  
For autocomplete for your shell run the following command:

```bash
kubetunnel completion --help
```

See all the available commands and options by running:

# Getting support

If you need support using KubeTunnel CLI, please [join our Slack channel](https://we-dcode.slack.com/archives/C047WAUR41M).

Please leave issues for any error or bug that you encounter.

# Known Limitations

* The current KubeTunnel version can only tunnel a single service per workstation. In the future, we will add support for multiple services.

# Contributing

If you are an external contributor, before working on any contributions, please first [contact us](https://dcode.tech) to discuss the issue or feature request with us.

---

Made with 💙 by Dcode
