package main

import (
	k8s "example.com/cdk8s-test/imports/k8s"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

// NginxChart defines the nginx deployment and service
func NewNginxChart(scope constructs.Construct) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String("nginx-deployment"), &cdk8s.ChartProps{})

	nginxLabels := map[string]*string{
		"app": jsii.String("nginx"),
	}

	var replicas int32 = 3
	k8s.NewKubeDeployment(chart, jsii.String("nginx-deployment"), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name:   jsii.String("nginx-deployment"), // Keep the name consistent
			Labels: &nginxLabels,
		},
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(replicas),
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

	return chart
}

// HeadlampChart defines the Headlamp deployment and service
func NewHeadlampChart(scope constructs.Construct) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String("headlamp-deployment"), &cdk8s.ChartProps{})

	headlampLabels := map[string]*string{
		"app": jsii.String("headlamp"),
	}

	var replicas int32 = 3 // As defined in TS code
	headlampDeploymentName := "headlamp-deployment"

	k8s.NewKubeDeployment(chart, jsii.String(headlampDeploymentName), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name:   jsii.String(headlampDeploymentName),
			Labels: &headlampLabels,
		},
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(replicas),
			Selector: &k8s.LabelSelector{
				MatchLabels: &headlampLabels,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &headlampLabels,
				},
				Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{
						{
							Name:  jsii.String("headlamp"),
							Image: jsii.String("ghcr.io/headlamp-k8s/headlamp:latest"),
							Args:  &[]*string{jsii.String("-in-cluster"), jsii.String("-plugins-dir=/headlamp/plugins")},
							Env: &[]*k8s.EnvVar{
								{Name: jsii.String("HEADLAMP_CONFIG_TRACING_ENABLED"), Value: jsii.String("true")},
								{Name: jsii.String("HEADLAMP_CONFIG_METRICS_ENABLED"), Value: jsii.String("true")},
								{Name: jsii.String("HEADLAMP_CONFIG_OTLP_ENDPOINT"), Value: jsii.String("otel-collector:4317")},
								{Name: jsii.String("HEADLAMP_CONFIG_SERVICE_NAME"), Value: jsii.String("headlamp")},
								{Name: jsii.String("HEADLAMP_CONFIG_SERVICE_VERSION"), Value: jsii.String("latest")},
							},
							Ports: &[]*k8s.ContainerPort{
								{ContainerPort: jsii.Number(4466), Name: jsii.String("http")},
								{ContainerPort: jsii.Number(9090), Name: jsii.String("metrics")},
							},
						},
					},
				},
			},
		},
	})

	// Create Headlamp service
	var port int32 = 80
	var targetPort int32 = 4466 // As defined in TS code
	k8s.NewKubeService(chart, jsii.String("headlamp-service"), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("headlamp-service"),
			Namespace: jsii.String("default"), // As defined in TS code
		},
		Spec: &k8s.ServiceSpec{
			Selector: &headlampLabels,
			Ports: &[]*k8s.ServicePort{
				{
					Port:       jsii.Number(port),
					TargetPort: k8s.IntOrString_FromNumber(jsii.Number(targetPort)),
				},
			},
		},
	})

	// Create Headlamp secret
	k8s.NewKubeSecret(chart, jsii.String("headlamp-admin"), &k8s.KubeSecretProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("headlamp-admin"),
			Namespace: jsii.String("default"), // As defined in TS code
			Annotations: &map[string]*string{
				"kubernetes.io/service-account.name": jsii.String("headlamp-admin"),
			},
		},
		Type: jsii.String("kubernetes.io/service-account-token"),
	})

	return chart
}

func NewKustomizationChart(scope constructs.Construct) cdk8s.Chart {
	chart := cdk8s.NewChart(scope, jsii.String("kustomization-deployment"), &cdk8s.ChartProps{})

	cdk8s.NewApiObject(chart, jsii.String("kustomization-deployment"), &cdk8s.ApiObjectProps{
		ApiVersion: jsii.String("kustomize.config.k8s.io/v1beta1"),
		Kind:       jsii.String("Kustomization"),
	})

	return chart
}

func main() {
	app := cdk8s.NewApp(&cdk8s.AppProps{})

	// Create separate charts for different components
	NewNginxChart(app)
	NewHeadlampChart(app)
	NewKustomizationChart(app)

	// Synthesize everything
	app.Synth()
}
