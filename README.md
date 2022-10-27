
<p align="center">
  Large Kubetunnel Logo centered
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
* You are able to connect to other services in the namespace from the local process.

You can read more about it [here](https/dcode.tech).
<p align="center">
  <img src="./images/how_it_works.svg" alt="How It Works"/>
</p>

---

## Quick Start

A few quick ways to start using KubeTunnel

Kubetunnel CLI can be installed through multiple channels.

## Install with Homebrew

Install Kubetunnel CLI from [Kubetunnel tap](https://github.com/kubetunnel/homebrew-tap) with [Homebrew](https://brew.sh) by running:

```bash
brew tap we-dcode/tap
brew install kubetunnel
```

## Install with apt

```bash
apt install kubetunnel
```

# Getting started with Kubetunnel CLI

Once you installed the KubeTunnel CLI, you can verify it's working by running:

```bash
kubetunnel --help
```

## Doing XX with Kubetunnel CLI


```bash
kubetunnel whatever
```

## More flags and options to try

Here are some flags that you might find useful:

- `--severity-threshold=low|medium|high|critical`

  Only report vulnerabilities of provided level or higher.

- `--json`

  Prints results in JSON format.

- `--all-projects`

  Auto-detect all projects in working directory

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
