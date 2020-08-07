// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigPatch) DeepCopyInto(out *ConfigPatch) {
	*out = *in
	if in.WorkloadSelectorLabels != nil {
		in, out := &in.WorkloadSelectorLabels, &out.WorkloadSelectorLabels
		*out = new(map[string]string)
		if **in != nil {
			in, out := *in, *out
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
		}
	}
	in.RateLimitProperty.DeepCopyInto(&out.RateLimitProperty)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigPatch.
func (in *ConfigPatch) DeepCopy() *ConfigPatch {
	if in == nil {
		return nil
	}
	out := new(ConfigPatch)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Descriptor) DeepCopyInto(out *Descriptor) {
	*out = *in
	out.RateLimit = in.RateLimit
	if in.Descriptors != nil {
		in, out := &in.Descriptors, &out.Descriptors
		*out = make([]DescriptorInternal, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Descriptor.
func (in *Descriptor) DeepCopy() *Descriptor {
	if in == nil {
		return nil
	}
	out := new(Descriptor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DescriptorInternal) DeepCopyInto(out *DescriptorInternal) {
	*out = *in
	out.RateLimit = in.RateLimit
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DescriptorInternal.
func (in *DescriptorInternal) DeepCopy() *DescriptorInternal {
	if in == nil {
		return nil
	}
	out := new(DescriptorInternal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimit) DeepCopyInto(out *RateLimit) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimit.
func (in *RateLimit) DeepCopy() *RateLimit {
	if in == nil {
		return nil
	}
	out := new(RateLimit)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimitProperty) DeepCopyInto(out *RateLimitProperty) {
	*out = *in
	if in.Descriptors != nil {
		in, out := &in.Descriptors, &out.Descriptors
		*out = make([]Descriptor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimitProperty.
func (in *RateLimitProperty) DeepCopy() *RateLimitProperty {
	if in == nil {
		return nil
	}
	out := new(RateLimitProperty)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiter) DeepCopyInto(out *RateLimiter) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiter.
func (in *RateLimiter) DeepCopy() *RateLimiter {
	if in == nil {
		return nil
	}
	out := new(RateLimiter)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimiter) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterConfig) DeepCopyInto(out *RateLimiterConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterConfig.
func (in *RateLimiterConfig) DeepCopy() *RateLimiterConfig {
	if in == nil {
		return nil
	}
	out := new(RateLimiterConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimiterConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterConfigList) DeepCopyInto(out *RateLimiterConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RateLimiterConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterConfigList.
func (in *RateLimiterConfigList) DeepCopy() *RateLimiterConfigList {
	if in == nil {
		return nil
	}
	out := new(RateLimiterConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimiterConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterConfigSpec) DeepCopyInto(out *RateLimiterConfigSpec) {
	*out = *in
	if in.ConfigPatches != nil {
		in, out := &in.ConfigPatches, &out.ConfigPatches
		*out = make([]ConfigPatch, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterConfigSpec.
func (in *RateLimiterConfigSpec) DeepCopy() *RateLimiterConfigSpec {
	if in == nil {
		return nil
	}
	out := new(RateLimiterConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterConfigStatus) DeepCopyInto(out *RateLimiterConfigStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterConfigStatus.
func (in *RateLimiterConfigStatus) DeepCopy() *RateLimiterConfigStatus {
	if in == nil {
		return nil
	}
	out := new(RateLimiterConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterList) DeepCopyInto(out *RateLimiterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RateLimiter, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterList.
func (in *RateLimiterList) DeepCopy() *RateLimiterList {
	if in == nil {
		return nil
	}
	out := new(RateLimiterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RateLimiterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterSpec) DeepCopyInto(out *RateLimiterSpec) {
	*out = *in
	if in.Port != nil {
		in, out := &in.Port, &out.Port
		*out = new(int32)
		**out = **in
	}
	if in.LogLevel != nil {
		in, out := &in.LogLevel, &out.LogLevel
		*out = new(LogLevel)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterSpec.
func (in *RateLimiterSpec) DeepCopy() *RateLimiterSpec {
	if in == nil {
		return nil
	}
	out := new(RateLimiterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RateLimiterStatus) DeepCopyInto(out *RateLimiterStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RateLimiterStatus.
func (in *RateLimiterStatus) DeepCopy() *RateLimiterStatus {
	if in == nil {
		return nil
	}
	out := new(RateLimiterStatus)
	in.DeepCopyInto(out)
	return out
}
