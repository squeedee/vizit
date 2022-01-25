package il

import (
	"github.com/bunniesandbeatings/vizit/blueprint"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Stage string

type Criteria struct {
	Selector *metav1.LabelSelector
	Inputs   []string
}
type Option struct {
	ResourceName string
	TemplateRef blueprint.ClusterResourceRef

	Criteria
}

type Resource struct {
	Name    string
	Options []Option
}

type Blueprint struct {
	Options map[string]Option
	Resources map[string]Resource
}

func (bp *Blueprint) Entrypoints() map[string]Option {
	ep := map[string]Option{}

	for name, option := range bp.Options {
		if len(option.Criteria.Inputs) < 1 {
			ep[name] = option
		}
	}

	return ep
}

func getRef(templateRef *blueprint.ClusterResourceRef, option *blueprint.Option) blueprint.ClusterResourceRef {
	if option != nil {
		return blueprint.ClusterResourceRef{
			Kind: templateRef.Kind,
			Name: option.Name,
		}
	}

	return blueprint.ClusterResourceRef{
		Kind: templateRef.Kind,
		Name: templateRef.Name,
	}

}

func getSelector(resource blueprint.Resource, option *blueprint.Option) *metav1.LabelSelector {
	if resource.Selector != nil {
		return resource.Selector.DeepCopy()
	}

	if option != nil && option.Selector != nil {
		return option.Selector.DeepCopy()
	}

	return nil
}

func getInputs(resource blueprint.Resource, option *blueprint.Option) []string {
	var inputs []string

	for _, config := range resource.Configs {
		inputs = append(inputs, config.Resource)
	}
	for _, image := range resource.Images {
		inputs = append(inputs, image.Resource)
	}
	for _, source := range resource.Sources {
		inputs = append(inputs, source.Resource)
	}
	if option != nil {
		for _, config := range option.Configs {
			inputs = append(inputs, config.Resource)
		}
		for _, image := range option.Images {
			inputs = append(inputs, image.Resource)
		}
		for _, source := range option.Sources {
			inputs = append(inputs, source.Resource)
		}
	}

	return inputs
}

func ParseBlueprint(in blueprint.Blueprint) Blueprint {
	resources := map[string]Resource{}
	options := map[string]Option{}

	for _, inResource := range in.Spec.Resources {
		resource := Resource{ Name: inResource.Name}
		if len(inResource.TemplateRef.Options) > 0 {
			for _, inOption := range inResource.TemplateRef.Options {
				opt := Option{
					ResourceName: inResource.Name,
					TemplateRef: getRef(inResource.TemplateRef, &inOption),
					Criteria: Criteria{
						Selector: getSelector(inResource, &inOption),
						Inputs:   getInputs(inResource, &inOption),
					},
				}
				options[inResource.Name+":"+inOption.Name] = opt
				resource.Options = append(resource.Options, opt)
			}

		} else {
			ref := getRef(inResource.TemplateRef, nil)
			opt := Option{
				ResourceName: inResource.Name,
				TemplateRef: ref,
				Criteria: Criteria{
					Selector: getSelector(inResource, nil),
					Inputs:   getInputs(inResource, nil),
				},
			}
			options[inResource.Name+":"+ref.Name] = opt
			resource.Options = append(resource.Options, opt)
		}
		resources[resource.Name] =  resource
	}

	return Blueprint{
		Options: options,
		Resources: resources,
	}
}
