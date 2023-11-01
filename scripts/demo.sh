#!/usr/bin/env bash

. ./scripts/demo-magic.sh

TYPE_SPEED=25
DEMO_PROMPT="${GREEN}➜ ${CYAN}\W ${COLOR_RESET}"
clear

pe "cat ./validation/ingress-c.yaml"
echo " "
pe "cat ./validation/ingress-ct.yaml"
echo " "
pe "kubectl get ingress"
echo " "
pe "cat ./scripts/incorrect-ingress.yaml"
echo " "
echo " "
pe "kubectl apply -f ./scripts/incorrect-ingress.yaml"
echo " "
pe "cat ./scripts/correct-ingress.yaml"
echo " "
echo " "
pe "kubectl apply -f ./scripts/correct-ingress.yaml"
