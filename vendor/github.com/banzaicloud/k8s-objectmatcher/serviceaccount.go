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
	corev1 "k8s.io/api/core/v1"
)

type serviceAccountMatcher struct {
	objectMatcher ObjectMatcher
}

func NewServiceAccountMatcher(objectMatcher ObjectMatcher) *serviceAccountMatcher {
	return &serviceAccountMatcher{
		objectMatcher: objectMatcher,
	}
}

// Match compares two corev1.ServiceAccount objects
func (m serviceAccountMatcher) Match(oldOrig, newOrig *corev1.ServiceAccount) (bool, error) {
	old := oldOrig.DeepCopy()
	new := newOrig.DeepCopy()

	type ServiceAccount struct {
		ObjectMeta
	}

	oldData, err := json.Marshal(ServiceAccount{
		ObjectMeta: m.objectMatcher.GetObjectMeta(old.ObjectMeta),
	})
	if err != nil {
		return false, emperror.WrapWith(err, "could not marshal old object", "name", old.Name)
	}
	newObject := ServiceAccount{
		ObjectMeta: m.objectMatcher.GetObjectMeta(new.ObjectMeta),
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
