package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "time/tzdata" // bundle IANA tz data so LoadLocation works everywhere (CI, Render)
)

func main() {
	loadDotenv(".env")

	port := env("PORT", "8080")
	allowedOrigin := env("ALLOWED_ORIGIN", "http://localhost:5173")
	configPath := env("CONFIG_PATH", "config.json")
	messagesPath := env("MESSAGES_PATH", "reasons.json")
	defaultCountry := env("DEFAULT_COUNTRY", "ID")
	ratePerMin := envInt("RATE_LIMIT_PER_MINUTE", 30)
	nowOverride := env("NOW_OVERRIDE", "") // RFC3339; demo only — fake the clock

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := loadMessages(messagesPath)
	if err != nil {
		log.Fatal(err)
	}

	s := &server{
		cfg:            cfg,
		msgs:           msgs,
		lim:            newLimiter(ratePerMin),
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
		defaultCountry: defaultCountry,
		nowOverride:    nowOverride,
	}

	log.Printf("listening on :%s (allowed origin %s)", port, allowedOrigin)
	log.Fatal(http.ListenAndServe(":"+port, s.router(allowedOrigin)))
}
