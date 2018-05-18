package stub

import (
	"github.com/jlegrone/deployment-test-operator/pkg/apis/deploy/v1alpha1"

	"encoding/json"
	"fmt"

	k8sclient "github.com/operator-framework/operator-sdk/pkg/k8sclient"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx sdk.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *appsv1.Deployment:
		return deploymentHandler(o)
	case *batchv1.Job:
		return jobHandler(o)
	}
	return nil
}

func deploymentHandler(deployment *appsv1.Deployment) error {
	deploymentRevision, hasDeploymentRevision := deployment.Annotations["deployment.kubernetes.io/revision"]
	deploymentTestName, hasDeploymentTestAnnotation := deployment.Annotations["k8s.jacob.work/deployment-test-name"]
	testStatusAnnotationKey := "k8s.jacob.work/deployment-test-status-revision-" + deploymentRevision
	_, testsInitialized := deployment.Annotations[testStatusAnnotationKey]

	if hasDeploymentTestAnnotation && hasDeploymentRevision && !testsInitialized {
		client, _, clientErr := k8sclient.GetResourceClient("deploy.k8s.jacob.work/v1alpha1", "DeploymentTest", deployment.Namespace)
		if clientErr != nil {
			return clientErr
		}
		deploymentTestUnstruc, getErr := client.Get(deploymentTestName, metav1.GetOptions{})
		if getErr != nil {
			return getErr
		}
		depJSON, _ := deploymentTestUnstruc.MarshalJSON()

		var deploymentTest *v1alpha1.DeploymentTest
		jsonErr := json.Unmarshal([]byte(depJSON), &deploymentTest)
		if jsonErr != nil {
			return jsonErr
		}

		// Create the Job from DeploymentTest spec
		createErr := sdk.Create(newDeploymentTestJob(deploymentTest, deployment))
		if createErr != nil && !errors.IsAlreadyExists(createErr) {
			return createErr
		}

		// Add annotation to Deployment
		deployment.Annotations[testStatusAnnotationKey] = "Pending"
		updateErr := sdk.Update(deployment)
		if updateErr != nil {
			return updateErr
		}

		logrus.Info(fmt.Sprint(
			deployment.Name,
			" revision ",
			deploymentRevision,
			" test started.",
		))

		eventErr := sdk.Create(newEvent(
			deployment.Namespace,
			"CreatedDeploymentTest", fmt.Sprint(
				"Started deployment test ",
				deploymentTest.Name,
				" for revision ",
				deploymentRevision,
			),
			"Normal",
			newDeploymentReference(deployment),
		))
		if eventErr != nil {
			return eventErr
		}
	}
	return nil
}

func jobHandler(job *batchv1.Job) error {
	testStatusReportedStr, isDeploymentTest := job.Annotations["k8s.jacob.work/deployment-test-status-reported"]

	if isDeploymentTest {
		statusReported := testStatusReportedStr == "True"
		if !statusReported && len(job.Status.Conditions) > 0 {
			return processJob(job)
		}
	}
	return nil
}

func processJob(job *batchv1.Job) error {
	for _, v := range job.Status.Conditions {
		if v.Type == "Complete" {
			notifyTestResult(job, true)
			break
		}
		if v.Type == "Failed" {
			notifyTestResult(job, false)
			break
		}
	}
	return nil
}

func notifyTestResult(job *batchv1.Job, succeeded bool) error {
	deploymentName := job.Annotations["k8s.jacob.work/deployment-name"]
	deploymentRevision := job.Annotations["k8s.jacob.work/deployment-revision"]
	deploymentTestName := job.Annotations["k8s.jacob.work/deployment-test-name"]
	testStatusAnnotationKey := "k8s.jacob.work/deployment-test-status-revision-" + deploymentRevision

	client, _, clientErr := k8sclient.GetResourceClient("apps/v1", "Deployment", job.Namespace)
	if clientErr != nil {
		return clientErr
	}
	deploymentUnstruc, getErr := client.Get(deploymentName, metav1.GetOptions{})
	if getErr != nil {
		return getErr
	}
	depJSON, _ := deploymentUnstruc.MarshalJSON()

	var deployment *appsv1.Deployment
	jsonErr := json.Unmarshal([]byte(depJSON), &deployment)
	if jsonErr != nil {
		return jsonErr
	}

	if succeeded {
		// Send event to Deployment
		sdk.Create(newEvent(
			deployment.Namespace,
			"SuccessfulDeploymentTest", fmt.Sprint(
				"Deployment test ",
				deploymentTestName,
				" for revision ",
				deploymentRevision,
				" completed successfully",
			),
			"Normal",
			newDeploymentReference(deployment),
		))
		logrus.Info(fmt.Sprint(
			deploymentName,
			" revision ",
			deploymentRevision,
			" test succeeded.",
		))
		deployment.Annotations[testStatusAnnotationKey] = "Complete"
		sdk.Update(deployment)
		job.Annotations["k8s.jacob.work/deployment-test-status-reported"] = "True"
		sdk.Update(job)
	} else {
		// Send event to Deployment
		sdk.Create(newEvent(
			deployment.Namespace,
			"FailedDeploymentTest", fmt.Sprint(
				"Deployment test ",
				deploymentTestName,
				" for revision ",
				deploymentRevision,
				" failed",
			),
			"Warning",
			newDeploymentReference(deployment),
		))
		logrus.Error(fmt.Sprint(
			deploymentName,
			" revision ",
			deploymentRevision,
			" test failed.",
		))
		deployment.Annotations[testStatusAnnotationKey] = "Failed"
		sdk.Update(deployment)
		job.Annotations["k8s.jacob.work/deployment-test-status-reported"] = "True"
		sdk.Update(job)
	}

	return nil
}

func newDeploymentReference(deployment *appsv1.Deployment) v1.ObjectReference {
	return v1.ObjectReference{
		APIVersion:      deployment.APIVersion,
		Kind:            deployment.Kind,
		Namespace:       deployment.Namespace,
		Name:            deployment.Name,
		ResourceVersion: deployment.ResourceVersion,
		UID:             deployment.UID,
	}
}

func newDeploymentTestJob(cr *v1alpha1.DeploymentTest, target *appsv1.Deployment) *batchv1.Job {
	deploymentRevision := target.Annotations["deployment.kubernetes.io/revision"]
	labels := map[string]string{
		"deployment-test": "true",
	}
	annotations := map[string]string{
		"k8s.jacob.work/deployment-name":                 target.Name,
		"k8s.jacob.work/deployment-revision":             deploymentRevision,
		"k8s.jacob.work/deployment-test-name":            cr.Name,
		"k8s.jacob.work/deployment-test-status-reported": "False",
	}
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprint(target.Name, "-test-revision-", deploymentRevision),
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "DeploymentTest",
				}),
			},
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: cr.Spec.JobTemplate,
	}
}

func newEvent(namespace string, reason string, message string, eventType string, involvedObject v1.ObjectReference) *v1.Event {
	return &v1.Event{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Event",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    namespace,
			GenerateName: "deployment-test-lifecycle-",
		},
		InvolvedObject: involvedObject,
		Reason:         reason,
		Message:        message,
		Type:           eventType,
		FirstTimestamp: metav1.Now(),
		LastTimestamp:  metav1.Now(),
		Source: v1.EventSource{
			Component: "deployment-test-operator",
		},
	}
}
