package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}
