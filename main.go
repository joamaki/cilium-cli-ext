package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/cilium/cilium-cli/cli"
	"github.com/cilium/cilium-cli/connectivity/check"
)

func main() {
	hooks := &myHooks{}

	if err := cli.NewCiliumCommand(hooks).Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type myHooks struct {
	cli.NopHooks
}

// AddConnectivityTestFlags implements cli.Hooks
func (*myHooks) AddConnectivityTestFlags(flags *pflag.FlagSet) {
	fmt.Println(">>> AddConnectivityTestFlags")
}

// AddConnectivityTests implements cli.Hooks
func (*myHooks) AddConnectivityTests(ct *check.ConnectivityTest) error {
	fmt.Println(">>> AddConnectivityTests")
	ct.NewTest("dummies").WithScenarios(
		&dummy{},
	)
	return nil
}

// dummy implements a Scenario.
type dummy struct {
	name string
}

func Dummy(name string) check.Scenario {
	return &dummy{
		name: name,
	}
}

func (s *dummy) Name() string {
	tn := "dummy"
	if s.name == "" {
		return tn
	}
	return fmt.Sprintf("%s:%s", tn, s.name)
}

func (s *dummy) Run(ctx context.Context, t *check.Test) {
	t.NewAction(s, "action-1", nil, nil, check.IPFamilyAny).Run(func(a *check.Action) {
		a.Log("logging")
		a.Debug("debugging")
		a.Info("informing")
	})

	t.NewAction(s, "action-2", nil, nil, check.IPFamilyAny).Run(func(a *check.Action) {
		a.Log("logging")
		a.Fatal("killing :(")
		a.Fail("failing (this should not be printed)")
	})
}
