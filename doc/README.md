## config-operator

### Overview
`config-operator` is a Kubernetes operator built with Kubebuilder (Go 1.22+, controller-runtime v0.21) that turns a custom resource into a `ConfigMap`.

### CRD
- **Group**: `ajay.dev`
- **Version**: `v1`
- **Kind**: `ConfigSpec`
- **Fields**:
  - `spec.value` (string): raw YAML content to store
  - `spec.time` (string, RFC3339): timestamp maintained by the controller
  - `spec.status` (string): one of `Pending`, `Applied`, `Error`

### Controller behavior
1. On create:
   - Validates `spec.value` as YAML
   - Creates a `ConfigMap` named `<cr-name>-config` in the same namespace with key `config.yaml`
   - Sets `spec.time` to now and `spec.status` to `Applied` on success; `Error` on failure
2. On update:
   - Updates the associated `ConfigMap` with the new `spec.value`
   - Updates `spec.time` to now and sets `spec.status` accordingly
3. On delete:
   - Finalizer removes the associated `ConfigMap`

### Key manifests and code
- API types: `api/v1/`
- Controller: `internal/controller/configspec_controller.go`
- CRD: `config/crd/bases/config.ajay.dev_configspecs.yaml`
- RBAC: `config/rbac/role.yaml` (includes permissions for core `configmaps`)
- Sample: `config/samples/config_v1_configspec.yaml`

### Quick links
- Usage and testing (local run): see `doc/quickstart.md`
- Usage with Docker image (build/push/deploy): see `doc/quickstart-docker.md`
- Development and tests: see `doc/development.md`
- Troubleshooting: see `doc/troubleshooting.md`
