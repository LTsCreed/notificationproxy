# Notification Proxy

## About

Notification Proxy is a simple tool that forwards (proxies) incoming messages to your selected upstream notification providers.

You don't need to open any inbound ports in your network.

![Alt text](media/notificationproxy.png?raw=true "Title")

## Installation (Docker)

1. Download required files

    ```shell
    mkdir notificationproxy && cd notificationproxy
    curl -O https://raw.githubusercontent.com/LTsCreed/notificationproxy/refs/heads/main/docker-compose.yml
    curl -L https://raw.githubusercontent.com/LTsCreed/notificationproxy/refs/heads/main/.env.example -o ".env"
    ```

2. Configure environment variables

    ```shell
    nano .env
    ```

3. Enable at least one notification service

    ```text
    # Discord notifications via Webhook
    # More info: https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks
    # Example:
    # NTFY_DISCORD_URL="https://discord.com/api/webhooks/..."

    ....
    ```

4. Run the service

    ```shell
    docker compose up
    ```

## Inbound services

Currently, two inbound services are supported

### Webhook

The Webhook service listens on port `8080` by default.

You can change this port using the environment variable `SERVER_HOOK_PORT`.

To send a message, configure your software to make a POST request to:

```text
http://notificationproxy.example.com/hook?host=gitea.example.com
```

You can also include additional query parameters — these will be included in the forwarded message.

### SMTP

The SMTP service listens on port `2525` by default.
You can change this using the environment variable `SERVER_SMTP_PORT`.

This service accepts unauthenticated requests.

## Outboud Services

### Discord

Send notifications to a Discord channel using a webhook URL.

### Email

Send notifications to an email inbox via SMTP.
You can use any standard email provider (e.g., Gmail, Yahoo).
