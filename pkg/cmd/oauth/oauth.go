package oauth

import (
	"context"

	"github.com/jenkins-x/jx-helpers/v3/pkg/cobras/helper"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/spf13/cobra"

	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"

	// apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"golang.org/x/oauth2/clientcredentials"

)

// Options the options for updating jx-boot secrets
type Options struct {
	kubeClient    kubernetes.Interface
	secret        *v1.Secret
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
	var err error
	if o.kubeClient == nil {
		o.kubeClient, err = kube.LazyCreateKubeClient(o.kubeClient)
		if err != nil {
			return errors.Wrapf(err, "failed to create kube client")
		}
	}

	return nil
}

func (o *Options) Run() error {
	err := o.Validate()
	if err != nil {
		return errors.Wrapf(err, "failed to validate options")
	}

	oldSecret, err := o.kubeClient.
		CoreV1().
		Secrets("jx-git-operator").
		Get(context.Background(), "jx-boot", metav1.GetOptions{})

	o.secret, err = updateSecret(oldSecret)
	o.kubeClient.CoreV1().
		Secrets("jx-git-operator").
		Update(context.Background(), o.secret, metav1.UpdateOptions{})
	return nil
}

func updateSecret(s *v1.Secret) (*v1.Secret, error) {
	config := clientcredentials.Config{
		ClientID:     "",
		ClientSecret: "",
		TokenURL:     "https://bitbucket.org/site/oauth2/access_token",
	}

	s.Data["username"] = []byte("x-token-auth")
	token, err := config.Token(context.Background())
	if err == nil {
		s.Data["password"] = []byte(token.AccessToken)
	}
	return s, nil
}
