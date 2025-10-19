package v1alpha1

import (
	"errors"
	"strconv"
	"strings"

	argov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service.Validate()", func() {
	var (
		service *Service

		actualErr error
	)

	BeforeEach(func() {
		service = &Service{}
		actualErr = nil
	})

	JustBeforeEach(func() {
		actualErr = service.Validate()
	})

	Context("a valid service", func() {
		BeforeEach(func() {
			service.Name = "foo"
			service.sourceValidatorFunc = func(source Source) error {
				return nil
			}
			service.destinationGroupsValidatorFunc = func(groups DestinationGroups) error {
				return nil
			}
		})

		It("should return no error", func() {
			Expect(actualErr).NotTo(HaveOccurred())
		})

		Context("except", func() {
			When("the name is empty", func() {
				BeforeEach(func() {
					service.Name = ""
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("service name is required"))
				})
			})

			When("the source fails to validate", func() {
				BeforeEach(func() {
					service.sourceValidatorFunc = func(source Source) error {
						return errors.New("source validation error")
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating source: source validation error"))
				})
			})

			When("the destinationGroups fail to validate", func() {
				BeforeEach(func() {
					service.destinationGroupsValidatorFunc = func(groups DestinationGroups) error {
						return errors.New("destination groups validation error")
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("destination groups validation error"))
				})
			})
		})

	})

	When("the service is nil", func() {
		BeforeEach(func() {
			service = nil
		})

		It("should return an error", func() {
			Expect(actualErr).To(MatchError("service is nil"))
		})
	})

	When("the service name is empty", func() {
		BeforeEach(func() {
			service.Name = ""
		})
	})
})

var _ = Describe("DestinationGroups.Validate()", func() {
	var (
		destinationGroups DestinationGroups

		actualErr error
	)

	BeforeEach(func() {
		destinationGroups = nil
	})

	JustBeforeEach(func() {
		actualErr = destinationGroups.Validate()
	})

	When("the destinationGroups is nil", func() {
		BeforeEach(func() {
			destinationGroups = nil
		})

		It("should return an error", func() {
			Expect(actualErr).To(MatchError("at least 1 destination group is required"))
		})
	})

	When("the destinationGroups is empty", func() {
		BeforeEach(func() {
			destinationGroups = make([]DestinationGroup, 0)
		})

		It("should return an error", func() {
			Expect(actualErr).To(MatchError("at least 1 destination group is required"))
		})
	})

	Context("a valid destination group set (1 group)", func() {
		BeforeEach(func() {
			destinationGroups = []DestinationGroup{
				{
					Name:                 "foo",
					Destinations:         nil,
					DestinationNamespace: "",
					destinationsValidatorFunc: func(ds Destinations) error {
						return nil
					},
				},
			}
		})

		It("should not return an error", func() {
			Expect(actualErr).NotTo(HaveOccurred())
		})

		Context("except", func() {
			When("a destination group name is empty", func() {
				BeforeEach(func() {
					for i := range destinationGroups {
						destinationGroups[i].Name = ""
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating destination group: <empty>: name is required"))
				})
			})

			When("the destinations validator fails", func() {
				BeforeEach(func() {
					for i := range destinationGroups {
						destinationGroups[i].destinationsValidatorFunc = func(ds Destinations) error {
							return errors.New("destination validation error")
						}
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating destination group: foo: validating destinations: destination validation error"))
				})
			})
		})
	})

	When("multiple groups exist", func() {
		BeforeEach(func() {
			destinationGroups = []DestinationGroup{
				{
					Name:                 "",
					Destinations:         nil,
					DestinationNamespace: "",
					destinationsValidatorFunc: func(ds Destinations) error {
						return nil
					},
				},
				{
					Name:                 "",
					Destinations:         nil,
					DestinationNamespace: "",
					destinationsValidatorFunc: func(ds Destinations) error {
						return nil
					},
				},
				{
					Name:                 "",
					Destinations:         nil,
					DestinationNamespace: "",
					destinationsValidatorFunc: func(ds Destinations) error {
						return nil
					},
				},
			}
		})

		When("each group name is unique", func() {
			BeforeEach(func() {
				for i := range destinationGroups {
					destinationGroups[i].Name = "foo" + strconv.Itoa(i)
				}
			})

			It("should not return an error", func() {
				Expect(actualErr).NotTo(HaveOccurred())
			})

			When("the destinations validator fails", func() {
				BeforeEach(func() {
					for i := range destinationGroups {
						destinationGroups[i].destinationsValidatorFunc = func(ds Destinations) error {
							return errors.New("destination validation error")
						}
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError(
						strings.Join([]string{
							"validating destination group: foo0: validating destinations: destination validation error",
							"validating destination group: foo1: validating destinations: destination validation error",
							"validating destination group: foo2: validating destinations: destination validation error",
						}, "\n"),
					))
				})
			})
		})

		When("each there's a duplicate name", func() {
			BeforeEach(func() {
				destinationGroups[0].Name = "foo"
				destinationGroups[1].Name = "bar"
				destinationGroups[2].Name = "foo"
			})

			It("should return an error", func() {
				Expect(actualErr).To(MatchError("duplicate destination group name: foo"))
			})
		})
	})
})

var _ = Describe("Destinations.Validate()", func() {
	var (
		destinations Destinations
		actualErr    error
	)

	BeforeEach(func() {
		destinations = nil
		actualErr = nil
	})

	JustBeforeEach(func() {
		actualErr = destinations.Validate()
	})

	Context("a validate destinations set", func() {
		BeforeEach(func() {
			destinations = Destinations{
				{
					ApplicationDestination: &argov1alpha1.ApplicationDestination{
						Server:    "",
						Namespace: "my-namespace",
						Name:      "http://kube.local",
					},
					ArgoCDLogin: ArgoCDLogin{
						Hostname: "foo",
					},
				},
			}
		})

		It("should not return an error", func() {
			Expect(actualErr).NotTo(HaveOccurred())
		})

		When("server provided, name empty", func() {
			BeforeEach(func() {
				for i := range destinations {
					destinations[i].Server = "server-hostname"
					destinations[i].Name = ""
				}
			})

			It("should not return an error", func() {
				Expect(actualErr).NotTo(HaveOccurred())
			})
		})

		When("server empty, name provided", func() {
			BeforeEach(func() {
				for i := range destinations {
					destinations[i].Server = ""
					destinations[i].Name = "server-name"
				}
			})

			It("should not return an error", func() {
				Expect(actualErr).NotTo(HaveOccurred())
			})
		})

		Context("except", func() {
			When("argocdLogin.hostname is empty", func() {
				BeforeEach(func() {
					for i := range destinations {
						destinations[i].ArgoCDLogin.Hostname = ""
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating destination: argocdLogin.hostname is required"))
				})
			})

			When("server & name are both provided", func() {
				BeforeEach(func() {
					for i := range destinations {
						destinations[i].Server = "server-hostname"
						destinations[i].Name = "server-name"
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating destination: only one of clusterName or clusterUrl may be specified"))
				})
			})

			When("neither server or name are provided", func() {
				BeforeEach(func() {
					for i := range destinations {
						destinations[i].Server = ""
						destinations[i].Name = ""
					}
				})

				It("should return an error", func() {
					Expect(actualErr).To(MatchError("validating destination: one of clusterName or clusterUrl required"))
				})
			})
		})
	})
})
