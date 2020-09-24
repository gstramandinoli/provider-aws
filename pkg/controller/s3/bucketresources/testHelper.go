/*
Copyright 2020 The Crossplane Authors.

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

package bucketresources

import (
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	corev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"

	"github.com/crossplane/provider-aws/apis/s3/v1beta1"
)

var (
	enabled               = "Enabled"
	suspended             = "Suspended"
	errBoom               = errors.New("boom")
	accelGetFailed        = "cannot get bucket accelerate configuration"
	accelPutFailed        = "cannot put bucket acceleration configuration"
	accelDeleteFailed     = "cannot delete bucket acceleration configuration"
	corsGetFailed         = "cannot get bucket CORS configuration"
	corsPutFailed         = "cannot put bucket cors"
	corsDeleteFailed      = "cannot delete bucket CORS configuration"
	lifecycleGetFailed    = "cannot get bucket lifecycle"
	lifecyclePutFailed    = "cannot put bucket lifecycle"
	lifecycleDeleteFailed = "cannot delete bucket lifecycle configuration"
)

type bucketModifier func(policy *v1beta1.Bucket)

func withConditions(c ...corev1alpha1.Condition) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Status.ConditionedStatus.Conditions = c }
}

func withAccelerationConfig(s *v1beta1.AccelerateConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.AccelerateConfiguration = s }
}

func withSSEConfig(s *v1beta1.ServerSideEncryptionConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.ServerSideEncryptionConfiguration = s }
}

func withVersioningConfig(s *v1beta1.VersioningConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.VersioningConfiguration = s }
}

func withCORSConfig(s *v1beta1.CORSConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.CORSConfiguration = s }
}

func withWebConfig(s *v1beta1.WebsiteConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.WebsiteConfiguration = s }
}

func withLoggingConfig(s *v1beta1.LoggingConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.LoggingConfiguration = s }
}

func withPayerConfig(s *v1beta1.PaymentConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.PayerConfiguration = s }
}

func withTaggingConfig(s *v1beta1.Tagging) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.BucketTagging = s }
}

func withReplConfig(s *v1beta1.ReplicationConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.ReplicationConfiguration = s }
}

func withLifecycleConfig(s *v1beta1.BucketLifecycleConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.LifecycleConfiguration = s }
}

func withNotificationConfig(s *v1beta1.NotificationConfiguration) bucketModifier { //nolint
	return func(r *v1beta1.Bucket) { r.Spec.ForProvider.NotificationConfiguration = s }
}

func bucket(m ...bucketModifier) *v1beta1.Bucket {
	cr := &v1beta1.Bucket{
		Spec: v1beta1.BucketSpec{
			ForProvider: v1beta1.BucketParameters{
				ACL:                               aws.String("private"),
				LocationConstraint:                aws.String("us-east-1"),
				GrantFullControl:                  nil,
				GrantRead:                         nil,
				GrantReadACP:                      nil,
				GrantWrite:                        nil,
				GrantWriteACP:                     nil,
				ObjectLockEnabledForBucket:        nil,
				ServerSideEncryptionConfiguration: nil,
				VersioningConfiguration:           nil,
				AccelerateConfiguration:           nil,
				CORSConfiguration:                 nil,
				WebsiteConfiguration:              nil,
				LoggingConfiguration:              nil,
				PayerConfiguration:                nil,
				BucketTagging:                     nil,
				ReplicationConfiguration:          nil,
				LifecycleConfiguration:            nil,
				NotificationConfiguration:         nil,
			},
		},
	}
	for _, f := range m {
		f(cr)
	}
	return cr
}

func createRequest(err error, data interface{}) *aws.Request {
	return &aws.Request{HTTPRequest: &http.Request{}, Retryer: aws.NoOpRetryer{}, Error: err, Data: data}
}

func copyTag(tag *v1beta1.Tag) *awss3.Tag {
	if tag == nil {
		return nil
	}
	return &awss3.Tag{
		Key:   aws.String(tag.Key),
		Value: aws.String(tag.Value),
	}
}

func copyTags(tags []v1beta1.Tag) []awss3.Tag {
	if tags == nil {
		return nil
	}
	out := make([]awss3.Tag, len(tags))
	for i := range tags {
		out[i] = *copyTag(&tags[i])
	}
	return out
}
