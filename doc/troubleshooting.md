## Troubleshooting

### Operator fails to start: address already in use :8081
- Another process is listening on the health probe port.
- Find and kill it:
```
ss -ltnp '( sport = :8081 )'
kill -9 <pid>
```

### ConfigMap not created/updated
- Check controller logs from `make run`.
- Validate the YAML in `spec.value` is valid. Invalid YAML sets `spec.status=Error` and logs an error.

### RBAC errors creating ConfigMaps
- Ensure the manager `ClusterRole` includes core `configmaps` permissions (this repo does).
- Re-run:
```
make manifests && make deploy IMG=<your-image>
```

### CRD/API mismatch
- API group is `ajay.dev/v1`. Ensure your CRs use `apiVersion: ajay.dev/v1`.

### Uninstall / reinstall
```
make undeploy || true
make uninstall || true
make install
```
