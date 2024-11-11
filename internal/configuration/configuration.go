package configuration

type Configuration interface {
	GetSecretValue() string
	GetMinImageSizeInMB() int
	GetMaxImageSizeInMB() int
}
