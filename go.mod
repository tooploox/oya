module github.com/tooploox/oya

go 1.16

replace github.com/gobuffalo/plush => github.com/bart84ek/plush v3.8.2-oya+incompatible

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.18.0+incompatible
	github.com/blang/semver v3.5.1+incompatible
	github.com/c-bata/go-prompt v0.2.4-0.20200321140817-d043be076398
	github.com/cucumber/godog v0.12.2
	github.com/cucumber/messages-go/v16 v16.0.1
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/go-test/deep v1.0.1
	github.com/gobuffalo/flect v0.2.2 // indirect
	github.com/gobuffalo/helpers v0.6.1 // indirect
	github.com/gobuffalo/plush v3.8.2+incompatible
	github.com/gobuffalo/tags v2.1.0+incompatible // indirect
	github.com/gobuffalo/validate/v3 v3.3.0 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.2 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/microcosm-cc/bluemonday v1.0.6 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mozilla/mig v0.0.0-20190703170622-33eefe9c974e
	github.com/pkg/errors v0.8.1
	github.com/pkg/term v0.0.0-20200520122047-c3ffed290a03 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	gopkg.in/src-d/go-billy.v4 v4.3.0
	gopkg.in/src-d/go-git-fixtures.v3 v3.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.9.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.0.0-20181015213631-60666be32c5d
	k8s.io/helm v2.12.3+incompatible
	mvdan.cc/sh/v3 v3.1.1
)

replace gopkg.in/urfave/cli.v1 => github.com/urfave/cli v1.20.0
