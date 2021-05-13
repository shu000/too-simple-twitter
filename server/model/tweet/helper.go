package tweet

import "github.com/google/uuid"

// UUID version 4 を生成
func GetUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
