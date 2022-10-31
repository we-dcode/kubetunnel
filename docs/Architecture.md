# KubeTunnel Architecture

KubeTunnel uses a Kubernetes Operator, KubeTunnel CustomResources and a local CLI to forward traffic to and from your local workstation and the cluster. 

The **operator** pod is composed of two different components:

1. A manager component that handles the KubeTunnel Custom Resources and deployments. 
2. A REST API that handles requests from the KubeTunnel deployments to patch the application services. This API will change the application service's labelSelectors to move traffic to the original application or to the tunnel depending on the connection of the tunnel to your station. This ensures that if your local process is down or that the KubeTunnel CRD still exists in the cluster but you aren't using KubeTunnel currently the service returns to forward traffic to the original application pod.

Each **KubeTunnel server** is composed of a deployment with the following components:

1. An FRP Server container that includes configuration from the KubeTunnel CLI.
2. A KubeTunnel Go server that checks that the FRP Server is connected to your local process. If it is connected, it sends a request to the operator to change the Kubernetes service's labelSelector to it, and if not, it sends the request to change the labelSelector back to the original pod.

The **KubeTunnel CLI** achieves the two way connection by doing the following:

From you to the cluster:
1. Forwards all the cluster services to your local station and modifies your `hosts` file.


From the cluster to you:
1. Creates a KubeTunnel CR which includes the tunnel from the local process to the cluster.
2. Changes the application's service labelSelectors to point to this tunnel, "stealing" the traffic from the Kubernetes deployment and moving it to your process.
