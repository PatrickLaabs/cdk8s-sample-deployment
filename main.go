package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"example.com/cdk8s-test/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// func NewChart(scope constructs.Construct, id string, ns string, appLabel string) cdk8s.Chart {

// 	chart := cdk8s.NewChart(scope, jsii.String(id), &cdk8s.ChartProps{
// 		Namespace: jsii.String(ns),
// 	})

// 	labels := map[string]*string{
// 		"app": jsii.String(appLabel),
// 	}

// 	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
// 		Spec: &k8s.DeploymentSpec{
// 			Replicas: jsii.Number(1),
// 			Selector: &k8s.LabelSelector{
// 				MatchLabels: &labels,
// 			},
// 			Template: &k8s.PodTemplateSpec{
// 				Metadata: &k8s.ObjectMeta{
// 					Labels: &labels,
// 				},
// 				Spec: &k8s.PodSpec{
// 					Containers: &[]*k8s.Container{{
// 						Name:  jsii.String("app-container"),
// 						Image: jsii.String("nginx:1.19.10"),
// 						Ports: &[]*k8s.ContainerPort{{
// 							ContainerPort: jsii.Number(80),
// 						}},
// 					}},
// 				},
// 			},
// 		},
// 	})

// 	return chart
// }

// func main() {
// 	app := cdk8s.NewApp(nil)

// 	NewChart(app, "getting-started", "default", "my-app")

// 	app.Synth()
// }

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	// Create nginx deployment
	nginxLabels := map[string]*string{
		"app": jsii.String("nginx"),
	}

	k8s.NewKubeDeployment(chart, jsii.String("nginx-deployment"), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String("nginx-deployment"),
		},
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(1),
			Selector: &k8s.LabelSelector{
				MatchLabels: &nginxLabels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &nginxLabels,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{
						{
							Name:  jsii.String("nginx"),
							Image: jsii.String("nginx:latest"),
							Ports: &[]*k8s.ContainerPort{
								{
									ContainerPort: jsii.Number(80),
								},
							},
						},
					},
				},
			},
		},
	})

	// Create nginx service
	var port int32 = 80
	var targetPort int32 = 80
	k8s.NewKubeService(chart, jsii.String("nginx-service"), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String("nginx-service"),
		},
		Spec: &k8s.ServiceSpec{
			Selector: &nginxLabels,
			Ports: &[]*k8s.ServicePort{
				{
					Port:       jsii.Number(port),
					TargetPort: k8s.IntOrString_FromNumber(jsii.Number(targetPort)),
				},
			},
			Type: jsii.String("ClusterIP"),
		},
	})

	// // Create kustomization.yaml
	// kustomization := cdk8s.NewApiObject(chart, jsii.String("kustomization"), &cdk8s.ApiObjectProps{
	// 	ApiVersion: jsii.String("kustomize.config.k8s.io/v1beta1"),
	// 	Kind:       jsii.String("Kustomization"),
	// 	Metadata: &cdk8s.ApiObjectMetadata{
	// 		Name: jsii.String("kustomization"),
	// 	},
	// })

	// kustomization.AddJsonPatch(cdk8s.JsonPatch_Add(jsii.String("/resources"), &[]interface{}{
	// 	"nginx-deployment.yaml",
	// 	"nginx-service.yaml",
	// 	"headlamp-deployment.yaml",
	// }))

	return chart
}

func downloadHeadlampManifest(outputPath string) error {
	// Create dist directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	// Download the Headlamp manifest
	headlampURL := "https://raw.githubusercontent.com/kubernetes-sigs/headlamp/main/kubernetes-headlamp.yaml"
	resp, err := http.Get(headlampURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the content
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "cdk8s-sample", nil)
	app.Synth()

	// After synthesis, download the Headlamp manifest
	outputPath := "./dist/headlamp-deployment.yaml"
	if err := downloadHeadlampManifest(outputPath); err != nil {
		fmt.Printf("Error downloading Headlamp manifest: %v\n", err)
	} else {
		fmt.Printf("Headlamp deployment manifest saved to %s\n", outputPath)
	}
}
