# Installation

These actions are performed under a Linux server.

### Prerequisites

```
mkdir -p $HOME/alert-prometheus && cd $HOME/alert-prometheus
```

```
vi config.yaml
```

```
easy_api_prom_sms_alert:
  simulation: true
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "http://localhost:5797"
    content_type: "application/json"
    authentication:
      enabled: false
      authorization_type: ''
      authorization_credential: ''
    parameters: 
      from: 
        param_name: "from"
        param_value: "+1234567890"
        param_method: "post"
      to:
        param_name: "to"
        param_value: "administration"
        param_method: "query"
      message: 
        param_name: "content"
    timeout: 10s
  recipients: 
  - name: "administration"
    members:
    - "+1234567890"
    - "+0987654321"
  - name: "urgence"
    members:
    - "+1122334455"
    - "+5544332211"
```


### Installation via a package under Linux

```
cd $HOME && mkdir -p alert-prometheus && cd alert-prometheus
```

```
curl -LO https://github.com/willbrid/easy_api_prom_sms_alert/releases/download/v<VERSION>/easy_api_prom_sms_alert_<VERSION>_linux_amd64.tar.gz
```

```
tar -xvzf easy_api_prom_sms_alert_<VERSION>_linux_amd64.tar.gz
```

```
./easy_api_prom_sms_alert_<VERSION>_linux_amd64 --config-file ./config.yaml
```

Replace **\<VERSION\>** with the desired version number (greater than or equal to **1.3.4**).

### Installation via a Docker container

--- **Installation using the default configuration file and enabling the https protocol**

```
docker run -d -p 8000:5957 \
       --name alert-sms-sender -\
       -e EASY_API_PROM_SMS_ALERT_PORT=5957 \
       -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true \
       willbrid/easy-api-prom-sms-alert:latest
```

In this example, the default port **5957** internal to the container is mapped to the external port **8000**.

--- **Installation using a persistent volume for the config.yaml file and enabling the https protocol**

The idea here is to allow customization of the **config.yaml** configuration file. Enabling the **https** protocol uses **server.crt** and **server.key** files by default, located in the **/etc/easy-api-prom-sms-alert/tls/server.key** directory within the container.

```
docker run -d -p 8000:5957 \
       --name alert-sms-sender \
       -v $HOME/alert-prometheus/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml \
       -e EASY_API_PROM_SMS_ALERT_PORT=5957 \
       -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true \
       willbrid/easy-api-prom-sms-alert:latest
```

--- **Installation using a persistent volume for the config.yaml file and enabling the https protocol with the tls files**

We will assume the existence of our tls files: **server.crt** for the certificate and **server.key** the private key, in the directory **$HOME/alert-prometheus**.

```
docker run -d -p 8000:5957 \
       --name alert-sms-sender \
       -v $HOME/alert-prometheus/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml \
       -v $HOME/alert-prometheus/server.crt:/etc/easy-api-prom-sms-alert/tls/server.crt \
       -v $HOME/alert-prometheus/server.key:/etc/easy-api-prom-sms-alert/tls/server.key \
       -e EASY_API_PROM_SMS_ALERT_PORT=5957 \
       -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true \
       -e EASY_API_PROM_SMS_ALERT_CERT_FILE=/etc/easy-api-prom-sms-alert/tls/server.crt \
       -e EASY_API_PROM_SMS_ALERT_KEY_FILE=/etc/easy-api-prom-sms-alert/tls/server.key \
       willbrid/easy-api-prom-sms-alert:latest
```