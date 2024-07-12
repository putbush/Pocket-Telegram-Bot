package storage

type Bucket string

const (
	AccessTokens  = "access_tokens"
	RequestTokens = "request_tokens"
)

type TokenRepository interface {
	Save(chatID int64, requestKey string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
}
