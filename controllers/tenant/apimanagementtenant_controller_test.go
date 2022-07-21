package controllers

import (
	"reflect"
	"testing"

	rhmiv1alpha1 "github.com/integr8ly/integreatly-operator/apis/v1alpha1"
	l "github.com/integr8ly/integreatly-operator/pkg/resources/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGetAPIManagementTenant(t *testing.T) {
	scheme := runtime.NewScheme()
	err := rhmiv1alpha1.SchemeBuilder.AddToScheme(scheme)
	if err != nil {
		t.Errorf("error adding the scheme: %v", err)
	}
	tenant := &rhmiv1alpha1.APIManagementTenant{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIManagementTenant",
			APIVersion: "integreatly.org/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-tenant-name",
			Namespace: "test-tenant-namespace",
		},
	}
	type fields struct {
		Client client.Client
		Scheme *runtime.Scheme
		mgr    controllerruntime.Manager
	}
	type args struct {
		crName      string
		crNamespace string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *rhmiv1alpha1.APIManagementTenant
		wantErr bool
	}{
		{
			name:   "Test - GetAPIManagementTenant() when tenant exists",
			fields: fields{Client: fakeclient.NewFakeClientWithScheme(scheme, tenant)},
			args: args{
				crName:      "test-tenant-name",
				crNamespace: "test-tenant-namespace",
			},
			want:    tenant,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		r := &TenantReconciler{
			Client: tt.fields.Client,
			Scheme: tt.fields.Scheme,
			mgr:    tt.fields.mgr,
			log:    l.Logger{},
		}
		got, err := r.getAPIManagementTenant(tt.args.crName, tt.args.crNamespace)
		if (err != nil) != tt.wantErr {
			t.Errorf("getAPIManagementTenant() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("getAPIManagementTenant() got = %v, want %v", got, tt.want)
		}
	}
}
