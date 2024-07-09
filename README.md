# Starometry
Starometry is a service designed for the collection and management of metrics for the c12s platform. This service gathers metrics from cAdvisor and node-exporter, enabling real-time monitoring of both machine states and virtualized Docker containers running on the machine.

## Configuration
Starometry supports two types of configurations that can be provided via environment variables. The first type pertains to the application itself, while the second type is related to the metrics.

### Application configuration

| Parameter | Description | Default value |
|--|--|--|
| APP_PORT | The port number at which the Starometry HTTP server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from APP_PORT + 1. | 8003 |
| GRPC_PORT | The port number at which the Starometry GRPC server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from GRPC_PORT + 1. | 50055 |
| NODE_EXPORTER_URL | The address without the port that allows Starometry to communicate with the node-exporter running on the same node. |
| NODE_EXPORTER_PORT | The port at which the node exporter is running. | 9100 |
| CADVISOR_URL | The address without the port that allows Starometry to communicate with the cAdvisor running on the same node. | cadvisor |
| CADVISOR_PORT | The port at which the cAdvisor is running. | 8081 |
| NATS_URL | The address without the port that allows Starometry to communicate with the NATS running on the Control Plane | nats |
| NATS_PORT | The port at which the Nats is running. | 4222 |

### Metrics configuration
| Parameter | Description | Default value |
|--|--|--|
| APP_METRICS_CONFIG | List of metrics that you want to scrape from cAdvisor or node-exporter, separated as CSV. | Link to default list of metrics |
| APP_METRICS_CRON_TIMER | Value for the cron job timer that defines how often the scrape for metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end. | 45s |
| APP_METRICS_EXTERNAL_CRON_TIMER | Value for the cron job timer that defines how often the scrape for external metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end. | 45s |

Small example of APP_METRICS_CONFIG would be: container_cpu_usage_seconds_total,container_spec_cpu_quota
## Usage

The Starometry for HTTP requests is, by default, available at [http://localhost:8003](http://localhost:8003). It can be accessed via any tool that allows you to send HTTP requests. For each instance, just add +1 to the port number.

The Starometry for gRPC requests is, by default, available at [127.0.0.1:50055](127.0.0.1:50055). For each instance of Starometry, just add +1 to the port number. Refer to the [start.sh](https://github.com/c12s/tools/blob/master/start.sh) for more information.

## Endpoints
There is two types of endpoints: gRPC and HTTP.
### HTTP Endpoints
#### Base HTTP Response
```json
{
    "status": 200,
    "data": {}
}
```

#### Base Error HTTP Response
```json
{
    "status": 400,
    "path": "path",
    "time": "2024-07-09",
    "error": "Error"
}
```
#### GET /latest

The endpoint for reading latest written metrics.

##### Request headers

None

#### Request body

None

#### Response - 200 OK

```json
{
    "nodeId": "e984c7e0-0f83-4870-81e7-0424595c90a5",
    "metrics": [
        {
            "metric_name": "container_network_transmit_bytes_total",
            "labels": {
                "id": "/",
                "interface": "br-0758707fa6ae"
            },
            "value": 22096194,
            "timestamp": 1720546066
        },
    ]
}
```

#### POST /place-new-config

The endpoint for adding new configuration metrics.

#### Request body

```json
{
    "queries": [
        "node_filesystem_avail_bytes"
    ]
}
```
|property| type  |                    description                      |
|-----|-----|----|
| `queries`    | array of strings  | Array of strings that are metric names. |

#### Response - 200 OK

```json
{
    "status": 200,
    "data": {
        "status": "OK"
    }
}
```

### gRPC Endpoints

#### /GetLatestMetrics

The endpoint for getting latest metrics.

#### Request body
None

#### Response - 0 OK

```json
{
    "data": {
        "metrics": [
            {
                "labels": {
                    "id": "/"
                },
                "metric_name": "container_memory_usage_bytes",
                "value": 4999598080,
                "timestamp": "1720547417"
            },
        ],
        "node_id": "038a427d-0c78-495d-8781-07cf2707798d"
    }
}
```

#### /PostNewMetrics

The endpoint for adding new metrics in configuration.

#### Request body

```json
{
    "metrics": [
        "node_filesystem_avail_bytes"
    ]
}
```
|property| type  |                    description                      |
|-----|-----|----|
| `metric`    | array of strings  | Array of strings that are metric names. |

#### Response - 0 OK

```json
{
    "data": {
        "metrics": [
            {
                "labels": {
                    "id": "/"
                },
                "metric_name": "container_spec_cpu_period",
                "value": 0,
                "timestamp": "1720547610"
            },
        ],
        "node_id": "038a427d-0c78-495d-8781-07cf2707798d"
    }
}
```
#### /PostNewExternalApplicationsList

The endpoint for adding new addresses for external applications.

#### Request body

```json
{
    "external_applications": [
        {
            "address": "example-address:8080"
        },
    ],
}
```
|property| type  |                    description                      |
|-----|-----|----|
| `external_applications`    | array of applications-url objects | Array of applications-url. |
| `address`    | string  | String value of URL. |


#### Response - 0 OK

```json
{
    "external_applications": [
        {
            "address": "external-app"
        }
    ]
}
```


