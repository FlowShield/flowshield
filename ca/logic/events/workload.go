package events

func NewWorkloadLifeCycle(op string, author string, cert CertOp) *Op {
	return &Op{
		Operator: author,
		Category: CategoryWorkloadLifecycle,
		Type:     op,
		Obj:      cert,
	}
}
