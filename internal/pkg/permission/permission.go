package permission

import (
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"google.golang.org/protobuf/proto"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

type Permissions struct {
	mu       sync.RWMutex
	roleMaps map[protoreflect.FullName]*pb.RoleOptions
}

func New() *Permissions {
	return &Permissions{
		mu:       sync.RWMutex{},
		roleMaps: make(map[protoreflect.FullName]*pb.RoleOptions),
	}
}

func (r *Permissions) Register(protoFile protoreflect.FileDescriptor) {
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