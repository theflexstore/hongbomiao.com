#!/usr/bin/env bash

set -e


# Patch Ingress
for cluster in west east; do
  echo "# Patch Ingress on: k3d-${cluster}"
  kubectl patch configmap ingress-nginx-controller \
    --context="k3d-${cluster}" \
    --namespace=ingress-nginx \
    --patch "$(cat kubernetes/patches/ingress-nginx-controller-configmap-patch.yaml)"
  kubectl patch deployment ingress-nginx-controller \
    --context="k3d-${cluster}" \
    --namespace=ingress-nginx \
    --patch "$(cat kubernetes/patches/ingress-nginx-controller-deployment-patch.yaml)"
  echo "=================================================="
done
