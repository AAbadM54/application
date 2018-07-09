package memcached_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/kubernetes-sigs/kubebuilder/test/projects/memcached-api-server/pkg/apis/myapps/v1alpha1"
	. "github.com/kubernetes-sigs/kubebuilder/test/projects/memcached-api-server/pkg/client/clientset/versioned/typed/myapps/v1alpha1"
)

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to implement controller logic tests

var _ = Describe("Memcached controller", func() {
	var instance Memcached
	var expectedKey types.ReconcileKey
	var client MemcachedInterface

	BeforeEach(func() {
		instance = Memcached{}
		instance.Name = "instance-1"
		expectedKey = types.ReconcileKey{
			Namespace: "default",
			Name:      "instance-1",
		}
	})

	AfterEach(func() {
		client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	Describe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			after := make(chan struct{})
			ctrl.AfterReconcile = func(key types.ReconcileKey, err error) {
				defer func() {
					// Recover in case the key is reconciled multiple times
					defer func() { recover() }()
					close(after)
				}()
				defer GinkgoRecover()
				Expect(key).To(Equal(expectedKey))
				Expect(err).ToNot(HaveOccurred())
			}

			// Create the instance
			client = cs.MyappsV1alpha1().Memcacheds("default")

			instance.Spec.Size = 1
			_, err := client.Create(&instance)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(MatchRegexp(".*spec.size in body should be greater than or equal to 5.*"))

			instance.Spec.Size = 101
			_, err = client.Create(&instance)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(MatchRegexp(".*spec.size in body should be less than or equal to 100.*"))

			instance.Spec.Size = 50
			_, err = client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for reconcile to happen
			Eventually(after, "10s", "100ms").Should(BeClosed())

			// INSERT YOUR CODE HERE - test conditions post reconcile
		})
	})
})
