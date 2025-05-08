To call the `IsActive` gRPC API with proper parameters using `grpcurl`, you'll need to:

1. First, ensure you have `grpcurl` installed (https://github.com/fullstorydev/grpcurl)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

2. Here's the command to call the `IsActive` endpoint:

```bash
grpcurl \
  -plaintext \  # since we're not using TLS
  -d '{
    "name": "your-scaled-object-name",
    "namespace": "your-namespace",
    "scalerMetadata": {
      "longitude": "-122.4194",
      "latitude": "37.7749"
    }
  }' \
  localhost:6000 \
  externalscaler.ExternalScaler/IsActive
```

### Explanation of parameters:

1. `-plaintext`: Disables TLS (since your server isn't using it)
2. `-d`: The request payload in JSON format containing:
   - `name`: Name of your ScaledObject
   - `namespace`: Kubernetes namespace where the ScaledObject exists
   - `scalerMetadata`: Map containing the required parameters (longitude/latitude in your case)
3. `localhost:6000`: Address of your gRPC server
4. `externalscaler.ExternalScaler/IsActive`: Fully qualified service/method name

### Expected Response:

You should get a response like:
```json
{
  "result": true
}
```

### Additional Tips:

1. If you need to list all available methods first:
   ```bash
   grpcurl -plaintext localhost:6000 list
   ```

2. To describe the service and methods:
   ```bash
   grpcurl -plaintext localhost:6000 describe externalscaler.ExternalScaler
   ```

3. For the other endpoints:
   - `StreamIsActive`: Similar to IsActive but will stream responses
   - `GetMetricSpec`: Takes same parameters but returns metric specs
   - `GetMetrics`: Requires metricName in addition to scaled object ref

4. If you're calling this from outside the cluster, you'll need to:
   - Port-forward to your service: `kubectl port-forward service/<your-service> 6000:6000`
   - Or expose it via Ingress/NodePort

Would you like me to provide examples for the other endpoints as well?
