# Mastro Helm Chart

## Development Info: testing the Helm-Chart

* `kubectl create namespace mastro`
* `export HELM_NAMESPACE=mastro`
* `helm install mongodb bitnami/mongodb --set auth.username=mastro,auth.password=mastro,auth.database=mastro`
* `helm install mastro mastro/`
* `kubectl get svc -n mastro`  
  ```
  NAME                  TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)     AGE  
  mastro-catalogue      ClusterIP   10.111.101.128   <none>        8085/TCP    33s  
  mastro-featurestore   ClusterIP   10.109.155.11    <none>        8085/TCP    33s  
  mastro-metricstore    ClusterIP   10.109.212.83    <none>        8085/TCP    33s  
  mastro-ui             ClusterIP   10.104.162.18    <none>        80/TCP      33s  
  mongodb               ClusterIP   10.103.255.227   <none>        27017/TCP   45s  
  ```
* `kubectl port-forward -n mastro svc/mastro-ui 8088:80`  
  ```
  Forwarding from 127.0.0.1:8088 -> 80
  Forwarding from [::1]:8088 -> 80
  Handling connection for 8088
  Handling connection for 8088
  Handling connection for 8088
  ```

## Packaging the Helm Chart

The script `pack-it.sh` will generate a `.tgz` file and update the `index.yaml` file accordingly.

## Using the repo and the Helm chart

Add the helm repo:
```
helm repo add mastro https://data-mill-cloud.github.io/mastro/helm-charts
```

Start using the Chart:
```
❯ helm search repo mastro
NAME            CHART VERSION   APP VERSION     DESCRIPTION            
mastro/mastro   0.1.0           0.3.1           A Helm chart for Mastro
```