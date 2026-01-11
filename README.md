# eewbot-go

**msky.haibala.com 用になっています**
**Discord通知機能は無効化されています**

### docker compose
```
services:
  app:
    image: docker.io/nexryai/eew-bot:latest
    restart: always
    environment:
      - MISSKEY_TOKEN=TOKEN
      - DISCORD_WEBHOOK=WEBHOOKURL
      # For Misskey v12 Servers（hotfix）
      - USE_CURL=1
      # 5弱以上じゃないとDiscordに通知しない設定
      - DISCORD_SILENT=1
```
