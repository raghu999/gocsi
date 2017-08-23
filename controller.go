package gocsi

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/codedellemc/gocsi/csi"
)

const (
	// FMCreateVolume is the full method name for the
	// eponymous RPC message.
	FMCreateVolume = "/" + Namespace +
		".Controller/" +
		"CreateVolume"

	// FMDeleteVolume is the full method name for the
	// eponymous RPC message.
	FMDeleteVolume = "/" + Namespace +
		".Controller/" +
		"DeleteVolume"

	// FMControllerPublishVolume is the full method name for the
	// eponymous RPC message.
	FMControllerPublishVolume = "/" + Namespace +
		".Controller/" +
		"ControllerPublishVolume"

	// FMControllerUnpublishVolume is the full method name for the
	// eponymous RPC message.
	FMControllerUnpublishVolume = "/" + Namespace +
		".Controller/" +
		"ControllerUnpublishVolume"

	// FMValidateVolumeCapabilities is the full method name for the
	// eponymous RPC message.
	FMValidateVolumeCapabilities = "/" + Namespace +
		".Controller/" +
		"ValidateVolumeCapabilities"

	// FMListVolumes is the full method name for the
	// eponymous RPC message.
	FMListVolumes = "/" + Namespace +
		".Controller/" +
		"ListVolumes"

	// FMGetCapacity is the full method name for the
	// eponymous RPC message.
	FMGetCapacity = "/" + Namespace +
		".Controller/" +
		"GetCapacity"

	// FMControllerGetCapabilities is the full method name for the
	// eponymous RPC message.
	FMControllerGetCapabilities = "/" + Namespace +
		".Controller/" +
		"ControllerGetCapabilities"
)

// CreateVolume issues a CreateVolume request to a CSI controller.
func CreateVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	name string,
	requiredBytes, limitBytes uint64,
	fsType string, mountFlags []string,
	params map[string]string,
	callOpts ...grpc.CallOption) (volume *csi.VolumeInfo, err error) {

	req := &csi.CreateVolumeRequest{
		Name:       name,
		Version:    version,
		Parameters: params,
	}

	if requiredBytes > 0 || limitBytes > 0 {
		req.CapacityRange = &csi.CapacityRange{
			LimitBytes:    limitBytes,
			RequiredBytes: requiredBytes,
		}
	}

	if fsType != "" || len(mountFlags) > 0 {
		cap := &csi.VolumeCapability_MountVolume{}
		cap.FsType = fsType
		if len(mountFlags) > 0 {
			cap.MountFlags = mountFlags
		}
		req.VolumeCapabilities = []*csi.VolumeCapability{
			&csi.VolumeCapability{
				Value: &csi.VolumeCapability_Mount{Mount: cap},
			},
		}
	}

	res, err := c.CreateVolume(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	return res.GetResult().VolumeInfo, nil
}

// DeleteVolume issues a DeleteVolume request to a CSI controller.
func DeleteVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeID *csi.VolumeID,
	volumeMetadata *csi.VolumeMetadata,
	callOpts ...grpc.CallOption) error {

	req := &csi.DeleteVolumeRequest{
		Version:        version,
		VolumeId:       volumeID,
		VolumeMetadata: volumeMetadata,
	}

	_, err := c.DeleteVolume(ctx, req, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

// ControllerPublishVolume issues a
// ControllerPublishVolume request
// to a CSI controller.
func ControllerPublishVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeID *csi.VolumeID,
	volumeMetadata *csi.VolumeMetadata,
	nodeID *csi.NodeID,
	readonly bool,
	callOpts ...grpc.CallOption) (
	*csi.PublishVolumeInfo, error) {

	if volumeID == nil {
		return nil, ErrVolumeIDRequired
	}

	req := &csi.ControllerPublishVolumeRequest{
		Version:        version,
		VolumeId:       volumeID,
		VolumeMetadata: volumeMetadata,
		NodeId:         nodeID,
		Readonly:       readonly,
	}

	res, err := c.ControllerPublishVolume(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	return res.GetResult().PublishVolumeInfo, nil
}

// ControllerUnpublishVolume issues a
// ControllerUnpublishVolume request
// to a CSI controller.
func ControllerUnpublishVolume(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeID *csi.VolumeID,
	volumeMetadata *csi.VolumeMetadata,
	nodeID *csi.NodeID,
	callOpts ...grpc.CallOption) error {

	if volumeID == nil {
		return ErrVolumeIDRequired
	}

	req := &csi.ControllerUnpublishVolumeRequest{
		Version:        version,
		VolumeId:       volumeID,
		VolumeMetadata: volumeMetadata,
		NodeId:         nodeID,
	}

	_, err := c.ControllerUnpublishVolume(ctx, req, callOpts...)
	if err != nil {
		return err
	}

	return nil
}

// ValidateVolumeCapabilities issues a ValidateVolumeCapabilities
// request to a CSI controller
func ValidateVolumeCapabilities(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	volumeInfo *csi.VolumeInfo,
	volumeCapabilities []*csi.VolumeCapability,
	callOpts ...grpc.CallOption) (*csi.ValidateVolumeCapabilitiesResponse_Result, error) {

	if volumeInfo == nil {
		return nil, ErrVolumeInfoRequired
	}
	if volumeCapabilities == nil {
		return nil, ErrVolumeCapabilityRequired
	}

	req := &csi.ValidateVolumeCapabilitiesRequest{
		Version:            version,
		VolumeInfo:         volumeInfo,
		VolumeCapabilities: volumeCapabilities,
	}

	res, err := c.ValidateVolumeCapabilities(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	return res.GetResult(), nil
}

// ListVolumes issues a ListVolumes request to a CSI controller.
func ListVolumes(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	maxEntries uint32,
	startingToken string,
	callOpts ...grpc.CallOption) (
	volumes []*csi.VolumeInfo, nextToken string, err error) {

	req := &csi.ListVolumesRequest{
		MaxEntries:    maxEntries,
		StartingToken: startingToken,
		Version:       version,
	}

	res, err := c.ListVolumes(ctx, req, callOpts...)
	if err != nil {
		return nil, "", err
	}

	result := res.GetResult()
	nextToken = result.NextToken
	entries := result.Entries

	// check to see if there are zero entries
	if len(result.Entries) == 0 {
		return nil, nextToken, nil
	}

	volumes = make([]*csi.VolumeInfo, len(entries))

	for x, e := range entries {
		if volumes[x] = e.GetVolumeInfo(); volumes[x] == nil {
			return nil, "", ErrNilVolumeInfo
		}
	}

	return volumes, nextToken, nil
}

// GetCapacity issues a GetCapacity
// request to a CSI controller
func GetCapacity(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	callOpts ...grpc.CallOption) (uint64, error) {

	req := &csi.GetCapacityRequest{
		Version: version,
	}

	res, err := c.GetCapacity(ctx, req, callOpts...)
	if err != nil {
		return 0, err
	}

	return res.GetResult().TotalCapacity, nil
}

// ControllerGetCapabilities issues a ControllerGetCapabilities request to a
// CSI controller.
func ControllerGetCapabilities(
	ctx context.Context,
	c csi.ControllerClient,
	version *csi.Version,
	callOpts ...grpc.CallOption) (
	capabilties []*csi.ControllerServiceCapability, err error) {

	req := &csi.ControllerGetCapabilitiesRequest{
		Version: version,
	}

	res, err := c.ControllerGetCapabilities(ctx, req, callOpts...)
	if err != nil {
		return nil, err
	}

	return res.GetResult().Capabilities, nil
}
