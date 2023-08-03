package service_type

type Type int

const (
	Region Type = iota
	Ecs
	Ces
	Vpc
	Ims
	Evs
	Iam
	Obs
	Sfs
	Rds
	Eip
)

func (c Type) String() string {
	return [...]string{"region", "ecs", "ces", "vpc", "ims", "evs", "iam", "obs", "sfs", "rds", "eip"}[c]
}
