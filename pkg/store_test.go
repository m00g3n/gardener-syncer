package seeker_test

import (
	"context"
	"fmt"
	"testing"

	seeker "github.com/kyma-project/gardener-syncer/pkg"
	"github.com/kyma-project/gardener-syncer/pkg/types"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	errPatchFailedTest = fmt.Errorf("patch failed test")
	errGetFailedTest   = fmt.Errorf("get failed test")
	errTestFailed      = fmt.Errorf("test failed")

	testName      = "test-name"
	testNamespace = "test-namespace"
	testData      = map[string]string{"test": `seedRegions:
- me
- plz`}
	testProviderRegions = types.Providers{
		"test": {
			SeedRegions: []string{
				"me", "plz",
			},
		},
	}

	testCM = corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: testNamespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "configmap",
			APIVersion: "v1",
		},
		Data: map[string]string{
			"test": `seedRegions:
- me
- plz`},
	}
)

func buildGetWithError(err error) seeker.Get {
	return func(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error {
		return err
	}
}

func buildGetNotFound(group, resource, name string) seeker.Get {
	return func(context.Context, client.ObjectKey, client.Object, ...client.GetOption) error {
		return errors.NewNotFound(schema.GroupResource{
			Group:    group,
			Resource: resource,
		}, name)
	}
}

func buildGet(out corev1.ConfigMap) seeker.Get {
	return func(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
		if key.Name != out.GetName() || key.Namespace != out.GetNamespace() {
			return fmt.Errorf("invalid test case")
		}

		obj = &out
		return nil
	}
}

func buildPatch(expectedName, expectedNamespace string, expectedData map[string]string) seeker.Patch {
	return func(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
		cm, ok := obj.(*corev1.ConfigMap)
		if !ok {
			return fmt.Errorf("%w: invalid type, expected *v1.ConfigMap", errTestFailed)
		}

		if expectedName != cm.Name {
			return fmt.Errorf("%w: invalid name, expected: '%s', actual: '%s'", errTestFailed, expectedName, cm.Name)
		}

		if expectedNamespace != cm.Namespace {
			return fmt.Errorf("%w: invalid namespace, expected: '%s', actual: '%s'", errTestFailed, expectedNamespace, cm.Namespace)
		}

		for k, v := range expectedData {
			actualValue, found := cm.Data[k]
			if !found {
				return fmt.Errorf("%w: key not found: '%s'", errTestFailed, k)
			}
			if cm.Data[k] != v {
				return fmt.Errorf("%w: invalid value, expected: '%s', actual: '%s'", errTestFailed, v, actualValue)
			}
		}

		return nil
	}
}

type stringMap = map[string]string

func TestBuildStoreFn(t *testing.T) {
	testCases := []struct {
		title       string
		key         client.ObjectKey
		get         seeker.Get
		patch       seeker.Patch
		data2Store  types.Providers
		expectedErr error
	}{
		{
			title: "GET:random fail",
			key: client.ObjectKey{
				Name:      testName,
				Namespace: testNamespace,
			},
			get:         buildGetWithError(errGetFailedTest),
			expectedErr: errGetFailedTest,
		},
		{
			title: "GET:not found",
			key: client.ObjectKey{
				Name:      testName,
				Namespace: testNamespace,
			},
			data2Store: testProviderRegions,
			get:        buildGetNotFound("", "configmap", testName),
			patch:      buildPatch(testName, testNamespace, testData),
		},
		{
			title: "OK",
			key: client.ObjectKey{
				Name:      testName,
				Namespace: testNamespace,
			},
			data2Store: testProviderRegions,
			get:        buildGet(testCM),
			patch:      buildPatch(testName, testNamespace, testData),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.title, func(t *testing.T) {
			// GIVEN
			store := seeker.BuildStoreFn(seeker.StoreOpts{
				Key:     testCase.key,
				Patch:   testCase.patch,
				Get:     testCase.get,
				Convert: seeker.ToConfigMap,
			})

			// WHEN
			err := store(testCase.data2Store)

			// THEN
			if testCase.expectedErr == nil {
				require.NoError(t, err)
			}

			// THEN
			if testCase.expectedErr != nil {
				require.EqualError(t, err, testCase.expectedErr.Error())
			}
		})
	}
}
