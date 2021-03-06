package gce

import (
	"context"
	"google.golang.org/api/googleapi"
	"net/http"
	"time"

	"github.com/supergiant/control/pkg/workflows/steps"
	"google.golang.org/api/compute/v1"
)

type computeService struct {
	getFromFamily       func(context.Context, steps.GCEConfig) (*compute.Image, error)
	getMachineTypes     func(context.Context, steps.GCEConfig) (*compute.MachineType, error)
	insertInstance      func(context.Context, steps.GCEConfig, *compute.Instance) (*compute.Operation, error)
	getInstance         func(context.Context, steps.GCEConfig, string) (*compute.Instance, error)
	setInstanceMetadata func(context.Context, steps.GCEConfig, string, *compute.Metadata) (*compute.Operation, error)
	deleteInstance      func(string, string, string) (*compute.Operation, error)

	insertTargetPool           func(context.Context, steps.GCEConfig, *compute.TargetPool) (*compute.Operation, error)
	insertAddress              func(context.Context, steps.GCEConfig, *compute.Address) (*compute.Operation, error)
	getAddress                 func(context.Context, steps.GCEConfig, string) (*compute.Address, error)
	insertForwardingRule       func(context.Context, steps.GCEConfig, *compute.ForwardingRule) (*compute.Operation, error)
	getForwardingRule       func(context.Context, steps.GCEConfig, string) (*compute.ForwardingRule, error)
	addInstanceToTargetGroup   func(context.Context, steps.GCEConfig, string, *compute.TargetPoolsAddInstanceRequest) (*compute.Operation, error)
	addInstanceToInstanceGroup func(context.Context, steps.GCEConfig, string, *compute.InstanceGroupsAddInstancesRequest) (*compute.Operation, error)
	insertInstanceGroup        func(context.Context, steps.GCEConfig, *compute.InstanceGroup) (*compute.Operation, error)
	insertBackendService       func(context.Context, steps.GCEConfig, *compute.BackendService) (*compute.Operation, error)
	getInstanceGroup           func(context.Context, steps.GCEConfig, string) (*compute.InstanceGroup, error)
	getTargetPool              func(context.Context, steps.GCEConfig, string) (*compute.TargetPool, error)
	getBackendService          func(context.Context, steps.GCEConfig, string) (*compute.BackendService, error)
	deleteForwardingRule       func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	deleteBackendService       func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	deleteInstanceGroup        func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	deleteTargetPool           func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	deleteIpAddress            func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	insertNetwork              func(context.Context, steps.GCEConfig, *compute.Network) (*compute.Operation, error)
	switchNetworkMode          func(context.Context, steps.GCEConfig, string) (*compute.Operation, error)
	getNetwork                 func(context.Context, steps.GCEConfig, string) (*compute.Network, error)

	insertHealthCheck          func(context.Context, steps.GCEConfig, *compute.HealthCheck) (*compute.Operation, error)
	addHealthCheckToTargetPool func(context.Context, steps.GCEConfig, string, *compute.TargetPoolsAddHealthCheckRequest) (*compute.Operation, error)
	getHealthCheck             func(context.Context, steps.GCEConfig, string) (*compute.HealthCheck, error)
}

func Init(getter accountGetter) {
	createInstance := NewCreateInstanceStep(time.Second*10, time.Minute*5)
	createTargetPool := NewCreateTargetPoolStep()
	createIPAddress := NewCreateAddressStep()
	createForwardingRules := NewCreateForwardingRulesStep()

	createInstanceGroup, _ := NewCreateInstanceGroupsStep(getter)
	createBackendService, _ := NewCreateBackendServiceStep()
	createHealthCheck := NewCreateHealthCheckStep()
	createNetworks := NewCreateNetworksStep()

	deleteCluster := NewDeleteClusterStep()
	deleteInstanceGroup, _ := NewDeleteInstanceGroupStep()
	deleteForwardingRules := NewDeleteForwardingRulesStep()
	deleteBackendService, _ := NewDeleteBackendServiceStep()
	deleteTargetPool := NewDeleteTargetPoolStep()
	deleteIpAddress := NewDeleteIpAddressStep()
	deleteNode := NewDeleteNodeStep()

	steps.RegisterStep(CreateHealthCheckStepName, createHealthCheck)
	steps.RegisterStep(DeleteInstanceGroupStepName, deleteInstanceGroup)
	steps.RegisterStep(DeleteBackendServicStepName, deleteBackendService)
	steps.RegisterStep(CreateInstanceGroupsStepName, createInstanceGroup)
	steps.RegisterStep(CreateBackendServiceStepName, createBackendService)
	steps.RegisterStep(CreateInstanceStepName, createInstance)
	steps.RegisterStep(DeleteClusterStepName, deleteCluster)
	steps.RegisterStep(DeleteNodeStepName, deleteNode)
	steps.RegisterStep(CreateTargetPullStepName, createTargetPool)
	steps.RegisterStep(CreateIPAddressStepName, createIPAddress)
	steps.RegisterStep(CreateForwardingRulesStepName, createForwardingRules)
	steps.RegisterStep(DeleteForwardingRulesStepName, deleteForwardingRules)
	steps.RegisterStep(DeleteTargetPoolStepName, deleteTargetPool)
	steps.RegisterStep(DeleteIpAddressStepName, deleteIpAddress)
	steps.RegisterStep(CreateNetworksStepName, createNetworks)
}

func isNotFound(err error) bool {
	if gErr, ok := err.(*googleapi.Error); ok {
		if gErr.Code == http.StatusNotFound {
			return true
		}
	}

	return false
}