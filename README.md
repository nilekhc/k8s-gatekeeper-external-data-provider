# Kubernetes API Provider

The Kubernetes API provider is used for querying Kubernetes resources. This can be a 
good alternative to caching, as the cache can grow larger as the cluster becomes
busier with many resources.

This provider offers an example of how to query Kubernetes objects where a
cluster-scoped view is required to execute a policy. The example here demonstrates
how uniqueness of the ingress host name is evaluated using this provider.

Please refer to the [Quick Start](#quick-start) to set up your local environment and see a demo of
this provider in action.

## Prerequisites

- [ ] [`docker`](https://docs.docker.com/get-docker/)
- [ ] [`helm`](https://helm.sh/)
- [ ] [`kind`](https://kind.sigs.k8s.io/)
- [ ] [`kubectl`](https://kubernetes.io/docs/tasks/tools/#kubectl)

## Quick Start

1. Run [local-setup.sh](./scripts/local-setup.sh) to setup Gatekeeper and k8s external data provider.

2. Run [demo.sh](./scripts/demo.sh) to run the k8s external data provider demo.

## Installation
Below installation steps helps run the _Kubernetes API Provider_ locally on KinD cluster.
-   Setup KinD cluster
    ```shell
    kind delete cluster --name gatekeeper
    kind create cluster --name gatekeeper
    ```
-   Install Gatekeeper 
    ```shell
    helm install gatekeeper/gatekeeper \
        --name-template=gatekeeper \
        --namespace gatekeeper-system \
        --create-namespace --wait --debug
    ```
-   Locally build provider image install it using helm chart
    ```shell
    ./scripts/generate-tls-cert.sh
    make docker-buildx
    make kind-load-image
    
    kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=gatekeeper-system:default
    
    helm install external-data-provider charts/external-data-provider \
    --set provider.tls.caBundle="$(cat certs/ca.crt | base64 | tr -d '\n\r')" \
    --namespace "${NAMESPACE:-gatekeeper-system}"
    ```
-   Install Constraint and ConstraintTemplate that calls the provider
    ```shell
    kubectl apply -f validation/ingress-ct.yaml
    kubectl apply -f validation/ingress-c.yaml
    ```

## Verification
-   Pre-populate ingress resource
    ```shell
    count=3
    for i in $(seq 1 $count); do
        {
        cat <<EOF | kubectl apply -f -
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
    name: ingress-$i
    annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    spec:
    rules:
    - host: test-$i.kind.local
      http:
      paths:
        - path: /
          pathType: Prefix
          backend:
          service:
          name: test
          port:
          number: 80
    ---
    EOF
        }
    done
    ```
  -   Correct Ingress
      ```shell
      kubectl apply -f ./scripts/correct-ingress.yaml
      ```
      - Request should be allowed
         ```shell
         ingress.networking.k8s.io/correct-ingress-test created
         ```
-   Incorrect Ingress
    ```shell
    kubectl apply -f ./scripts/incorrect-ingress.yaml"
    ```
    - Request should be rejected
       ```shell
       Error from server (Forbidden): error when creating "./scripts/incorrect-ingress.yaml": admission webhook
       "validation.gatekeeper.sh" denied the request: [deny-ingress-with-duplicate-host]
       invalid response: {"errors": [["test-1.kind.local", "Duplicate Ingress host found test-1.kind.local_invalid"]], "responses": [], "status_code": 200, "system_error": ""}
       ```
      
To get started refer to the scripts in [Quick Start](#quick-start) section.