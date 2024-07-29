package permission

import (
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"go.uber.org/fx"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"google.golang.org/protobuf/proto"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

type Permission struct {
	mu       sync.RWMutex
	roleMaps map[protoreflect.FullName]*pb.RoleOptions
	rbac     *casbin.Enforcer
}

type Params struct {
	fx.In

	RBAC *casbin.Enforcer
}

func New(p Params) *Permission {
	return &Permission{
		mu:       sync.RWMutex{},
		roleMaps: make(map[protoreflect.FullName]*pb.RoleOptions),
		rbac:     p.RBAC,
	}
}

func (r *Permission) Register(protoFile protoreflect.FileDescriptor) {
	r.mu.Lock()
	defer r.mu.Unlock()

	services := protoFile.Services()
	for i := 0; i < services.Len(); i++ {
		methods := services.Get(i).Methods()

		for j := 0; j < methods.Len(); j++ {
			methodDescriptor := methods.Get(j)
			mOpts, ok := methodDescriptor.Options().(*descriptorpb.MethodOptions)
			if !ok {
				continue
			}

			if proto.HasExtension(mOpts, pb.E_Roles) {
				ext, ok2 := proto.GetExtension(mOpts, pb.E_Roles).(*pb.RoleOptions)
				if !ok2 {
					continue
				}

				methodName := methodDescriptor.FullName()
				r.roleMaps[methodName] = ext
			}
		}
	}
}

func (r *Permission) AutoMigrate() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.roleMaps) == 0 {
		return nil
	}

	policies := r.parsePolicies()
	if len(policies) > 0 {
		_, err := r.rbac.AddPolicies(policies)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Permission) parsePolicies() [][]string {
	if len(r.roleMaps) == 0 {
		return nil
	}

	policies := make([][]string, 0)

	for key, val := range r.roleMaps {
		if len(val.Defaults) == 0 {
			continue
		}

		var roles []string
		for _, role := range val.Defaults {
			roles = append(roles, role.String())
		}

		var policy []string
		policy = append(policy, string(key), strings.Join(roles, "|"))
		policies = append(policies, policy)
	}

	return policies
}
