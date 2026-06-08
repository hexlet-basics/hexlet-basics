# codebasics Helm chart

Deploys codebasics workloads and gateway resources.

## Validate

```bash
helm lint k8s/app-chart
```

## Render manifests

```bash
helm template codebasics k8s/app-chart -n codebasics
```

## Upgrade/install

```bash
helm upgrade --install codebasics k8s/app-chart -n codebasics --create-namespace
```
