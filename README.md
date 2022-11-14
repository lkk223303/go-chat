# Go-chat with Line bot

This is the repository for candidate tech stack project.
![go-chat](https://user-images.githubusercontent.com/10274839/201698449-d7993bfb-49d1-4818-b7de-b3e550fde587.png)

- Built in Go version 1.17
- Uses the [gin server](https://github.com/go-chi/chi)
- Uses [MongoDB](https://github.com/mongodb/mongo-go-driver)
- Uses [Line Bot Go SDK](https://github.com/line/line-bot-sdk-go) 
- Uses [Viper](https://github.com/spf13/viper)

1. Requires Docker and Docker-compose setup
2. Setup config_template.yaml file and Run ./run.sh
3. Setup [ngrok](https://ngrok.com) for port forwarding
4. Setup Line developer webhook url as ngrok url 
5. using repository pattern for handlers, utils and database
