package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/open-policy-agent/frameworks/constraint/pkg/externaldata"
	"github.com/open-policy-agent/gatekeeper-external-data-provider/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func Handler(w http.ResponseWriter, req *http.Request, clientset *kubernetes.Clientset) {
	// only accept POST requests
	if req.Method != http.MethodPost {
		utils.SendResponse(nil, "only POST is allowed", w)
		return
	}

	// read request body
	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		utils.SendResponse(nil, fmt.Sprintf("unable to read request body: %v", err), w)
		return
	}
	klog.InfoS("received request", "body", requestBody)

	ingressHosts := getExistingIngressHosts(clientset)

	// parse request body
	var providerRequest externaldata.ProviderRequest
	err = json.Unmarshal(requestBody, &providerRequest)
	if err != nil {
		utils.SendResponse(nil, fmt.Sprintf("unable to unmarshal request body: %v", err), w)
		return
	}

	results := make([]externaldata.Item, 0)
	// iterate over all keys
	for _, key := range providerRequest.Request.Keys {
		// check if key exists in ingressHosts, error if it does
		for _, host := range ingressHosts {
			if key == host {
				results = append(results, externaldata.Item{
					Key:   key,
					Error: "Duplicate Ingress host found " + key + "_invalid",
				})
			}
		}
	}

	utils.SendResponse(&results, "", w)
}

func getExistingIngressHosts(clientset *kubernetes.Clientset) []string {
	// list all the ingress hosts
	ingresses, err := clientset.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	ingressHosts := make([]string, 0)
	for _, ingress := range ingresses.Items {
		for _, rule := range ingress.Spec.Rules {
			ingressHosts = append(ingressHosts, rule.Host)
		}
	}

	return ingressHosts
}
