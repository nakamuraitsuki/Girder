package entity

import (
	"fmt"
	"strings"
	"girder/internal/domain/value"
)

type NodeState string

const (
	NodeStateProvisioning NodeState = "provisioning"
	NodeStateReady        NodeState = "ready"
	NodeStateDraining     NodeState = "draining"
	NodeStateOffline      NodeState = "offline"
)

// Node は Girder クラスタ内の物理Nodeを表すエンティティ。ノードは仮想マシンのホストとして機能し、CPU、メモリ、ディスク容量などのリソースを提供する。
type Node struct {
	ID          NodeID
	Name        string
	Description string
	State       NodeState
	Capacity    NodeCapacity
	Allocated   NodeAllocation
	Labels      map[string]string
}

// Nodeのコンストラクタ。
func NewNode(name string, capacity NodeCapacity, labels map[string]string) (Node, error) {
	// 空文字、負のリソース値はエラーとして扱う。
	if strings.TrimSpace(name) == "" {
		return Node{}, fmt.Errorf("%w: node name", ErrEmptyName)
	}
	if capacity.CPU < 0 || capacity.MemoryMiB < 0 || capacity.DiskGiB < 0 {
		return Node{}, fmt.Errorf("%w: node capacity", ErrInvalidResource)
	}

	return Node{
		ID:       NewNodeID(),
		Name:     strings.TrimSpace(name),
		State:    NodeStateProvisioning,
		Capacity: capacity,
		Labels:   value.CloneStringMap(labels),
	}, nil
}

// Name Setter
func (n *Node) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w: node name", ErrEmptyName)
	}
	n.Name = strings.TrimSpace(name)
	return nil
}

// Description Setter
func (n *Node) SetDescription(description string) {
	n.Description = strings.TrimSpace(description)
}

// State Setter
func (n *Node) SetState(state NodeState) {
	n.State = state
}

// Label Setters
func (n *Node) SetLabel(key, val string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("%w: label key", ErrEmptyName)
	}
	if n.Labels == nil {
		n.Labels = map[string]string{}
	}
	n.Labels[key] = strings.TrimSpace(val)
	return nil
}
func (n *Node) RemoveLabel(key string) {
	if n.Labels == nil {
		return
	}
	delete(n.Labels, strings.TrimSpace(key))
}

// Capacity and Allocation Management
func (n Node) RemainingCapacity() NodeCapacity {
	return NodeCapacity{
		CPU:       n.Capacity.CPU - n.Allocated.CPU,
		MemoryMiB: n.Capacity.MemoryMiB - n.Allocated.MemoryMiB,
		DiskGiB:   n.Capacity.DiskGiB - n.Allocated.DiskGiB,
	}
}

// CanHost は指定されたリソース要求をホストできるかどうかを判定する。
func (n Node) CanHost(request NodeAllocation) bool {
	return n.RemainingCapacity().CanHost(request)
}
// Reserve は指定されたリソース要求をノードに割り当てる。要求がノードの残り容量を超える場合、エラーを返す。
func (n *Node) Reserve(request NodeAllocation) error {
	if !n.CanHost(request) {
		return fmt.Errorf("%w: node %s", ErrInsufficientCapacity, n.ID)
	}
	n.Allocated.CPU += request.CPU
	n.Allocated.MemoryMiB += request.MemoryMiB
	n.Allocated.DiskGiB += request.DiskGiB
	return nil
}

// Release は指定されたリソース要求をノードから解放する。要求がノードの割り当てを超える場合、エラーを返す。
func (n *Node) Release(request NodeAllocation) error {
	if request.CPU < 0 || request.MemoryMiB < 0 || request.DiskGiB < 0 {
		return fmt.Errorf("%w: node release", ErrInvalidResource)
	}
	if request.CPU > n.Allocated.CPU || request.MemoryMiB > n.Allocated.MemoryMiB || request.DiskGiB > n.Allocated.DiskGiB {
		return fmt.Errorf("%w: node release", ErrInvalidState)
	}
	n.Allocated.CPU -= request.CPU
	n.Allocated.MemoryMiB -= request.MemoryMiB
	n.Allocated.DiskGiB -= request.DiskGiB
	return nil
}