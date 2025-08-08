## Quickstart

Prereqs: kubectl, Docker, Kubebuilder toolchain, Go 1.22+ available locally and kubeconfig pointing to a cluster.

### 1) Install CRDs
```
make install
```

### 2) Run controller locally
```
make run
```
If port :8081 is busy, kill the existing process bound to it and retry.

### 3) Apply sample CR
```
kubectl apply -f config/samples/config_v1_configspec.yaml
```

### 4) Verify
```
kubectl get configspecs.ajay.dev configspec-sample -o jsonpath='{.spec.status}{"\n"}{.spec.time}{"\n"}'
kubectl get cm configspec-sample-config -o yaml
```

Update the CR:
```
kubectl patch configspecs.ajay.dev configspec-sample --type merge -p '{"spec":{"value":"key: new\n"}}'
kubectl get cm configspec-sample-config -o yaml
```

Delete the CR (finalizer removes the ConfigMap):
```
kubectl delete configspecs.ajay.dev configspec-sample
```

### 5) Uninstall CRDs
```
make uninstall
```
