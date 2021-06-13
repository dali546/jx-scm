package oauth

import (
	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"
	"github.com/jenkins-x/jx-git-operator/pkg/constants"

	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	// apierrors "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Options the options for updating jx-boot secrets
type Options struct {
	kubeClient kubernetes.Interface
	ns         string
	selector   string
}

// NewCmdUpdateSecret refreshes the jx-boot secret
func NewCmdUpdateSecret() (*cobra.Command, *Options) {
	o := &Options{}

	command := &cobra.Command{
		Use:   "update-secret",
		Short: "Refresh Git Authentication Secret",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run()
		},
	}
	return command, o
}

func (o *Options) Validate() error {
	if o.kubeClient == nil {
		f := kubeclient.NewFactory()
		cfg, err := f.CreateKubeConfig()
		if err != nil {
			return errors.Wrapf(err, "failed to create kube config")
		}

		o.kubeClient, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			return errors.Wrapf(err, "failed to create the kube client")
		}

		if o.ns == "" {
			o.ns, err = kubeclient.CurrentNamespace()
			if err != nil {
				return errors.Wrapf(err, "failed to find the current namespace")
			}
		}
	}
	return nil
}

func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}
	return nil
}
