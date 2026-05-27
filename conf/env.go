package conf

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type EnvKey string

func (e EnvKey) String() string {
	return string(e)
}

const (
	BOT_TOKEN            EnvKey = "BOT_TOKEN"
	BOT_OWNER            EnvKey = "BOT_OWNER"
	BOT_ADMINS           EnvKey = "BOT_ADMINS"
	BOT_DEBUG            EnvKey = "BOT_DEBUG"
	BOT_VERSION          EnvKey = "BOT_VERSION"
	BOT_IS_DEVELOPMENT   EnvKey = "BOT_IS_DEVELOPMENT"
	BOT_PLATFORM_API_KEY EnvKey = "BOT_PLATFORM_API_KEY"
)

type Env struct {
	botToken         string
	botOwner         int64
	botAdmins        []int64
	botDebug         bool
	botVersion       string
	botIsDevelopment bool
	// treat this as the personal access API key token for automation or personal account management
	botPlatformAPIKey string
}

func (e *Env) Load() {
	botToken := os.Getenv(BOT_TOKEN.String())
	if botToken == "" {
		log.Printf("WARN! %s not set ", BOT_TOKEN)
	}
	botOwner, err := strconv.Atoi(os.Getenv(BOT_OWNER.String()))
	if err != nil {
		log.Printf("WARN! %s not set or invalid ", BOT_OWNER)
	}
	botAdmins := os.Getenv(BOT_ADMINS.String())
	if botAdmins == "" {
		log.Printf("WARN! %s not set ", BOT_ADMINS)
	}
	var admins []int64
	for admin := range strings.SplitSeq(botAdmins, ",") {
		adminID, err := strconv.Atoi(admin)
		if err != nil {
			log.Printf("WARN! %s not set or invalid ", BOT_ADMINS)
			continue
		}
		admins = append(admins, int64(adminID))
	}
	botPlatformAPIKey := os.Getenv(BOT_PLATFORM_API_KEY.String())
	if botPlatformAPIKey == "" {
		log.Printf("WARN! %s not set ", BOT_PLATFORM_API_KEY)
	}
	e.botToken = botToken
	e.botOwner = int64(botOwner)
	e.botAdmins = admins
	e.botDebug = os.Getenv(BOT_DEBUG.String()) == "true"
	e.botVersion = os.Getenv(BOT_VERSION.String())
	e.botIsDevelopment = os.Getenv(BOT_IS_DEVELOPMENT.String()) == "true"
	e.botPlatformAPIKey = botPlatformAPIKey
}

func (e *Env) Get(key EnvKey) any {
	switch key {
	case BOT_TOKEN:
		return e.botToken
	case BOT_OWNER:
		return e.botOwner
	case BOT_ADMINS:
		return e.botAdmins
	case BOT_DEBUG:
		return e.botDebug
	case BOT_VERSION:
		return e.botVersion
	case BOT_IS_DEVELOPMENT:
		return e.botIsDevelopment
	default:
		log.Printf("WARN! %s not found in env ", key)
		return nil
	}
}

func (e *Env) GetString(key EnvKey) string {
	value := e.Get(key)
	if value == nil {
		return ""
	}
	strValue, ok := value.(string)
	if !ok {
		log.Printf("WARN! %s is not a string ", key)
		return ""
	}
	return strValue
}

func (e *Env) GetInt64(key EnvKey) int64 {
	value := e.Get(key)
	if value == nil {
		return 0
	}
	intValue, ok := value.(int64)
	if !ok {
		log.Printf("WARN! %s is not an int64 ", key)
		return 0
	}
	return intValue
}

func (e *Env) GetBool(key EnvKey) bool {
	value := e.Get(key)
	if value == nil {
		return false
	}
	boolValue, ok := value.(bool)
	if !ok {
		log.Printf("WARN! %s is not a bool ", key)
		return false
	}
	return boolValue
}
