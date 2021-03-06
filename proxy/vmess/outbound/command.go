package outbound

import (
	"time"

	"v2ray.com/core/common/loader"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy/vmess"
)

func (this *VMessOutboundHandler) handleSwitchAccount(cmd *protocol.CommandSwitchAccount) {
	account := &vmess.Account{
		Id:      cmd.ID.String(),
		AlterId: uint32(cmd.AlterIds),
	}

	user := &protocol.User{
		Email:   "",
		Level:   cmd.Level,
		Account: loader.NewTypedSettings(account),
	}
	dest := v2net.TCPDestination(cmd.Host, cmd.Port)
	until := time.Now().Add(time.Duration(cmd.ValidMin) * time.Minute)
	this.serverList.AddServer(protocol.NewServerSpec(vmess.NewAccount, dest, protocol.BeforeTime(until), user))
}

func (this *VMessOutboundHandler) handleCommand(dest v2net.Destination, cmd protocol.ResponseCommand) {
	switch typedCommand := cmd.(type) {
	case *protocol.CommandSwitchAccount:
		if typedCommand.Host == nil {
			typedCommand.Host = dest.Address
		}
		this.handleSwitchAccount(typedCommand)
	default:
	}
}
