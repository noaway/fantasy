package apns

type APNsapi interface {
	Aps() []byte
	Token() string
	Topic() string
}
