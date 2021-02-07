package check

import (
	"fmt"
	"strings"
	"sync"

	"github.com/fission/fission/pkg/fission-cli/cliwrapper/cli"
	"github.com/fission/fission/pkg/fission-cli/cmd"
	"github.com/fission/fission/pkg/fission-cli/cmd/check/resources"
	flagkey "github.com/fission/fission/pkg/fission-cli/flag/key"
	"github.com/fission/fission/pkg/fission-cli/util"
)

type CheckSubCommand struct {
	cmd.CommandActioner
}

func Check(input cli.Input) error {
	return (&CheckSubCommand{}).do(input)
}

func (opts *CheckSubCommand) do(input cli.Input) error {
	pre := input.Bool(flagkey.CheckPre)
	kubeContext := input.String(flagkey.KubeContext)

	mark := map[bool]string{
		true:  "✓",
		false: "✗",
	}

	_, k8sClient, err := util.GetKubernetesClient(kubeContext)
	if err != nil {
		return err
	}

	var checks []resources.Resource
	if pre {
		checks = []resources.Resource{
			resources.NewKubernetesVersion(k8sClient),
		}
	} else {
		checks = []resources.Resource{
			resources.NewFissionVersion(opts.Client()),
			resources.NewKubernetesPodStatus(k8sClient),
		}
	}

	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}

	// perform checks
	for _, res := range checks {
		wg.Add(1)
		go func(res resources.Resource) {
			results := res.Check()
			mtx.Lock()
			fmt.Printf("\n[%s]\n%s\n", res.GetLabel(), strings.Repeat("-", len(res.GetLabel())+2))
			for _, result := range results {
				fmt.Printf("%s %s\n", mark[result.Ok], result.Description)
			}
			mtx.Unlock()
			wg.Done()
		}(res)
	}
	wg.Wait()
	return nil
}
