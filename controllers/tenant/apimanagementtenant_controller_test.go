package controllers

import (
	"context"
	rhmiconfigv1alpha1 "github.com/integr8ly/integreatly-operator/apis/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("APIManagementTenant controller", func() {
	const (
		TenantName      = "test-tenant-name"
		TenantNamespace = "test-tenant-namespace-dev"

		timeout  = time.Second * 30
		interval = time.Millisecond * 250
	)

	Context("When updating APIManagementTenant Status", func() {
		It("Should reconcile APIManagementTenant after some time", func() {
			By("By creating a new APIManagementTenant CR")
			ctx := context.Background()

			// Create a new APIManagementTenant CR
			tenant := &rhmiconfigv1alpha1.APIManagementTenant{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "integreatly.org/v1alpha1",
					Kind:       "APIManagementTenant",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      TenantName,
					Namespace: TenantNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, tenant)).Should(Succeed())

			// Confirm that the new APIManagementTenant CR was actually created
			tenantLookupKey := types.NamespacedName{Name: TenantName, Namespace: TenantNamespace}
			createdTenant := &rhmiconfigv1alpha1.APIManagementTenant{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, tenantLookupKey, createdTenant)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			// Confirm that the APIManagementTenant CR was created properly
			Expect(createdTenant.Status.ProvisioningStatus).ShouldNot(Equal(""))

			By("By checking that the APIManagementTenant CR has finished reconciling")
			Eventually(func() (rhmiconfigv1alpha1.ProvisioningStatus, error) {
				err := k8sClient.Get(ctx, tenantLookupKey, createdTenant)
				if err != nil {
					return "", err
				}

				return createdTenant.Status.ProvisioningStatus, nil
			}, timeout, interval).Should(ConsistOf(rhmiconfigv1alpha1.ThreeScaleAccountReady))
		})
	})

})
