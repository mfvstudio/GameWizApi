package helpers

import (
	"crypto/rand"
	"log"

	"github.com/google/uuid"
	"github.com/mfvstudio/gamewizapi/cmd/gen"
)

func CreateSessionFromSessionRequest(s *gen.CreateSession) (*gen.Session, error) {
	uuid, err := generateUUID()
	if err != nil {
		log.Printf("Failed to generate UUID for session creation")
		return nil, err
	}
	sesh := gen.Session{
		GameId:         s.GameId,
		GameState:      nil,
		MaxPlayerCount: int(s.MaxPlayerCount),
		MetaData:       s.MetaData,
		MinPlayerCount: int(s.MinPlayerCount),
		Players:        []string{s.HostUID},
		SessionId:      *uuid,
		Status:         gen.WAITINGFORPLAYERS,
		UpNext:         s.UpNext,
	}
	return &sesh, nil
}

func generateUUID() (*string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	str := id.String()
	return &str, nil
}

func GenerateSessionInviteCode() string {
	codeLength := 6
	code := rand.Text()
	return code[0:codeLength]
}
