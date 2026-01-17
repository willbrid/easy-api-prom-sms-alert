# Complete example

### Setting up a sandbox

We will set up a vagrant **Rocky Linux 9** server on an Ubuntu 24.04 (or >= 24.04) host machine with **virtualbox7** and **vagrant** tools already installed.

```
mkdir $HOME/easy-api-prom-sms-alert && cd $HOME/easy-api-prom-sms-alert
```

```
wget https://download.virtualbox.org/virtualbox/7.0.24/VBoxGuestAdditions_7.0.24.iso
```

```
vi Vagrantfile
```

```
# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vbguest.auto_update = false
  config.vbguest.no_remote = true
  config.vbguest.iso_path = "./VBoxGuestAdditions_7.0.24.iso"

  # General Vagrant VM configuration
  config.ssh.insert_key = false
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.provider :virtualbox do |v|
    v.memory = 4096
    v.cpus = 2
    v.linked_clone = true
  end

  # Monitoring Server
  config.vm.define "monitoring-server" do |srv|
    srv.vm.box = "rockylinux/9"
    srv.vm.box_version = "6.0.0"
    srv.vm.hostname = "monitoring-server"
    srv.vm.network :private_network, ip: "192.168.56.211"
  end
end
```

```
vagrant up
```

### Installation of prometheus, alertmanager, node-exporter and easy-api-prom-sms-alert on our monitoring-server

In this section, we will install the following components in a containerized format on the vagrant **monitoring-server** :

- **prometheus**: for server monitoring
- **alertmanager**: to receive prometheus alerts and trigger notifications
- **node-exporter**: to collect system metrics on the vagrant **monitoring-server**
- **easy-api-prom-sms-alert**: for the webhook that will be configured within **alertmanager** to send SMS alerts

```
vagrant ssh monitoring-server
```

##### Setting up node-exporter

```
podman run -d --net="host" \
       --pid="host" \
       -v "/:/host:ro,rslave" \
       -u root \
       --name node-exporter \
       quay.io/prometheus/node-exporter:v1.7.0 \
       --path.rootfs=/host
```

##### Setting up alertmanager

```
mkdir -p $HOME/monitoring/alertmanager && mkdir $HOME/monitoring/alertmanager/data && cd $HOME/monitoring/alertmanager
```

```
vi alertmanager.yml
```

```
route:
  receiver: 'admin'
  repeat_interval: 1h

receivers:
- name: 'admin'
  webhook_configs:
  - url: 'https://192.168.56.211:5797/api-alert'
    send_resolved: false
    http_config: 
      basic_auth:
        username: test
        password: test@test
      tls_config:
        insecure_skip_verify: true
```

```
podman run -d --net=host \
       -v $HOME/monitoring/alertmanager/alertmanager.yml:/config/alertmanager.yml:z \
       -v $HOME/monitoring/alertmanager/data:/data:z \
       --name alertmanager \
       prom/alertmanager:v0.26.0 \
       --config.file=/config/alertmanager.yml \
       --log.level=debug
```

##### Setting up prometheus

```
mkdir -p $HOME/monitoring/prometheus && mkdir $HOME/monitoring/prometheus/data && mkdir $HOME/monitoring/prometheus/rules && cd $HOME/monitoring/prometheus
```

```
vi prometheus.yml
```

```
global:
  scrape_interval:     10s 
  evaluation_interval: 10s
  external_labels:
    cluster: CLUSTER_A
    replica: 0

rule_files:
- "/etc/prometheus/rules/*"

alerting:
 alertmanagers:
 - static_configs:
   - targets: ['192.168.56.211:9093']

scrape_configs:
  - job_name: 'prometheus'
    scrape_interval: 10s
    static_configs:
    - targets: ['192.168.56.211:9090']

  - job_name: 'node-exporter'
    static_configs:
    - targets: ['192.168.56.211:9100']
```

```
vi rules/nodes-rules.yml
```

```
groups:
- name: NODE_CLUSTER_A
  rules:
  - alert: NodeDown
    expr: up{job="node-exporter"} == 0
    for: 1m
    labels:
      severity: critical
      team: urgence
    annotations:
      summary: "Node is down"
      description: "The node has been down for the last 1 minute."
```

```
podman run -d --net=host \
    -v $HOME/monitoring/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:z \
    -v $HOME/monitoring/prometheus/data:/prometheus:z \
    -v $HOME/monitoring/prometheus/rules:/etc/prometheus/rules:z \
    -u root \
    --name prometheus \
    quay.io/prometheus/prometheus:v2.38.0 \
    --config.file=/etc/prometheus/prometheus.yml \
    --storage.tsdb.retention.time=4h \
    --storage.tsdb.path=/prometheus \
    --storage.tsdb.max-block-duration=2h \
    --storage.tsdb.min-block-duration=2h \
    --web.listen-address=:9090 \
    --web.external-url=http://192.168.56.211:9090 \
    --web.enable-lifecycle \
    --web.enable-admin-api
```

##### Setting up easy-api-prom-sms-alert

```
mkdir -p $HOME/monitoring/alert && cd $HOME/monitoring/alert
```

Next, you would need to use one of the configuration files (which you can customize) below to set up the integration with a provider : **Twilio**, **WhatsApp**, **Slack**,...

```
vi $HOME/monitoring/alert/config.yaml
```

- **Twilio integration**

```
easy_api_prom_sms_alert:
  simulation: false
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "https://api.twilio.com/2010-04-01/Accounts/XXXXXXX/Messages"
    content_type: "application/x-www-form-urlencoded"
    authentication:
      enabled: true
      authorization_type: 'Basic'
      authorization_credential: 'YYYYYYY'
    parameters: 
      from: 
        param_name: "From"
        param_value: "+xxxxxxx"
        param_method: "post"
      to:
        param_name: "To"
        param_value: "urgence"
        param_method: "post"
      message: 
        param_name: "Body"
    timeout: 10s
  recipients: 
  - name: "urgence"
    members:
    - "+yyyyyyy"
    - "+zzzzzzz"
```

**XXXXXXX** is the **SID** string to retrieve from the **Twilio** platform. <br>
**YYYYYYY** is the base64 hash of the **SID:TOKEN** string to retrieve from the **Twilio** platform. <br>
**+xxxxxxx** is the sender's number. <br>
**+yyyyyyy** and **+zzzzzzz** are the phone numbers that will receive SMS alerts.

**Reference** : [Twilio Documentation](https://www.twilio.com/en-us/blog/send-sms-twilio-shell-script-curl)

- **WhatsApp integration**

```
easy_api_prom_sms_alert:
  simulation: false
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "https://api.wassenger.com/v1/messages"
    content_type: "application/json"
    authentication:
      enabled: true
      authorization_type: 'Token'
      authorization_credential: 'API_TOKEN'
    parameters: 
      from: 
        param_name: "reference"
        param_value: "prometheus"
        param_method: "post"
      to:
        param_name: "phone"
        param_value: "urgence"
        param_method: "post"
      message: 
        param_name: "message"
    timeout: 10s
  recipients: 
  - name: "urgence"
    members:
    - "+xxxxxxx"
    - "+yyyyyyy"
```

**API_TOKEN** is the token that can be retrieved on the **Wassenger** platform. <br>
**+xxxxxxx** and **+yyyyyyy** are the WhatsApp accounts that will receive SMS alerts.

**Reference** : [Wassenger Documentation](https://app.wassenger.com/docs/)

- **Slack integration**

```
easy_api_prom_sms_alert:
  simulation: false
  auth:
    enabled: true
    username: test
    password: test@test
  provider:
    url: "https://slack.com/api/chat.postMessage"
    content_type: "application/json"
    authentication:
      enabled: true
      authorization_type: 'Bearer'
      authorization_credential: 'API_TOKEN'
    parameters: 
      from: 
        param_name: "username"
        param_value: "prometheus"
        param_method: "post"
      to:
        param_name: "channel"
        param_value: "administration"
        param_method: "post"
      message: 
        param_name: "text"
    timeout: 10s
  recipients: 
  - name: "administration"
    members:
    - "alerts-channel"
```

**API_TOKEN** is the token that can be retrieved on the **Slack** platform. <br>
**alerts-channel** is provided as an example of an existing Slack channel. You can create and use your own channel instead.

**Reference** : [Slack Documentation](https://docs.slack.dev/apis/)

- **Starting easy-api-prom-sms-alert with one of the integration contents**

```
podman run -d --net=host \
       --name alert-sms-sender \
       -v $HOME/monitoring/alert/config.yaml:/etc/easy-api-prom-sms-alert/config.yaml:z \
       -e EASY_API_PROM_SMS_ALERT_PORT=5957 \
       -e EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS=true \ 
       willbrid/easy-api-prom-sms-alert:latest
```

### Test

In order to simulate the shutdown of the **monitoring-server**, we stop the **node-exporter** container.

```
podman container stop node-exporter
```

After one minute, we will see an SMS alert on our phone (**Twilio integration**) or on our WhatsApp account (**WhatsApp integration**).