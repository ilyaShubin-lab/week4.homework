package postgres

type Config struct {
	ImageName string
	Database  string
	User      string
	Password  string
}

type Option func(*Config)

func WithImage(image string) Option {
	return func(c *Config) { c.ImageName = image }
}

func WithDatabase(db string) Option {
	return func(c *Config) { c.Database = db }
}

func WithAuth(user, password string) Option {
	return func(c *Config) {
		c.User = user
		c.Password = password
	}
}
