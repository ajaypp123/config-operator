/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv1 "github.com/ajayp/config-operator/api/v1"
)

var _ = Describe("ConfigSpec Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		configspec := &configv1.ConfigSpec{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind ConfigSpec")
			err := k8sClient.Get(ctx, typeNamespacedName, configspec)
			if err != nil && errors.IsNotFound(err) {
				resource := &configv1.ConfigSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: configv1.ConfigSpecSpec{
						Value:  "key: value\n",
						Time:   time.Now().UTC().Format(time.RFC3339),
						Status: "Pending",
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// Cleanup logic after each test, ignore if already deleted
			resource := &configv1.ConfigSpec{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err != nil {
				Expect(errors.IsNotFound(err)).To(BeTrue())
				return
			}
			By("Cleanup the specific resource instance ConfigSpec")
			Expect(ctrlclient.IgnoreNotFound(k8sClient.Delete(ctx, resource))).To(Succeed())
		})
		It("should successfully reconcile create/update/delete", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ConfigSpecReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// Assert ConfigMap created
			cm := &corev1.ConfigMap{}
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-config", Namespace: "default"}, cm)).To(Succeed())
			Expect(cm.Data["config.yaml"]).To(Equal("key: value\n"))

			// Update the CR's spec and reconcile again
			res := &configv1.ConfigSpec{}
			Expect(k8sClient.Get(ctx, typeNamespacedName, res)).To(Succeed())
			res.Spec.Value = "key: new\n"
			Expect(k8sClient.Update(ctx, res)).To(Succeed())

			_, err = controllerReconciler.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
			Expect(err).NotTo(HaveOccurred())

			// ConfigMap should be updated
			Expect(k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-config", Namespace: "default"}, cm)).To(Succeed())
			Expect(cm.Data["config.yaml"]).To(Equal("key: new\n"))

			// Delete the CR and ensure CM removed through finalizer
			Expect(k8sClient.Delete(ctx, res)).To(Succeed())
			// Simulate reconcile on deletion
			_, _ = controllerReconciler.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})

			// ConfigMap should be gone eventually
			err = k8sClient.Get(ctx, types.NamespacedName{Name: resourceName + "-config", Namespace: "default"}, &corev1.ConfigMap{})
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})
})
