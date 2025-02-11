// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectmatch

import (
	"encoding/json"

	"github.com/goph/emperror"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

type crdMatcher struct {
	objectMatcher ObjectMatcher
}

func NewCRDMatcher(objectMatcher ObjectMatcher) *crdMatcher {
	return &crdMatcher{
		objectMatcher: objectMatcher,
	}
}

// Match compares two extensionsobj.CustomResourceDefinition objects
func (m crdMatcher) Match(oldOrig, newOrig *extensionsobj.CustomResourceDefinition) (bool, error) {
	old := oldOrig.DeepCopy()
	new := newOrig.DeepCopy()

	extensionsobj.SetObjectDefaults_CustomResourceDefinition(new)

	if old.Spec.AdditionalPrinterColumns == nil && newOrig.Spec.AdditionalPrinterColumns == nil {
		new.Spec.AdditionalPrinterColumns = nil
	}
	if old.Spec.Versions == nil && newOrig.Spec.Versions == nil {
		new.Spec.Versions = nil
	}

	type CRD struct {
		ObjectMeta
		Spec extensionsobj.CustomResourceDefinitionSpec
	}

	oldData, err := json.Marshal(CRD{
		ObjectMeta: m.objectMatcher.GetObjectMeta(old.ObjectMeta),
		Spec:       old.Spec,
	})
	if err != nil {
		return false, emperror.WrapWith(err, "could not marshal old object", "name", old.Name)
	}
	newObject := CRD{
		ObjectMeta: m.objectMatcher.GetObjectMeta(new.ObjectMeta),
		Spec:       new.Spec,
	}
	newData, err := json.Marshal(newObject)
	if err != nil {
		return false, emperror.WrapWith(err, "could not marshal new object", "name", new.Name)
	}

	matched, err := m.objectMatcher.MatchJSON(oldData, newData, newObject)
	if err != nil {
		return false, emperror.WrapWith(err, "could not match objects", "name", new.Name)
	}

	return matched, nil
}
