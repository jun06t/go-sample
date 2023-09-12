package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/slog"
)

type Config struct {
	Addr     string `split_words:"true"`
	Port     int    `split_words:"true"`
	Password string `split_words:"true"`
}

func (c Config) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("addr", c.Addr)
	enc.AddInt("port", c.Port)
	enc.AddString("password", "****") // パスワードをマスク
	//enc.AddString("password", hash(c.Password))
	return nil
}

func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("addr", c.Addr),
		slog.Int("port", c.Port),
		slog.String("password", "****"), // パスワードをマスク
	)
}

var (
	conf    Config
	logger  *zap.Logger
	slogger *slog.Logger
)

func init() {
	envconfig.Process("myapp", &conf)
	logger, _ = zap.NewDevelopment()
	slogger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func main() {
	defer logger.Sync()
	logger.Info("zap", zap.Object("config", conf))
	slogger.Info("slog", "config", conf)
}

const (
	hashPrefixLength = 8
)

// hash returns sha256 hash value of input string.
// You can confirm hash value on terminal as follows:
// $ echo -n "default" | shasum -a 256
func hash(in string) string {
	hash := sha256.Sum256([]byte(in))
	return hex.EncodeToString(hash[:])[:hashPrefixLength] + "****"
}
