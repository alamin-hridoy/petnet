module brank.as/RBAC/deploy/testsite

go 1.16

require (
	brank.as/rbac v0.0.0-20210526103834-182555235847
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/gorilla/csrf v1.7.0
	github.com/gorilla/schema v1.2.0
	github.com/gorilla/sessions v1.2.1
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/kenshaw/goji v0.2.0
	github.com/kenshaw/jwt v0.2.1 // indirect
	github.com/kenshaw/sentinel v1.0.2
	github.com/kenshaw/stringid v0.1.1 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.1
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914
	google.golang.org/grpc v1.44.0
)

replace brank.as/rbac => ../../
