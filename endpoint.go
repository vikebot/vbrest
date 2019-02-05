package main

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"strings"

	"go.uber.org/zap"
)

type endpoint struct {
	Name       string
	Handler    handler
	ExactMatch bool
}

func allEndpoints(log *zap.Logger) []endpoint {
	buf := make([]byte, 128)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		log.Fatal("unable to generate genjwtkey secrete key", zap.Error(err))
	}

	secret := strings.ToLower(base32.StdEncoding.EncodeToString(buf))
	log.Warn("secret generated for genjwtkey",
		zap.String("secret", secret),
		zap.String("link", "/v0/admin/genjwtkey/"+secret+"/<DESIRED-USER-ID>"))

	// save it
	genjwtsecret = secret

	return []endpoint{
		{"/v0/admin/genjwtkey/" + secret + "/", v0AdminGenjwtkey, false},

		{"/v1/test", v1Test, true},
		{"/v1/user/get", v1UserGet, true},
		{"/v1/user/get/id/", v1UserGetPublicByID, false},
		{"/v1/user/get/username/", v1UserGetPublicByUsername, false},
		{"/v1/user/update", v1UserUpdate, true},
		{"/v1/round/active", v1RoundActive, true},
		{"/v1/round/join/", v1RoundJoin, false},
		{"/v1/roundentry/active", v1RoundentryActive, true},
		{"/v1/roundentry/connectinfo/", v1RoundentryConnectinfo, false},
		{"/v1/roundentry/watchresolve/", v1RoundentryWatchresolve, false},
		{"/v1/register/confirm", v1RegisterConfirm, false},
		{"/v1/watch/userinfo/", v1UserGetByWatchtoken, false},
		{"/v1/watch/players/", v1RoundPlayers, false},
	}
}
