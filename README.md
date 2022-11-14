# Go-chat with Line bot

This is the repository for candidate tech stack project.

- Built in Go version 1.17
- Uses the [Gin server](https://github.com/go-chi/chi)
- Uses [MongoDB](https://github.com/mongodb/mongo-go-driver)
- Uses [Redis](https://github.com/go-redis/redis)
- Uses [Line Bot Go SDK](https://github.com/line/line-bot-sdk-go) 
- Uses [Viper](https://github.com/spf13/viper)

1. Docker and Docker-compose setup required
2. Setup config_template.yaml file and Run ./run.sh
3. Setup [ngrok](https://ngrok.com) for port forwarding
4. Setup Line developer webhook url using ngrok url 
5. using repository pattern for handlers, utils and database
6. using redis for message caching queue