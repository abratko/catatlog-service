package config

import (
	"fmt"
	"net"

	"gitlab.com/gorib/pry"
	"gitlab.com/gorib/pry/channels"
	"gitlab.com/gorib/server/host"
	"gitlab.com/gorib/try"
	grpcCommand "gitlab.com/gorib/waffle/tools/grpc"
	"gitlab.trgdev.com/gotrg/white-label/modules/health"
	healthProto "gitlab.trgdev.com/gotrg/white-label/modules/health/proto"
	"gitlab.trgdev.com/gotrg/white-label/services/catalog/config/dc"

	"gitlab.com/gorib/di"
	"gitlab.com/gorib/env"
)

func InitDi() {

	di.MustWire[pry.Logger](func() (pry.Logger, error) {
		dsn := env.Value("SENTRY_DSN", "")
		environment := fmt.Sprintf("%s_%s", env.NeedValue[string]("GITLAB_ENVIRONMENT_NAME"), env.NeedValue[string]("ENVIRONMENT"))
		return pry.New(
			env.Value("LOGLEVEL", "info"),
			pry.ToChannels(try.Throw(channels.Sentry("error", channels.WithSentryConnection(dsn, environment)))),
		)
	})

	di.Wire[net.Addr](host.NewListener, di.For[grpcCommand.GrpcCommand](), di.Defaults(map[int]any{
		0: "tcp",
		1: env.Value("DAEMON_HOST", "0.0.0.0"),
		2: env.Value("DAEMON_PORT", "8009"),
	}))

	di.Wire[healthProto.HealthServiceServer](
		health.NewHealthServiceServer,
		di.Tag(grpcCommand.ServiceTag),
	)

	di.Wire[grpcCommand.GrpcCommand](
		dc.SearchGrpcController.Use,
		di.Tag(grpcCommand.ServiceTag),
	)
	grpcCommand.InitGrpc()
}
