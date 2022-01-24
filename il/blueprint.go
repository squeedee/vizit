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

func getRef(resource blueprint.Resource, option *blueprint.Option) blueprint.ClusterResourceRef {
	if resource.TemplateRef != nil {
		return blueprint.ClusterResourceRef{
			Kind: resource.TemplateRef.Kind,
			Name: resource.TemplateRef.Name,
		}
	}

	return blueprint.ClusterResourceRef{
		Kind: resource.Kind,
		Name: option.Name,
	}
}

func getSelector(resource blueprint.Resource, option *blueprint.Option) *metav1.LabelSelector {
	if resource.Selector != nil {
		return resource.Selector.DeepCopy()
	}

	return option.Selector.DeepCopy()
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
	options := map[string]Option{}

	for _, inResource := range in.Spec.Resources {
		if len(inResource.Options) > 0 {
			for _, inOption := range inResource.Options {
				options[inResource.Name+":"+inOption.Name] = Option{
					ResourceName: inResource.Name,
					TemplateRef: getRef(inResource, &inOption),
					Criteria: Criteria{
						Selector: getSelector(inResource, &inOption),
						Inputs:   getInputs(inResource, &inOption),
					},
				}
			}
		} else {
			ref := getRef(inResource, nil)
			options[inResource.Name+":"+ref.Name] = Option{
				ResourceName: inResource.Name,
				TemplateRef: ref,
				Criteria: Criteria{
					Selector: getSelector(inResource, nil),
					Inputs:   getInputs(inResource, nil),
				},
			}

		}
	}

	return Blueprint{
		Options: options,
	}
}
