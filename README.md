# CloudEvents Converting Proxy(CECP)
This proxy converts a CloudEvents contents "application/cloudevents+json" to a simple http json content "application/json".

# Environment Values

| Key | Default | Description |
| ---- | ---- |  ---- |  
| HTTP_HOST | 0.0.0.0 | Proxy host IP which retrieves CloudEvents. |
| HTTP_HOST | 8080 | Proxy port which retrives CloudEvents.|
| TARGET_URL | http://localhost:8888 | Destination server which retrives converted Json message. (i.e. Fluetnt-bit URL) |
| PROMETHEUS_ADDRESS | 0.0.0.0 | Prometheus exporter host IP. |
| PROMETHEUS_PORT | 9100 | Prometheus exporter host port. |
| PROMETHEUS_URL_PATH | /metric | Prometheus exporter URL path. |