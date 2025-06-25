package lksdk

type Track struct{}

type LocalParticipant struct{}

func (p *LocalParticipant) PublishTrack(track Track, options interface{}) (interface{}, error) {
	return nil, nil
}

type Room struct {
	LocalParticipant *LocalParticipant
}

func ConnectToRoomWithToken(url, token string, config interface{}) (*Room, error) {
	return &Room{LocalParticipant: &LocalParticipant{}}, nil
}

func (r *Room) Disconnect() {}

func NewLocalFileTrack(path string) (Track, error) { return Track{}, nil }
