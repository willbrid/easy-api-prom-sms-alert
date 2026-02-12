# Easy-api-prom-sms-alert

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/willbrid/easy_api_prom_sms_alert/blob/main/LICENSE) [![CI](https://github.com/willbrid/easy-api-prom-sms-alert/actions/workflows/ci.yml/badge.svg)](https://github.com/willbrid/easy-api-prom-sms-alert/actions/workflows/ci.yml)

**Easy-api-prom-sms-alert** is a *prometheus* webhook that allows sending SMS alerts using any SMS provider exposing an API.

### Issue

When *prometheus* detects abnormal conditions in the monitored systems, it triggers alerts to notify Alertmanager to send SMS notifications. However, there are many types of SMS providers, each with its own specifications. Integrating several of them directly into Alertmanager would make the configuration complex and difficult to manage.

### Solution

With **Easy-api-prom-sms-alert**, users can choose any SMS service provider that exposes an HTTP POST API. This gives them the freedom to select the provider that best meets their needs in terms of cost and reliability.

This tool supports multiple providers, such as **Twilio**, **WhatsApp** and **Slack**.

### Documentation

1- [Installation](https://github.com/willbrid/easy-api-prom-sms-alert/blob/main/docs/installation.md) <br>
2- [Configuration](https://github.com/willbrid/easy-api-prom-sms-alert/blob/main/docs/configuration.md) <br>
3- [Usage](https://github.com/willbrid/easy-api-prom-sms-alert/blob/main/docs/usage.md) <br>
4- [Complete example](https://github.com/willbrid/easy-api-prom-sms-alert/blob/main/docs/complete-example.md)

### License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/willbrid/easy-api-prom-sms-alert/blob/main/LICENSE) file for more details.