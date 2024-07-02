A bot for continuously trying to launch an instance on Hetzner. Here are the example environment variables for launching the bot:

```ini
HETZNER_API_KEY=YOUR_HETZNER_API_KEY
HETZNER_INSTANCE_TYPE=YOUR_INSTANCE_NAME (example: cx32)
HETZNER_DATA_CENTER=YOUR_DATA_CENTER (example: nbg1-dc3)
SSH_KEY_ID=ID_OF_YOUR_SSH_KEY (example: 12345678)
ERROR_WAIT_SECONDS=60
PORT=5000
CGO_ENABLED=0
```
