# java-benchmarks

Benchmarks for backend Java frameworks

```bash
cd spring-kernel-threads
kubectl apply -f deployment/postgres-deployment.yaml
kubectl apply -f deployment/app-deployment.yaml
kubectl port-forward service/backend-api-service 8080:8080
```

```bash
> curl -X POST -d "First" localhost:8080/api/message
1
> curl -X GET -d "First" localhost:8080/api/message/1
First
```
