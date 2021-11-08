module code.vegaprotocol.io/vegawallet

go 1.17

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.12.0
	github.com/spf13/cobra => github.com/spf13/cobra v1.2.1
	github.com/spf13/viper => github.com/spf13/viper v1.8.1
	gopkg.in/ini.v1 => github.com/go-ini/ini v1.63.2
)

require (
	code.vegaprotocol.io/protos v0.43.1-0.20211004102311-d2ad34ada37b
	code.vegaprotocol.io/shared v0.0.0-20211015074835-9ed837d93090
	github.com/blang/semver/v4 v4.0.0
	github.com/cenkalti/backoff/v4 v4.1.1
	github.com/golang-jwt/jwt/v4 v4.1.0
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.2
	github.com/julienschmidt/httprouter v1.3.0
	github.com/mattn/go-isatty v0.0.14
	github.com/muesli/termenv v0.9.0
	github.com/oasisprotocol/curve25519-voi v0.0.0-20210716083614-f38f8e8b0b84
	github.com/rs/cors v1.7.0
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/vegaprotocol/go-slip10 v0.1.0
	github.com/zannen/toml v0.3.2
	go.uber.org/zap v1.19.1
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
	google.golang.org/grpc v1.40.0
)

require (
	github.com/adrg/xdg v0.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mwitkow/go-proto-validators v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
	golang.org/x/sys v0.0.0-20210917161153-d61c044b1678 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210917145530-b395a37504d4 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
