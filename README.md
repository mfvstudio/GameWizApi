### GameWizApi

##### Generating Server Client from OpenApi spec
- run the command `go tool oapi-codegen -config codegen_config.yaml GameWizOpenApiSpec.yaml`

#### Building and running 
- build package with `go build cmd/app/*.go` and run with `./main`


GameWizApi is a game server for turn-based games. It is general enough to work with many game formats (ex. TicTacToe, Chess, Darts). It can be used in game engine clients to create multiplayer game sessions. this api handle user account creation and user login by utilizing Firebase Authentication. After logging in, users are given an ID, which can be passed to api requests for authNZ. The bearer token has claims to verify if users can perform certain actions. Users can create Game Sessions and then share a 6 character UID to other users so they can join. There is no limit on the number of players that can join a session(Which I think I will implement one if issues arise). I wanted to test out the Google Cloud Platform, so I stuck with their services for this project. Also, I wanted to learn Golang (I heard great things!) so I used it to create the backend server. So far, this has been an exciting project. GCP is a great platform and may use it for other projects soon!
