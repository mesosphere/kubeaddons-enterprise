module github.com/mesosphere/kubeaddons-enterprise/test

go 1.13

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/docker/docker v1.4.2-0.20190916154449-92cc603036dd
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/mesosphere/ksphere-testing-framework v0.0.0-20200320204306-f29e7880920f
	github.com/mesosphere/kubeaddons v0.13.0
	go.uber.org/atomic v1.5.1 // indirect
	go.uber.org/multierr v1.4.0 // indirect
	go.uber.org/zap v1.13.0 // indirect
	golang.org/x/crypto v0.0.0-20200117160349-530e935923ad // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20200122134326-e047566fdf82 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	golang.org/x/tools v0.0.0-20200122042241-dc16b66866f1 // indirect
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200121175148-a6ecf24a6d71 // indirect
	k8s.io/apimachinery v0.17.4
	k8s.io/utils v0.0.0-20200117235808-5f6fbceb4c31 // indirect
	sigs.k8s.io/kind v0.7.0
)

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
