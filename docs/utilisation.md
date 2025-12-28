# Usage

Our previous configuration file shows that **Basic** authentication is enabled. Therefore, we need to generate the **base64** version of the **username:password** string in order to use it in the **Authorization** header.

```
echo -n test:test@test | base64
```

```
dGVzdDp0ZXN0QHRlc3Q=
```

#### **Test to be performed with curl**

```
curl -k --location 'https://localhost:8000/api-alert' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic dGVzdDp0ZXN0QHRlc3Q=' \
--data '{
  "version": "4",
  "groupKey": "{alertname=\"InstanceDown\"}",
  "status": "firing",
  "receiver": "webhook",
  "groupLabels": {
    "alertname": "InstanceDown"
  },
  "commonLabels": {
    "alertname": "InstanceDown",
    "job": "myjob",
    "severity": "critical"
  },
  "commonAnnotations": {
    "summary": "Instance xxx down",
    "description": "The instance xxx is down."
  },
  "externalURL": "http://prometheus.example.com",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "InstanceDown",
        "instance": "localhost:9090",
        "team": "urgence",
        "job": "myjob",
        "severity": "critical"
      },
      "annotations": {
        "summary": "Instance localhost:9090 down",
        "description": "The instance localhost:9090 is down."
      },
      "startsAt": "2023-05-20T14:28:00.000Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus.example.com/graph?g0.expr=up%3D%3D0&g0.tab=1"
    }
  ]
}'
```

#### **Integration with Alertmanager**

To integrate **Easy-api-prom-sms-alert** into **Alertmanager**, you need to configure a webhook by adding a webhook receiver to your **Alertmanager** configuration.

```
receivers:
- name: 'admin'
  webhook_configs:
  - url: 'https://localhost:8000/api-alert'
    send_resolved: false
    http_config: 
      basic_auth:
        username: test
        password: test@test
      tls_config:
        insecure_skip_verify: true
```

To view the result in **simulation** mode, you need to consult the container logs.

```
docker container logs alert-sms-sender
```