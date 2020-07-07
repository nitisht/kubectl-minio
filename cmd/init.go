/*
 * Copyright (C) 2020, MinIO, Inc.
 *
 * This code is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License, version 3,
 * as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License, version 3,
 * along with this program.  If not, see <http://www.gnu.org/licenses/>
 *
 */

package cmd

import (
	"errors"
	"io"

	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/kudobuilder/kudo/pkg/kudoctl/clog"
	"github.com/kudobuilder/kudo/pkg/kudoctl/kube"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"
)

const (
	initDesc = `
	This command installs MinIO Operator onto your Kubernetes cluster. It discovers Kubernetes clusters by reading $KUBECONFIG (default '~/.kube/config') 
	and using the default context. When installing  MinIO Operator, 'minio init' will attempt to install the latest released version. You can specify an 
	alternative image with '--image' which is the fully qualified image name replacement. To dump a manifest containing the deployment YAML, combine the 
	'--dry-run' and '--o' flags.
	`
	initExample = `kubectl minio init`
)

type initCmd struct {
	out            io.Writer
	errOut         io.Writer
	image          string
	dryRun         bool
	output         bool
	nsToWatch      string
	clusterDomain  string
	serviceAccount string
}

func newInitCmd(out io.Writer, errOut io.Writer, client *kube.Client) *cobra.Command {
	i := &initCmd{out: out, errOut: errOut}

	cmd := &cobra.Command{
		Use:     "init",
		Short:   "Initialize MinIO Operator on the kubernetes cluster",
		Long:    initDesc,
		Example: initExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("this command does not accept arguments")
			}
			if err := i.validate(cmd.Flags()); err != nil {
				return err
			}
			return i.run()
		},
	}

	f := cmd.Flags()
	f.StringVarP(&i.image, "image", "i", "", "Override MinIO Operator image")
	f.StringVarP(&i.nsToWatch, "namespace-to-watch", "", "", "Namespace where MinIO Operator looks for MinIO Instances")
	f.StringVarP(&i.serviceAccount, "service-account", "", "", "Override for the default serviceAccount kudo-manager")
	f.StringVarP(&i.clusterDomain, "cluster-domain", "d", "", "Cluster domain ")
	f.BoolVar(&i.dryRun, "dry-run", false, "Do not install")
	f.BoolVar(&i.output, "output", false, "Output the yaml to be used for this command")

	return cmd
}

// TODO: Add validation for flags
func (initCmd *initCmd) validate(flags *pflag.FlagSet) error {
	return nil
}

// run initializes local config and installs MinIO Operator to Kubernetes cluster.
func (initCmd *initCmd) run() error {
	// initialize server
	clog.V(4).Printf("initializing server")

	if !initCmd.dryRun {
		config, err := GetKubeClient()
		if err != nil {
			klog.Errorf("could not get Kubernetes client: %s", err.Error())
			return nil
		}

		kubeClient, err := apiextension.NewForConfig(config)
		if err != nil {
			return nil
		}

		crd := embeddedCRD()
		if _, err = kubeClient. .ApiextensionsV1().CustomResourceDefinitions().Create(crd); err != nil {
			if apierrors.IsAlreadyExists(err) {
				clog.V(4).Printf("crd %v already exists", crd.Name)
				return nil
			}
			klog.Errorf("Error in creating CRD: %s", err.Error())
			return nil
		}
	}

	return nil
}
