package plugin

import "github.com/hashicorp/go-plugin"

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BOTKUBE_MAGIC_COOKIE",
	MagicCookieValue: "BOTKUBE_BASIC_PLUGIN",
}
