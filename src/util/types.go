package util

import (
	"encoding/json"
)


// RawPatentRecords defines raw PEDS bulk data.
// View sample_record to see its structure.
type RawPatentRecords []struct {
	PatentCaseMetadata struct {
    ApplicationNumberText struct {
      Value, ElectronicText string
    } `json:"applicationNumberText"`
    FilingDate string
    ApplicationTypeCategory string
    PartyBag struct {
      // Could be primaryExaminerOrAssistantExaminerOrAuthorizedOfficer, applicant, inventorOrDeceasedInventor.
      ApplicantBagOrInventorBagOrOwnerBag []map[string]*json.RawMessage
    }
    GroupArtUnitNumber struct {
      Value, ElectronicText string
    }
    ApplicationConfirmationNumber, ApplicantFileReference string
    PriorityClaimBag struct {
      PriorityClaim []struct {
        IpOfficeName, FilingDate, SequenceNumber string
        ApplicationNumber struct {
          ApplicationNumberText string
        }
      }
    }
    PatentClassificationBag struct {
      CpcClassificationBagOrIPCClassificationOrECLAClassificationBag []struct {
        IpOfficeCode string
        MainNationalClassification struct {
          NationalClass, NationalSubclass string
        }
      }
    }
    BusinessEntityStatusCategory, FirstInventorToFileIndicator string
    InventionTitle struct {
      Content []string
    }
    ApplicationStatusCategory, ApplicationStatusDate, OfficialFileLocationCategory string
    RelatedDocumentData struct {
      ParentDocumentDataOrChildDocumentData []struct {
        DescriptionText, ApplicationNumberText, FilingDate, ParentDocumentStatusCode, PatentNumber string
      }
    }
  }

	ProsecutionHistoryDataBag struct {
		ProsecutionHistoryData []struct {
			EventDate, EventDescriptionText string
		}
	}

  St96Version, IpoVersion string
}

type EntityName struct{
  PersonFullName string `json:"personFullName"`
  PersonStructuredName struct {
    FirstName string `json:"firstName`
    MiddleName string `json:"middleName"`
    LastName string `json:"lastName"`
  } `json:"personStructuredName"`
}

type Contact struct {
  Name struct {
    PersonNameOrOrganizationNameOrEntityName []EntityName `json:"personNameOrOrganizationNameOrEntityName"`
  } `json:"name"`
  PhoneNumberBag struct {
    PhoneNumber []struct {
      Value string `json:"value"`
    } `json:"phoneNumber"`
  } `json:"phoneNumberBag"`
  CityName string `json:"cityName"`
  GeographicRegionName struct {
    Value string `json:"value"`
    GeographicRegionCategory string `json:"geographicRegionCategory"`
  } `json:"geographicRegionName"`
  CountryCode string `json:"countryCode"`
}

type Examiner []Contact

type Applicant []struct {
  ContactOrPublicationContact []Contact `json:"contactOrPublicationContact"`
}

type Inventor []struct {
  ContactOrPublicationContact []Contact `json:"contactOrPublicationContact"`
}

type Practitioner []struct {
  RegisteredPractitionerRegistrationNumber string `json:"registeredPractitionerRegistrationNumber"`
  RegisteredPractitionerCategory string `json:"registeredPractitionerCategory"`
  ContactOrPublicationContact []Contact `json:"contactOrPublicationContact"`
  ActiveIndicator bool `json:"activeIndicator"`
}
