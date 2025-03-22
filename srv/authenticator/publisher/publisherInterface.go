package publisher

import "crypto"

type IAuthPublisher interface {
	PublishKey(puk crypto.PublicKey, issuer string) error
}
