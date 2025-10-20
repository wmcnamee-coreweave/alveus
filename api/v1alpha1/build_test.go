package v1alpha1

import (
	"fmt"

	"github.com/cakehappens/gocto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service.Inflate()", func() {
	var (
		service Service
	)

	BeforeEach(func() {
		service = Service{}
	})

	JustBeforeEach(func() {
		service.Inflate()
	})

	It("should set github.on.dispatch", func() {
		Expect(service.Github.On.Dispatch).To(Equal(&gocto.OnDispatch{}))
	})

	Context("destination.namespace", func() {
		BeforeEach(func() {
			service.DestinationGroups = DestinationGroups{
				{
					Destinations: Destinations{
						{},
					},
				},
			}
		})

		type TableEntry struct {
			serviceLevel string
			groupLevel   string
			destLevel    string

			expected string
		}

		for _, entry := range []TableEntry{
			{"top", "", "", "top"},
			{"", "group", "", "group"},
			{"top", "group", "", "group"},
			{"", "", "dest", "dest"},
			{"", "group", "dest", "dest"},
			{"top", "group", "dest", "dest"},
		} {
			Context("entry", func() {
				BeforeEach(func() {
					service.DestinationNamespace = entry.serviceLevel
					service.DestinationGroups[0].DestinationNamespace = entry.groupLevel
					service.DestinationGroups[0].Destinations[0].Namespace = entry.destLevel
				})

				It(fmt.Sprintf("should set namespace to %s", entry.expected), func() {
					for _, group := range service.DestinationGroups {
						for _, destination := range group.Destinations {
							Expect(destination.Namespace).To(Equal(entry.expected))
						}
					}
				})
			})
		}
	})

	Context("argocd.hostname", func() {
		BeforeEach(func() {
			service.DestinationGroups = DestinationGroups{
				{
					Destinations: Destinations{
						{},
					},
				},
			}
		})

		type TableEntry struct {
			serviceLevel string
			groupLevel   string
			destLevel    string

			expected string
		}

		for _, entry := range []TableEntry{
			{"top", "", "", "top"},
			{"", "group", "", "group"},
			{"top", "group", "", "group"},
			{"", "", "dest", "dest"},
			{"", "group", "dest", "dest"},
			{"top", "group", "dest", "dest"},
		} {
			Context("entry", func() {
				BeforeEach(func() {
					service.ArgoCD.Hostname = entry.serviceLevel
					service.DestinationGroups[0].ArgoCD.Hostname = entry.groupLevel
					service.DestinationGroups[0].Destinations[0].ArgoCD.Hostname = entry.destLevel
				})

				It(fmt.Sprintf("should set namespace to %s", entry.expected), func() {
					for _, group := range service.DestinationGroups {
						for _, destination := range group.Destinations {
							Expect(destination.ArgoCD.Hostname).To(Equal(entry.expected))
						}
					}
				})
			})
		}
	})
})
