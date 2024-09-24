[![GitHub](https://img.shields.io/badge/GitHub-Repository-green)](https://github.com/chazari-x/training-api-bot)
[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-chazari%2Ftraining--api-blue)](https://hub.docker.com/r/chazari/training-api)
[![Discord](https://img.shields.io/badge/Discord-Server-blue)](https://czo.ooo/invite)

# Discord API Training Bot

Этот Discord бот предназначен для взаимодействия с API сервера Training в игре SA-MP (San Andreas Multiplayer).

## SA-MP

SA-MP (San Andreas Multiplayer) - это бесплатный мод для игры Grand Theft Auto: San Andreas (TM) от Rockstar Games. SA-MP позволяет игрокам играть вместе в режиме онлайн, создавая собственные сервера с уникальными модификациями и правилами. Бот обеспечивает взаимодействие с сервером Training, который является одним из множества серверов SA-MP.

## API

API сервера предоставляет доступ к информации об администрации и его пользователях. Бот использует этот API для получения данных об администраторах и пользователях сервера Training.

## Команды

- `/admins`: Получение списка администраторов сервера Training. Команда предоставляет их имена, последнюю авторизацию и количество выданных варнов.
- `/user <nickname>`: Получение информации о конкретном пользователе сервера Training. Команда предоставляет никнейм пользователя, его accID, статус аккаунта, активность и список полученных наказаний.

## Установка

1. Убедитесь, что на вашем сервере установлен Docker.
2. Создайте файл `docker-compose.yml` со следующим содержимым:

```yaml
version: '3'
services:
  discord-api-training-bot:
    container_name: discord-api-training-bot
    image: "chazari/training-api:latest"
    restart: always
    command: ["/app/main", "discord", "--token=your-discord-bot-token"]
```

3. Замените `"your-discord-bot-token"` на токен вашего Discord бота.
4. Запустите контейнер с помощью команды:

```
docker-compose up -d
```

## Использование

После успешной установки и запуска бота, он готов к использованию на вашем Discord сервере. Для получения списка администраторов сервера Training используйте команду `/admins`. Для получения информации о конкретном пользователе сервера используйте команду `/user <nickname>`.

## Вклад

Мы приветствуем ваши предложения и вклад в развитие проекта. Если у вас есть идеи или обнаружили ошибку, создайте issue или отправьте pull request.

## Лицензия

Этот проект лицензируется по лицензии GNU General Public License v3.0. См. файл [LICENSE](https://github.com/chazari-x/training-api-bot/blob/master/LICENSE) для получения дополнительной информации.
