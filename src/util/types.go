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
      ApplicantBagOrInventorBagOrOwnerBag []json.RawMessage
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
