package util

import (
	"encoding/json"
)

// PatentCaseMetadata defines main meta data info.
type PatentCaseMetadata struct {
	ApplicationNumberText struct {
		Value, ElectronicText string
	} `json:"applicationNumberText"`
	FilingDate              string
	ApplicationTypeCategory string
	PartyBag                struct {
		// Could be primaryExaminerOrAssistantExaminerOrAuthorizedOfficer, applicant, inventorOrDeceasedInventor.
		ApplicantBagOrInventorBagOrOwnerBag []map[string]*json.RawMessage
	}
	GroupArtUnitNumber struct {
		Value, ElectronicText string
	}
	ApplicationConfirmationNumber, ApplicantFileReference string
	PriorityClaimBag                                      struct {
		PriorityClaim []struct {
			IPOfficeName      string `json:"ipOfficeName"`
			FilingDate        string `json:"filingDate"`
			SequenceNumber    string `json:"sequenceNumber"`
			ApplicationNumber struct {
				ApplicationNumberText string
			}
		}
	}
	PatentClassificationBag struct {
		CpcClassificationBagOrIPCClassificationOrECLAClassificationBag []struct {
			IPOfficeCode               string `json:"ipOfficeName"`
			MainNationalClassification struct {
				NationalClass, NationalSubclass string
			}
		}
	}
	BusinessEntityStatusCategory, FirstInventorToFileIndicator string
	InventionTitle                                             struct {
		Content []string
	}
	ApplicationStatusCategory       string `json:"applicationStatusCategory"`
	ApplicationStatusDate           string `json:"applicationStatusDate"`
	OfficialFileLocationCategory    string `json:"officialFileLocationCategory"`
	PatentPublicationIdentification struct {
		PublicationNumber string `json:"publicationNumber"`
		PublicationDate   string `json:"publicationDate"`
	} `json:"patentPublicationIdentification"`
	PatentGrantIdentification struct {
		PatentNumber string `json:"patentNumber"`
		GrantDate    string `json:"grantDate"`
	} `json:"patentGrantIdentification"`
	RelatedDocumentData struct {
		ParentDocumentDataOrChildDocumentData []struct {
			DescriptionText          string `json:"descriptionText"`
			ApplicationNumberText    string `json:"applicationNumberText"`
			FilingDate               string `json:"filingDate"`
			AiaIndicator             bool   `json:"aiaIndicator"`
			ParentDocumentStatusCode string `json:"parentDocumentStatusCode"`
			PatentNumber             string `json:"patentNumber"`
		} `json:"parentDocumentDataOrChildDocumentData"`
	} `json:"relatedDocumentData"`
}

// ProsecutionHistoryDataBag defines transaction histories.
type ProsecutionHistoryDataBag struct {
	ProsecutionHistoryData []struct {
		EventDate            string `json:"eventDate"`
		EventCode            string `json:"eventCode"`
		EventDescriptionText string `json:"eventDescriptionText"`
	}
}

type AssignmentContact struct {
  ExecutionDate string `json:"executionDate"`
  ContactOrPublicationContact Contact `json:"contactOrPublicationContact"`
}

type Assignment struct {
  ReelNumber string `json:"reelNumber"`
  FrameNumber string `json:"frameNumber"`
  DocumentReceivedDate string `json:"documentReceivedDate"`
  RecordedDate string `json:"recordedDate"`
  MailDate string `json:"mailDate"`
  PageTotalQuantity int `json:"pageTotalQuantity"`
  ConveyanceText string `json:"conveyanceText"`
  AssignorBag struct {
    Assignor []AssignmentContact `json:"assignor"`
  } `json:"assignorBag"`
  AssigneeBag struct {
    Assignee []AssignmentContact `json:"assignee"`
  } `json:"assigneeBag"`
  CorrespondenceAddress  `json:"correspondenceAddress"`
}

type AssignmentDataBag struct {
  AssignmentData []Assignment `json:"assignmentData"`
  AssignmentTotalQuantity int `json:"assignmentTotalQuantity"`
}

// RawPatentRecord defines raw PEDS bulk data.
// View sample_record.json.
type RawPatentRecord struct {
	PatentCaseMetadata        PatentCaseMetadata        `json:"patentCaseMetadata"`
	ProsecutionHistoryDataBag ProsecutionHistoryDataBag `json:"prosecutionHistoryDataBag"`
  AssignmentDataBag AssignmentDataBag `json:"assignmentDataBag"`

	St96Version string `json:"st96Version"`
	IpoVersion  string `json:"ipoVersion"`
}

// EntityName defines a person or an entity name, full or structured.
type EntityName struct {
	PersonFullName       string `json:"personFullName"`
	PersonStructuredName struct {
		FirstName  string `json:"firstName"`
		MiddleName string `json:"middleName"`
		LastName   string `json:"lastName"`
		NameSuffix string `json:"nameSuffix"`
	} `json:"personStructuredName"`
	OrganizationStandardName struct {
		Content []string `json:"content"`
	} `json:"organizationStandardName"`
}

// Contact defines an entity's contact information.
type Contact struct {
	Name struct {
		PersonNameOrOrganizationNameOrEntityName []EntityName `json:"personNameOrOrganizationNameOrEntityName"`
	} `json:"name"`
	PhoneNumberBag struct {
		PhoneNumber []struct {
			Value string `json:"value"`
		} `json:"phoneNumber"`
	} `json:"phoneNumberBag"`
	CityName             string `json:"cityName"`
	GeographicRegionName struct {
		Value                    string `json:"value"`
		GeographicRegionCategory string `json:"geographicRegionCategory"`
	} `json:"geographicRegionName"`
	CountryCode string `json:"countryCode"`
  // Disable for now.
  // PostalAddressBag struct {
  //   PostalAddress []struct {
  //     PostalAddressText []struct {
  //       SequenceNumber string `json:"sequenceNumber"`
  //       Value string `json:"value"`
  //     } `json:"postalAddressText"`
  //   } `json:"postalAddress"`
  // } `json:"postalAddressBag"`
}

// Examiner defines some examiners.
type Examiner []Contact

// Applicant defines some applicants.
type Applicant []struct {
	ContactOrPublicationContact []Contact `json:"contactOrPublicationContact"`
}

// Inventor defines some inventors.
type Inventor []struct {
	ContactOrPublicationContact []Contact `json:"contactOrPublicationContact"`
}

// Practitioner defines some practitioners.
type Practitioner []struct {
	RegisteredPractitionerRegistrationNumber string    `json:"registeredPractitionerRegistrationNumber"`
	RegisteredPractitionerCategory           string    `json:"registeredPractitionerCategory"`
	ContactOrPublicationContact              []Contact `json:"contactOrPublicationContact"`
	ActiveIndicator                          bool      `json:"activeIndicator"`
}
