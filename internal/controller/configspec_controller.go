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
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	configv1 "github.com/ajayp/config-operator/api/v1"
	"sigs.k8s.io/yaml"
)

// ConfigSpecReconciler reconciles a ConfigSpec object
type ConfigSpecReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=config.ajay.dev,resources=configspecs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ajay.dev,resources=configspecs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ajay.dev,resources=configspecs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=ajay.dev,resources=configspecs/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// Behavior:
// - Create/Update a ConfigMap named "<cr-name>-config" with the CR's spec.value in key "config.yaml".
// - Maintain spec.time (RFC3339) on writes and set spec.status to Applied/Error accordingly.
// - Delete the associated ConfigMap when the CR is deleted.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *ConfigSpecReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)

	// Fetch the ConfigSpec instance
	var instance configv1.ConfigSpec
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		if apierrors.IsNotFound(err) {
			// CR deleted; nothing to do
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Define finalizer name
	const finalizerName = "configspec.finalizers.ajay.dev"

	// Handle deletion
	if !instance.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(&instance, finalizerName) {
			// Delete associated ConfigMap
			cmName := fmt.Sprintf("%s-config", instance.Name)
			cm := &corev1.ConfigMap{}
			if err := r.Get(ctx, types.NamespacedName{Name: cmName, Namespace: instance.Namespace}, cm); err == nil {
				_ = r.Delete(ctx, cm)
			}
			// Remove finalizer and update
			controllerutil.RemoveFinalizer(&instance, finalizerName)
			if err := r.Update(ctx, &instance); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Ensure finalizer present
	if !controllerutil.ContainsFinalizer(&instance, finalizerName) {
		controllerutil.AddFinalizer(&instance, finalizerName)
		if err := r.Update(ctx, &instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Ensure ConfigMap exists/updated
	cmName := fmt.Sprintf("%s-config", instance.Name)
	// Validate YAML in spec.Value
	var validated interface{}
	if err := yaml.Unmarshal([]byte(instance.Spec.Value), &validated); err != nil {
		logger.Error(err, "invalid YAML in spec.value")
		_ = r.updateSpecStatus(ctx, &instance, "Error")
		return ctrl.Result{}, err
	}

	desiredData := map[string]string{"config.yaml": instance.Spec.Value}

	var cm corev1.ConfigMap
	err := r.Get(ctx, types.NamespacedName{Name: cmName, Namespace: instance.Namespace}, &cm)
	if apierrors.IsNotFound(err) {
		// Create
		cm = corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmName,
				Namespace: instance.Namespace,
			},
			Data: desiredData,
		}
		if err := controllerutil.SetControllerReference(&instance, &cm, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, &cm); err != nil {
			logger.Error(err, "failed to create ConfigMap")
			_ = r.updateSpecStatus(ctx, &instance, "Error")
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// Update if needed
		if cm.Data == nil || cm.Data["config.yaml"] != instance.Spec.Value {
			if cm.Data == nil {
				cm.Data = map[string]string{}
			}
			cm.Data["config.yaml"] = instance.Spec.Value
			if err := r.Update(ctx, &cm); err != nil {
				logger.Error(err, "failed to update ConfigMap")
				_ = r.updateSpecStatus(ctx, &instance, "Error")
				return ctrl.Result{}, err
			}
		}
	}

	// Update spec.time and spec.status to Applied
	if err := r.updateSpecApplied(ctx, &instance); err != nil {
		return ctrl.Result{}, err
	}

	// Optionally update observed status info
	now := metav1.NewTime(time.Now().UTC())
	if instance.Status.LastAppliedTime == nil || instance.Status.LastAppliedTime.Time.Before(now.Time) {
		instance.Status.LastAppliedTime = &now
		instance.Status.ConfigMapName = cmName
		if err := r.Status().Update(ctx, &instance); err != nil {
			// Requeue on conflict or retryable errors
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigSpecReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configv1.ConfigSpec{}).
		Owns(&corev1.ConfigMap{}).
		Named("configspec").
		Complete(r)
}

// updateSpecApplied sets spec.time to now and spec.status to Applied, persisting via Update.
func (r *ConfigSpecReconciler) updateSpecApplied(ctx context.Context, instance *configv1.ConfigSpec) error {
	instance.Spec.Time = time.Now().UTC().Format(time.RFC3339)
	instance.Spec.Status = "Applied"
	return r.Update(ctx, instance)
}

// updateSpecStatus sets spec.status and updates the resource.
func (r *ConfigSpecReconciler) updateSpecStatus(ctx context.Context, instance *configv1.ConfigSpec, status string) error {
	instance.Spec.Status = status
	return r.Update(ctx, instance)
}
