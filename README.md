<p align="center">
  <img height="100" src="https://via.placeholder.com/150">
</p>

# KubeTunnel

[KubeTunnel](https://dcode.tech/) helps you develop microservices locally while being connected to your Kubernetes environment.

![Large GIF of Kubetunnel example](https://aaa) <br><br><br>
Website: [https://www.dcode.tech/](https://www.dcode.tech)  
Slack: [Discuss](https://we-dcode.slack.com/archives/C047WAUR41M)

**With KubeTunnel:**

* You run one service locally using your favorite IDE
* You run the rest of your application in Kubernetes, not limited to resources and compute power.

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

Install Kubetunnel CLI from [Kubetunnel tap](https://github.com/kubetunnel/homebrew-tap) with [Homebrew](https://brew.sh) by running:

```bash
brew tap we-dcode/tap
brew install kubetunnel
```

### Windows

Download the latest CLI binary - https://github.com/we-dcode/kubetunnel/releases


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

To install the operator and CRD, run the following command:

```bash
kubetunnel install 
```
At this point, the KubeTunnel Kubernetes Operator is successfully installed. Once the KubeTunnel Operator pod is running, you are able to start tunneling processes to your cluster. 

For each service to tunnel you want to create, run the following command:

```bash
  sudo -E kubetunnel create-tunnel -p '8080:80' svc_name
```

This command does the following things:
1. Create a KubeTunnel custom resource in the cluster which creates kubetunnel server in the cluster.
2. Starts a local frp client to connect to the kubetunnel server with your wanted process on the given local port.
3. Forwards all the namespace services to your localhost and adds their DNS to your `/etc/hosts` file.
4. Starts forwarding traffic to your local process by changing the Kubernetes service to forward to the new frp server pod instead of the app.


##  Autocomplete with Kubetunnel CLI


```bash
kubetunnel completion <zsh/bash/fish/powershell> >> <>
```

[See all the available commands and options](./help/cli-commands) by running `--help`:

```bash
kubetunnel --help
```

# Getting support

If you need support using KubeTunnel CLI, please [join our Slack channel](https://we-dcode.slack.com/archives/C047WAUR41M).

We do not actively monitor GitHub Issues so any issues there may go unnoticed.

# Contributing

If you are an external contributor, before working on any contributions, please first [contact us](https://dcode.tech) to discuss the issue or feature request with us.

If you are contributing to KubeTunnel CLI, see [our contributing guidelines](CONTRIBUTING.md)

For information on how KubeTunnel CLI is implemented, see [our design decisions](help/_about-this-project/README.md).

---

Made with 💙 by Dcode
