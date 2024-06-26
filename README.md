# GitHub Sponsors Discord Webhook

Simple webhook server that allows you to receive Discord Webhook notifications
for your GitHub sponsors.

![](https://github.com/Vendicated/SponsorsWebhook/assets/45497981/623e2ba6-8f2c-4ead-b8ba-03ce75e16f4e)

## Setup

1. Clone the repo
    ```
    git clone https://github.com/Vendicated/SponsorsWebhook
    cd SponsorsWebhook
    ```
2. Copy `.env.example` to `.env` and fill out its values
3. Copy `sponsors-webhook.service` to your systemd service directory and start the service
4. Set up a domain for the server, for example via Caddy using the sample `Caddyfile`
5. Visit your domain in the browser and you should see a basic website
6. Go to https://github.com/sponsors/YourNameHere/dashboard/webhooks and press `Add webhook`
7. Inside payload url, enter your domain followed by `/webhook`, like `https://sponsors-webhook.example.com/webhook`
8. Set Content type to application/json
9. Paste your secret (from your .env) file in the Secret box
10. Hit save and you're all done!
