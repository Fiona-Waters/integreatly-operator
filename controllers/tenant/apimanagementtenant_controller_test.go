package controllers

import (
	"context"
	"time"

	rhmiconfigv1alpha1 "github.com/integr8ly/integreatly-operator/apis/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	usersv1 "github.com/openshift/api/user/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("APIManagementTenant controller", func() {
	const (
		UserName        = "test-user"
		TenantName      = "example"
		TenantNamespace = "test-user-dev"

		timeout  = time.Second * 30
		interval = time.Millisecond * 250
	)

	Context("APIManagementTenant CRs", func() {
		It("Created APIManagementTenants should reconcile after some time", func() {
			By("Creating a new APIManagementTenant CR")
			ctx := context.Background()

			// create test namespace
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: TenantNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, ns)).Should(Succeed())

			// create test user
			user := &usersv1.User{
				TypeMeta: metav1.TypeMeta{
					APIVersion: usersv1.GroupVersion.String(),
					Kind:       "User",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: UserName,
				},
			}

			Expect(k8sClient.Create(ctx, user)).Should(Succeed())

			// Create a new APIManagementTenant CR
			tenant := &rhmiconfigv1alpha1.APIManagementTenant{
				TypeMeta: metav1.TypeMeta{
					APIVersion: rhmiconfigv1alpha1.GroupVersion.String(),
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

			By("Checking that the APIManagementTenant CR has finished reconciling")
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
