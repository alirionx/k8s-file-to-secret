### Run on Docker -> Example
```
docker run \
  --rm \
  -v ~/.kube/config:/app/kubeconfig \
  -e KUBECONFIG=/app/kubeconfig \
  ghcr.io/alirionx/k8s-file-to-secret:latest
```

### Build Go Binary
```
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../dist/k8s-file-to-secret .
```