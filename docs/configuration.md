# Configuration

### Configuration options

- **Binary mode**

|Option          |Mandatory|Description|
|----------------|---------|-----------|
`--config-file`      |yes|option to specify the location of the configuration file
`--port`|no|option to specify the port (default : `5957`)
`--enable-https`     |no|Option to enable or disable TLS communication (default : `false`)
`--cert-file`|no|option to specify the location of the certificate file (required if the `--enable-https` option is set to `true`)
`--key-file`|no|option to specify the location of the private key file (required if the `--enable-https` option is set to `true`)

- **Container mode**

|Environment variable|Mandatory|Description|
|--------------------|---------|-----------|
`EASY_API_PROM_SMS_ALERT_CONFIG_FILE`|no|a variable that specifies the location of the configuration file within the container (default: `/etc/easy-api-prom-sms-alert/config.yaml`). It can be overwritten by an external file if the latter is mounted on a volume with the same name and in the same location.
`EASY_API_PROM_SMS_ALERT_PORT`|no|variable to specify the port (default: `5957`)
`EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS`|no|variable to enable or disable TLS communication (default: `true`)
`EASY_API_PROM_SMS_ALERT_CERT_FILE`|no|variable to specify the location of the certificate file (required if the variable `EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS` is set to `true`, default: `/etc/easy-api-prom-sms-alert/tls/server.crt`)
`EASY_API_PROM_SMS_ALERT_KEY_FILE`|no|variable to specify the location of the private key file (required if the variable `EASY_API_PROM_SMS_ALERT_ENABLE_HTTPS` is set to `true`, default: `/etc/easy-api-prom-sms-alert/tls/server.key`)

### Configuration file

```
# Documentation about the configuration file
easy_api_prom_sms_alert:
  # Webhook simulation mode: true -> SMS messages are written to logs
  # and false (production value) -> SMS messages are sent via the provider API
  simulation: true
  
  # Webhook authentication parameters
  auth:
    # Enable authentication: true -> username and password parameters are required
    # To authenticate using Basic header, you must generate the base64 of the string username:password
    enabled: true
    # Username
    username: test
    # Password
    password: test@test

  # SMS provider parameters
  provider:
    # Provider API URL
    url: "http://localhost:5797"
    # Content-Type header accepted by the provider
    # Possible values: "application/json", "application/x-www-form-urlencoded"
    content_type: "application/json"
    # Certificate verification parameter
    # - true -> the HTTPS provider API certificate will not be verified
    # - false -> the HTTPS provider API certificate will be verified (default value)
    insecure_skip_verify: false
    # Provider API authentication parameters
    authentication:
      # Enable provider API authentication:
      # - true -> the provider API requires authentication, and in this case
      #   the authorization_type and authorization_credential parameters are mandatory
      # - false -> the provider API does not require authentication
      enabled: false
      # HTTP Authorization header parameter
      # Example header types: Bearer, Basic, ApiKey
      authorization_type: ''
      # String representing the secret key
      authorization_credential: ''
    # HTTP request body field parameters for the provider API
    parameters:
      # Sender field
      from:
        # Sender field name
        param_name: "from"
        # Sender field value
        param_value: "+1234567890"
        # Method used to send the sender field: post or query
        param_method: "post"
      # Recipient field
      to:
        # Recipient field name
        param_name: "to"
        # Default recipient field value, which must match one of the configured recipient group names
        # in case the team field is missing in an alert field
        param_value: "administration"
        # Method used to send the recipient field: post or query
        param_method: "query"
      # SMS content field
      message:
        # SMS content field name
        param_name: "content"
      # Additional provider parameters. They may be mandatory or optional depending on the provider integration specifications
      # The values below are examples. You should read the provider documentation for proper configuration
      extra_params:
      - param_name: "pn1"
        param_value: "pv1"
        param_method: "post"
      - param_name: "pn2"
        param_value: "pv2"
        param_method: "query"
    # Timeout parameter for calling the provider API (default: 10s)
    timeout: 10s

  # Parameters for the different recipients who will receive SMS messages
  recipients:
  # Recipient group name
  - name: "administration"
    # Members of the recipient group
    members:
    - "+1234567890"
    - "+0987654321"
```