module github.com/tooploox/oya

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.18.0+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/cucumber/godog v0.9.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-test/deep v1.0.1
	github.com/gobuffalo/helpers v0.2.2 // indirect
	github.com/gobuffalo/plush v3.8.2+incompatible
	github.com/gobwas/glob v0.2.3
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/google/uuid v1.1.0 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mozilla/mig v0.0.0-20190703170622-33eefe9c974e
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586 // indirect
	golang.org/x/net v0.0.0-20190813141303-74dc4d7220e7 // indirect
	golang.org/x/sys v0.0.0-20220517195934-5e4e11fc645e // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.0
	gopkg.in/src-d/go-git-fixtures.v3 v3.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.9.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/apimachinery v0.0.0-20181015213631-60666be32c5d
	k8s.io/helm v2.12.3+incompatible
	mvdan.cc/sh/v3 v3.0.0-alpha2.0.20190827105346-6af96bc17993
)

replace github.com/gobuffalo/plush => github.com/bart84ek/plush v3.8.2-oya+incompatible

replace gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.20.0
