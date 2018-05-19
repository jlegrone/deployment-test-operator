# Deployment Test Operator

[![Operator Framework](https://img.shields.io/badge/Built%20With-Operator%20Framework-red.svg)](https://github.com/operator-framework/operator-sdk)
[![Proof of Concept](https://img.shields.io/badge/Status-Proof%20of%20Concept-yellow.svg)](https://github.com/jlegrone/deployment-test-operator)

[![CircleCI](https://circleci.com/gh/jlegrone/deployment-test-operator.svg?style=svg)](https://circleci.com/gh/jlegrone/deployment-test-operator)

This project is currently a proof of concept.  The goal is to declaratively specify tests for kubernetes [deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/) (and in the future potentially statefulsets), and to expose hooks for taking actions based on test results -- such as rolling back failed deployments, or promotion to higher environments.

## Usage

```bash
# Create DeploymentTest CRD, RBAC role, and operator deployment
kubectl apply -f deploy
# Create example Deployment, Service, and DeploymentTest
kubectl apply -f examples/nginx
```

## Development Workflow

Prerequisites:
  - locally running kubernetes or OpenShift cluster
  - [Operator SDK](https://github.com/operator-framework/operator-sdk#quick-start) command line tools installed

Instructions:

1. Deploy CRD and RBAC manifests:
    ```bash
    kubectl apply -f deploy/crd.yaml
    kubectl apply -f deploy/rbac.yaml
    ```

1. Create an example deployment:
    ```bash
    kubectl apply -f examples/nginx
    ```

1. On each change, run:
    ```bash
    WATCH_NAMESPACE=default operator-sdk up local
    ```

Publishing Updates:

```bash
# Build operator image
operator-sdk build quay.io/jlegrone/deployment-test-operator:<tag>
# Push to remote repository
docker push quay.io/jlegrone/deployment-test-operator:<tag>
```
