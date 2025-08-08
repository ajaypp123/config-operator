## Development

### Setup
- Go 1.22+
- Docker
- kubectl
- Kubebuilder tooling (controller-gen, kustomize, envtest are auto-installed by `make`)

Install dependencies and generate manifests:
```
make manifests generate
```

### Run unit tests
```
make test
```

### Common make targets
- `make manifests`: regenerate CRDs/RBAC from markers
- `make generate`: DeepCopy and other generated code
- `make fmt && make vet`: format and vet
- `make test`: run unit tests (with envtest)
- `make install` / `make uninstall`: install/uninstall CRDs
- `make run`: run manager locally
- `make build`: build binary to `bin/manager`

### Build and deploy an image
Set your registry (example uses local registry at `localhost:32000`):
```
make docker-build IMG=localhost:32000/config-operator:dev
make docker-push IMG=localhost:32000/config-operator:dev
```

Deploy controller to cluster:
```
make deploy IMG=localhost:32000/config-operator:dev
```

Undeploy:
```
make undeploy
```

### Sample CR
`config/samples/config_v1_configspec.yaml`:
```
apiVersion: ajay.dev/v1
kind: ConfigSpec
metadata:
  name: configspec-sample
spec:
  value: |
    logLevel: info
    featureFlags:
      alpha: true
  time: "2025-01-01T00:00:00Z"
  status: Pending
```

### Notes
- Controller updates `spec.time` and `spec.status` automatically.
- ConfigMap name: `<cr-name>-config`, key: `config.yaml`.
- On deletion, finalizer removes the associated ConfigMap.
