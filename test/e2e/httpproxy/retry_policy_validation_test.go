// Copyright Project Contour Authors
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

//go:build e2e
// +build e2e

package httpproxy

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	contourv1 "github.com/projectcontour/contour/apis/projectcontour/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func testRetryPolicyValidation(namespace string) {
	Specify("retry policy is validated on create", func() {
		t := f.T()

		p := &contourv1.HTTPProxy{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      "invalid-retry-on-condition",
			},
			Spec: contourv1.HTTPProxySpec{
				Routes: []contourv1.Route{
					{
						Services: []contourv1.Service{
							{
								Name: "foo",
								Port: 80,
							},
						},
						RetryPolicy: &contourv1.RetryPolicy{
							RetryOn: []contourv1.RetryOn{
								"foobar",
							},
						},
					},
				},
			},
		}
		err := f.Client.Create(context.TODO(), p)
		require.Error(t, err)

		// Kubernetes 1.24 adds array indexes to the error message, so allow
		// either format for now.
		isExpectedErr := func(err error) bool {
			return strings.Contains(err.Error(), "spec.routes.retryPolicy.retryOn: Unsupported value") ||
				strings.Contains(err.Error(), "spec.routes[0].retryPolicy.retryOn[0]: Unsupported value")
		}
		assert.True(t, isExpectedErr(err))
	})
}
