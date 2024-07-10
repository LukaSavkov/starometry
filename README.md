# Starometry
Starometry is a service designed for the collection and management of metrics for the c12s platform. This service gathers metrics from cAdvisor and node-exporter, enabling real-time monitoring of both machine states and virtualized Docker containers running on the machine.

## Configuration
Starometry supports two types of configurations that can be provided via environment variables. The first type pertains to the application itself, while the second type is related to the metrics.

### Application configuration

| Parameter | Description | Default value |
|--|--|--|
| APP_PORT | The port number at which the Starometry HTTP server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from APP_PORT + 1. | 8003 |
| GRPC_PORT | The port number at which the Starometry GRPC server is listening. If there are multiple Starometry agents running, each of them will be assigned a port number starting from GRPC_PORT + 1. | 50055 |
| NODE_EXPORTER_URL | The address without the port that allows Starometry to communicate with the node-exporter running on the same node. | node_exporter |
| NODE_EXPORTER_PORT | The port at which the node exporter is running. | 9100 |
| CADVISOR_URL | The address without the port that allows Starometry to communicate with the cAdvisor running on the same node. | cadvisor |
| CADVISOR_PORT | The port at which the cAdvisor is running. | 8081 |
| NATS_URL | The address without the port that allows Starometry to communicate with the NATS running on the Control Plane | nats |
| NATS_PORT | The port at which the Nats is running. | 4222 |

### Metrics configuration
| Parameter | Description | Default value |
|--|--|--|
| APP_METRICS_CONFIG | List of metrics that you want to scrape from cAdvisor or node-exporter, separated as CSV. | [Link to default list of metrics](#default-metrics) |
| APP_METRICS_CRON_TIMER | Value for the cron job timer that defines how often the scrape for metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end. | 45s |
| APP_METRICS_EXTERNAL_CRON_TIMER | Value for the cron job timer that defines how often the scrape for external metrics will be executed. It is important to note that you must add 's' for seconds or 'm' for minutes at the end. | 45s |

Small example of `APP_METRICS_CONFIG` would be: `container_cpu_usage_seconds_total,container_spec_cpu_quota`
## Usage

The Starometry for HTTP requests is, by default, available at [http://localhost:8003](http://localhost:8003). It can be accessed via any tool that allows you to send HTTP requests. For each instance, just add +1 to the port number.

The Starometry for gRPC requests is, by default, available at [127.0.0.1:50055](127.0.0.1:50055). For each instance of Starometry, just add +1 to the port number. Refer to the [start.sh](https://github.com/c12s/tools/blob/master/start.sh) for more information.

## Endpoints
There are two types of endpoints: gRPC and HTTP.

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
| `address`    | string  | String value of the URL. |


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

## Default Metrics(#default-metrics)

Metrics listed below are must have and always included.

1. container_cpu_usage_seconds_total:
    - **Description**: Total cumulative CPU usage of the container.
    - **Details**: Measures the total CPU time consumed by the container in seconds. This includes both user and system CPU time.
2. container_spec_cpu_quota:
    - **Description**: CPU quota limit set for the container.
    - **Details**: Indicates the maximum amount of CPU time that the container can use during a given period. This is specified in microseconds.
3. container_memory_usage_bytes:
    - **Description**: Current memory usage of the container.
    - **Details**: Shows the total memory usage in bytes, including all memory required by the container's processes, cache, and buffers.
4. container_spec_memory_limit_bytes:
    - **Description**: Memory limit set for the container.
    - **Details**: Specifies the maximum amount of memory the container is allowed to use, in bytes.
5. container_fs_usage_bytes:
    - **Description**: File system usage by the container.
    - **Details**: Represents the total disk space used by the container's filesystem, in bytes.
6. container_spec_cpu_period:
    - **Description**: CPU period for container scheduling.
    - **Details**: Defines the length of the time period in microseconds for CPU allocation, used in conjunction with container_spec_cpu_quota to control CPU resource allocation.
7. container_network_receive_bytes_total:
    - **Description**: Total bytes received by the container.
    - **Details**: Measures the total number of bytes received over the network interfaces of the container.
8. container_network_transmit_bytes_total:
    - **Description**: Total bytes transmitted by the container.
    - **Details**: Measures the total number of bytes sent over the network interfaces of the container.
9. node_cpu_seconds_total:
    - **Description**: Total CPU usage of the node.
    - **Details**: Represents the cumulative CPU time used by all processes on the node, in seconds.
10. node_memory_MemTotal_bytes:
    - **Description**: Total memory available on the node.
    - **Details**: Shows the total amount of physical memory (RAM) available on the node, in bytes.
11. node_memory_MemAvailable_bytes:
    - **Description**: Available memory on the node.
    - **Details**: Indicates the amount of memory that is available for use by processes on the node, in bytes. This includes free memory and reclaimable memory from caches and buffers.
12. node_filesystem_size_bytes:
    - **Description**: Total size of the node's filesystem.
    - **Details**: Represents the total capacity of the node's filesystem, in bytes.
13. node_filesystem_free_bytes:
    - **Description**: Free space in the node's filesystem.
    - **Details**: Indicates the amount of unused space in the node's filesystem, in bytes.
14. node_network_receive_bytes_total:
    - **Description**: Total bytes received by the node.
    - **Details**: Measures the total number of bytes received over all network interfaces on the node.
15. node_network_transmit_bytes_total:
    - **Description**: Total bytes transmitted by the node.
    - **Details**: Measures the total number of bytes sent over all network interfaces on the node.

## Custom metrics that get created in Starometry
Some metrics are calculated and categorized as custom. There are two types of categories: those calculated for containers and those calculated for nodes.

##### Node metrics

1. custom_node_cpu_usage_percentage
2. custom_node_ram_available_mb
3. custom_node_ram_total_mb
4. custom_node_disk_usage_gb
5. custom_node_disk_total_gb
6. custom_node_network_receieve_mb
7. custom_node_network_transmit_mb

##### Service (containers) metrics
1. custom_service_cpu_usage
2. custom_service_ram_usage_mb
3. custom_service_disk_usage_mb
4. custom_service_network_receive_mb
5. custom_service_network_transmit_mb


## Nats communication

Starometry communicates with the Healthcheck service from Protostar via NATS, where Starometry sends the latest scraped metrics to the Healthcheck.