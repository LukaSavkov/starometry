# Starometry

Starometry is a service designed for the collection and management of metrics for the c12s platform. This service gathers metrics from cAdvisor and node-exporter, enabling real-time monitoring of both machine states and virtualized Docker containers running on the machine.

### Configuration

Starometry supports two types of configurations that can be provided via environment variables. The first type pertains to the application itself, while the second type is related to the metrics.

### Application configuration

| Parameter          | Description                                                                                                                                                                                 | Default value |
| ------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- |
| APP_PORT           | The port number at which the Starometry HTTP server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from APP_PORT + 1.  | 8003          |
| GRPC_PORT          | The port number at which the Starometry GRPC server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from GRPC_PORT + 1. | 50055         |
| NODE_EXPORTER_URL  | The address without the port that allows Starometry to communicate with the node-exporter running on the same node.                                                                         |
| NODE_EXPORTER_PORT | The port at which the node exporter is running.                                                                                                                                             | 9100          |
| CADVISOR_URL       | The address without the port that allows Starometry to communicate with the cAdvisor running on the same node.                                                                              | cadvisor      |
| CADVISOR_PORT      | The port at which the cAdvisor is running.                                                                                                                                                  | 8081          |
| NATS_URL           | The address without the port that allows Starometry to communicate with the NATS running on the Control Plane                                                                               | nats          |
| NATS_PORT          | The port at which the Nats is running.                                                                                                                                                      | 4222          |

### Metrics configuration

| Parameter                       | Description                                                                                                                                                                                    | Default value                   |
| ------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| APP_METRICS_CONFIG              | List of metrics that you want to scrape from cAdvisor or node-exporter, separated as CSV.                                                                                                      | Link to default list of metrics |
| APP_METRICS_CRON_TIMER          | Value for the cron job timer that defines how often the scrape for metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end.          | 45s                             |
| APP_METRICS_EXTERNAL_CRON_TIMER | Value for the cron job timer that defines how often the scrape for external metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end. | 45s                             |
