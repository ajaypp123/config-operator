## Quickstart (Docker image + deploy)

Prereqs: kubectl, Docker, Go 1.22+, Kubebuilder toolchain. Cluster must be able to pull from your registry (example uses `localhost:32000`).

### 1) Build and push the controller image
Set your registry/tag (edit as needed):
```
export IMG=localhost:32000/config-operator:dev
make docker-build IMG=$IMG
make docker-push IMG=$IMG
```

Note: Ensure the cluster can pull from `localhost:32000` (configure as an insecure registry if needed).

### 2) Deploy the operator to the cluster
This applies CRDs, RBAC, namespace, and the controller Deployment.
```
make deploy IMG=$IMG
```

Check rollout:
```
kubectl -n config-operator-system get deploy,po
```

Optional logs:
```
kubectl -n config-operator-system logs deploy/config-operator-controller-manager -c manager -f
```

### 3) Test with the sample CR
Apply the sample:
```
kubectl apply -f config/samples/config_v1_configspec.yaml
```

Verify status/time and ConfigMap:
```
kubectl get configspecs.ajay.dev configspec-sample -o jsonpath='{.spec.status}{"\n"}{.spec.time}{"\n"}'
kubectl get cm configspec-sample-config -o yaml
```

Update the CR and confirm the ConfigMap updates:
```
kubectl patch configspecs.ajay.dev configspec-sample --type merge -p '{"spec":{"value":"key: new\n"}}'
kubectl get cm configspec-sample-config -o yaml
```

Delete the CR (finalizer removes its ConfigMap):
```
kubectl delete configspecs.ajay.dev configspec-sample
```

### 4) Undeploy and uninstall
Remove the controller Deployment and RBAC:
```
make undeploy
```

Uninstall CRDs:
```
make uninstall
```

### Notes
- API group is `ajay.dev/v1`; sample uses `apiVersion: ajay.dev/v1`.
- Controller writes `spec.time` (RFC3339) and `spec.status` (Applied/Error).
- ConfigMap name: `<cr-name>-config`, key: `config.yaml`.
