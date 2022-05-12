package kubernetes

import (
	"context"

	"github.com/Coflnet/pr-controller/internal/model"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Create(pr *model.Pr) error {
	log.Info().Msgf("creating deployment for repo %s branch %s", pr.Repo, pr.Branch)
	err := CreateDeployment(pr)
	if err != nil {
		return err
	}

	err = CreateService(pr)
	if err != nil {
		return err
	}

	err = CreateIngress(pr)
	if err != nil {
		return err
	}

	return nil
}

func CreateDeployment(pr *model.Pr) error {

	deploymentsClient := clientset.AppsV1().Deployments("pr-env")
	var replicaCount int32 = 1

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pr.KubernetesResourceName(),
			Namespace: "pr-env",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": pr.KubernetesResourceName(),
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": pr.KubernetesResourceName(),
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  pr.KubernetesResourceName(),
							Image: pr.FullImageWithTag(),
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 3000,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "BASE_PATH",
									Value: pr.DomainPath(),
								},
								{
									Name:  "COMMAND_ENDPOINT",
									Value: "http://sky-commands:8008/command",
								},
								{
									Name:  "API_ENDPOINT",
									Value: "http://api-service/api",
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info().Msgf("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("error when creating kubernetes deployment")
		return err
	}
	log.Info().Msgf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func CreateService(pr *model.Pr) error {
	serviceClient := clientset.CoreV1().Services("pr-env")

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pr.KubernetesResourceName(),
			Namespace: "pr-env",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": pr.KubernetesResourceName(),
			},
			Ports: []apiv1.ServicePort{
				{
					TargetPort: intstr.FromInt(3000),
					Port:       3000,
				},
			},
		},
	}

	log.Info().Msgf("Creating service...")
	result, err := serviceClient.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("error when creating kubernetes service")
		return err
	}
	log.Info().Msgf("Created service %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func CreateIngress(pr *model.Pr) error {
	ingressClient := clientset.NetworkingV1().Ingresses("pr-env")

	var pathPrefix = v1.PathTypePrefix

	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pr.KubernetesResourceName(),
			Namespace: "pr-env",
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "nginx",
			},
		},
		Spec: v1.IngressSpec{
			Rules: []v1.IngressRule{
				{
					Host: pr.DomainName(),
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     pr.DomainPath(),
									PathType: &pathPrefix,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: pr.KubernetesResourceName(),
											Port: v1.ServiceBackendPort{
												Number: 3000,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	log.Info().Msgf("Creating service...")

	result, err := ingressClient.Create(context.TODO(), ingress, metav1.CreateOptions{})

	if err != nil {
		log.Error().Err(err).Msgf("error when creating kubernetes ingress")
		return err
	}

	log.Info().Msgf("Created ingress %q.\n", result.GetObjectMeta().GetName())
	return nil
}
