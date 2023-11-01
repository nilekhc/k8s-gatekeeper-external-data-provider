#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail


kind delete cluster --name gatekeeper
kind create cluster --name gatekeeper
helm install gatekeeper/gatekeeper \
    --set enableExternalData=true \
    --name-template=gatekeeper \
    --namespace gatekeeper-system \
    --create-namespace --wait --debug
./scripts/generate-tls-cert.sh
make docker-buildx
make kind-load-image
kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=gatekeeper-system:default
helm install external-data-provider charts/external-data-provider \
    --set provider.tls.caBundle="$(cat certs/ca.crt | base64 | tr -d '\n\r')" \
    --namespace "${NAMESPACE:-gatekeeper-system}"
kubectl apply -f validation/ingress-ct.yaml
kubectl apply -f validation/ingress-c.yaml
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
