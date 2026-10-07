package config

type LoggerConfig interface {
	Level() string //// debug / info / warn / error
	AsJSON() bool  // true — JSON (для Docker и прода), false — читаемый текст
}

// настройки gRPC-сервера payment.
type PaymentGRPCConfig interface {
	Address() string
}
