package keystone

import (
        comv1 "github.com/openstack-k8s-operators/keystone-operator/pkg/apis/keystone/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Deployment(cr *comv1.KeystoneApi, cmName string) *appsv1.Deployment {
	runAsUser := int64(0)

	labels := map[string]string{
		"app": "keystone-api",
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Replicas: &cr.Spec.Replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "keystone-operator",
					Containers: []corev1.Container{
						{
							Name:  "keystone-api",
							Image: cr.Spec.ContainerImage,
							SecurityContext: &corev1.SecurityContext{
								RunAsUser: &runAsUser,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "KOLLA_CONFIG_STRATEGY",
									Value: "COPY_ALWAYS",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/var/lib/config-data",
									ReadOnly:  true,
									Name:      "config-data",
								},
								{
									MountPath: "/var/lib/kolla/config_files",
									ReadOnly:  true,
									Name:      "kolla-config",
								},
								{
									MountPath: "/var/lib/fernet-keys",
									ReadOnly:  true,
									Name:      "fernet-keys",
								},
							},
						},
					},
				},
			},
		},
	}
	deployment.Spec.Template.Spec.Volumes = getVolumes(cmName)
	return deployment
}
