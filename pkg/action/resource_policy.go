/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package action

import (
	"strings"

	"helm.sh/helm/v3/pkg/releaseutil"
)

// resourcePolicyAnno is the annotation name for a resource policy
const resourcePolicyAnno = "helm.sh/resource-policy"

// keepPolicy is the resource policy type for keep
//
// This resource policy type allows resources to skip being deleted
// during an uninstallRelease action.
const keepPolicy = "keep"

// preferExistingPolicy is the resource policy tpye for prefer-existing.
//
// This resource policy type indicates to Helm that this resource should only be created if it does not already exist.
const preferExistingPolicy = "prefer-existing"

func filterManifestsToKeep(manifests []releaseutil.Manifest) (keep, remaining []releaseutil.Manifest) {
	return filterManifestsByPolicyType(manifests, keepPolicy)
}

func filterManifestsToPreferExisting(manifests []releaseutil.Manifest) (preferExisting, remaining []releaseutil.Manifest) {
	return filterManifestsByPolicyType(manifests, preferExistingPolicy)
}

func filterManifestsByPolicyType(manifests []releaseutil.Manifest, filteredPolicyType string) (included, excluded []releaseutil.Manifest) {
	for _, m := range manifests {
		if m.Head.Metadata == nil || m.Head.Metadata.Annotations == nil || len(m.Head.Metadata.Annotations) == 0 {
			excluded = append(excluded, m)
			continue
		}

		resourcePolicyType, ok := m.Head.Metadata.Annotations[resourcePolicyAnno]
		if !ok {
			excluded = append(excluded, m)
			continue
		}

		resourcePolicyType = strings.ToLower(strings.TrimSpace(resourcePolicyType))
		if resourcePolicyType == filteredPolicyType {
			included = append(included, m)
		}

	}
	return included, excluded
}

