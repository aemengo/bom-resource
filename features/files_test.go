package features_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/bom-resource/features"
)

var _ = Describe("feature files", func() {
	Context("Generates a feature file", func() {
		expectedKeys := []string{"foo", "hello"}
		It("given a default file produces a map with same content", func() {
			output, err := features.Generate(true, []string{"./fixtures/default.yml"}, expectedKeys)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).ToNot(BeNil())
			Expect(output).To(HaveKeyWithValue("foo", "bar"))
			Expect(output).To(HaveKeyWithValue("hello", "world"))
		})

		It("given a default file and single override produces a map with overridden content", func() {
			output, err := features.Generate(true, []string{"./fixtures/default.yml", "./fixtures/prod.yml"}, expectedKeys)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).ToNot(BeNil())
			Expect(output).To(HaveKeyWithValue("foo", "foo-1"))
			Expect(output).To(HaveKeyWithValue("hello", "world"))
		})

		It("given a default file and multiple override produces a map with overridden content", func() {
			output, err := features.Generate(true, []string{"./fixtures/default.yml", "./fixtures/us.yml", "./fixtures/prod.yml"}, expectedKeys)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).ToNot(BeNil())
			Expect(output).To(HaveKeyWithValue("foo", "foo-1"))
			Expect(output).To(HaveKeyWithValue("hello", "america"))
		})

		It("given a default file and empty override produces a map with same content", func() {
			output, err := features.Generate(true, []string{"./fixtures/default.yml", "./fixtures/empty.yml"}, expectedKeys)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).ToNot(BeNil())
			Expect(output).To(HaveKeyWithValue("foo", "bar"))
			Expect(output).To(HaveKeyWithValue("hello", "world"))
		})

		It("given a default file and missing override produces a map with same content", func() {
			output, err := features.Generate(true, []string{"./fixtures/default.yml", "./fixtures/missing.yml"}, expectedKeys)
			Expect(err).ToNot(HaveOccurred())
			Expect(output).ToNot(BeNil())
			Expect(output).To(HaveKeyWithValue("foo", "bar"))
			Expect(output).To(HaveKeyWithValue("hello", "world"))
		})

		It("errors if all expected keys are not provided", func() {
			_, err := features.Generate(true, []string{"./fixtures/default.yml"}, []string{"foo", "hello", "new_key", "new_key2"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Missing keys"))
		})

		It("skip key validation", func() {
			_, err := features.Generate(false, []string{"./fixtures/default.yml"}, []string{"foo", "hello", "new_key", "new_key2"})
			Expect(err).To(Not(HaveOccurred()))
		})

		It("errors keys other than expected are provided", func() {
			_, err := features.Generate(true, []string{"./fixtures/default.yml"}, []string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("were provided but not expected"))
		})
	})
})
