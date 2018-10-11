package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	completeSample = `[{"patentCaseMetadata":{"applicationNumberText":{"value":"15861673","electronicText":"15861673"},"filingDate":"2018-01-04","applicationTypeCategory":"Utility","partyBag":{"applicantBagOrInventorBagOrOwnerBag":[{"primaryExaminerOrAssistantExaminerOrAuthorizedOfficer":[{"name":{"personNameOrOrganizationNameOrEntityName":[{"personFullName":"NEGIN, RUSSELL SCOTT"}]}}]},{"inventorOrDeceasedInventor":[{"contactOrPublicationContact":[{"name":{"personNameOrOrganizationNameOrEntityName":[{"personStructuredName":{"lastName":"BERMEISTER, Kevin"}}]},"cityName":"Sherman Oaks","geographicRegionName":{"value":"CA","geographicRegionCategory":"STATE"},"countryCode":"US"}]}]},{"applicant":[{"contactOrPublicationContact":[{"name":{"personNameOrOrganizationNameOrEntityName":[{"personStructuredName":{"firstName":"","middleName":"CODONDEX LLC","lastName":""}}]},"cityName":"Sherman Oaks","geographicRegionName":{"value":"CA","geographicRegionCategory":"STATE"},"countryCode":""}]}]},{"partyIdentifierOrContact":[{"contactText":{"content":["75948"]}}]}]},"groupArtUnitNumber":{"value":"1631","electronicText":"1631"},"applicationConfirmationNumber":"7469","applicantFileReference":"4040-0008-U2","patentClassificationBag":{"cpcClassificationBagOrIPCClassificationOrECLAClassificationBag":[{"mainNationalClassification":{"nationalClass":"703","nationalSubclass":"703/011000"}}]},"businessEntityStatusCategory":"SMALL","firstInventorToFileIndicator":"true","inventionTitle":{"content":["SYSTEMS, METHODS, AND DEVICES FOR ANALYSIS OF GENETIC MATERIAL"]},"applicationStatusCategory":"Docketed New Case - Ready for Examination","applicationStatusDate":"2018-06-13","officialFileLocationCategory":"ELECTRONIC","patentPublicationIdentification":{"ipOfficeCode":"US","publicationNumber":"0203975","patentDocumentKindCode":"A1","publicationDate":"2018-07-19"},"patentGrantIdentification":{},"hagueAgreementData":{}},"prosecutionHistoryDataBag":{"prosecutionHistoryData":[{"eventDate":"2018-07-19","eventDescriptionText":"PG-Pub Issue Notification , PG-ISSUE"},{"eventDate":"2018-06-13","eventDescriptionText":"Case Docketed to Examiner in GAU , DOCK"},{"eventDate":"2018-04-11","eventDescriptionText":"Filing Receipt - Corrected , FLRCPT.C"},{"eventDate":"2018-03-16","eventDescriptionText":"Application Dispatched from OIPE , OIPE"},{"eventDate":"2018-03-16","eventDescriptionText":"FITF set to YES - revise initial setting , FTFS"},{"eventDate":"2018-03-14","eventDescriptionText":"Patent Term Adjustment - Ready for Examination , PTA.RFE"},{"eventDate":"2018-03-19","eventDescriptionText":"Application Is Now Complete , COMP"},{"eventDate":"2018-03-19","eventDescriptionText":"Filing Receipt - Updated , FLRCPT.U"},{"eventDate":"2018-03-14","eventDescriptionText":"Payment of additional filing fee/Preexam , FLFEE"},{"eventDate":"2018-03-01","eventDescriptionText":"Notice of Incomplete Reply , INCR"},{"eventDate":"2018-02-19","eventDescriptionText":"A set of symbols and procedures, provided to the PTO on a set of computer listings, that describe in , SEQLIST"},{"eventDate":"2018-02-19","eventDescriptionText":"CRF Disk Has Been Received by Preexam / Group / PCT , CRFL"},{"eventDate":"2018-02-22","eventDescriptionText":"CRF Is Good Technically / Entered into Database , CRFE"},{"eventDate":"2018-02-08","eventDescriptionText":"Application ready for PDX access by participating foreign offices , CCRDY"},{"eventDate":"2018-02-07","eventDescriptionText":"Notice Mailed--Application Incomplete--Filing Date Assigned  , INCD"},{"eventDate":"2018-02-07","eventDescriptionText":"Filing Receipt , FLRCPT.O"},{"eventDate":"2018-01-04","eventDescriptionText":"PTO/SB/69-Authorize EPO Access to Search Results , SREXR141"},{"eventDate":"2018-01-04","eventDescriptionText":"Applicants have given acceptable permission for participating foreign  , APPERMS"},{"eventDate":"2018-02-07","eventDescriptionText":"Applicant Has Filed a Verified Statement of Small Entity Status in Compliance with 37 CFR 1.27 , SMAL"},{"eventDate":"2018-01-23","eventDescriptionText":"Cleared by L&R (LARS) , L128"},{"eventDate":"2018-01-22","eventDescriptionText":"Referred to Level 2 (LARS) by OIPE CSR , L198"},{"eventDate":"2018-01-04","eventDescriptionText":"IFW Scan & PACR Auto Security Review , SCAN"},{"eventDate":"2018-01-04","eventDescriptionText":"ENTITY STATUS SET TO UNDISCOUNTED (INITIAL DEFAULT SETTING OR STATUS CHANGE) , BIG."},{"eventDate":"2018-01-04","eventDescriptionText":"Initial Exam Team nn , IEXX"}]},"st96Version":"V3_0","ipoVersion":"US_V8_0_D7"}]`
)

func convertStringToRecord(rawString string) RawPatentRecords {
	var rawRecord RawPatentRecords
	json.Unmarshal([]byte(rawString), &rawRecord)
	return rawRecord
}

func TestExtractApplID(t *testing.T) {
	testRecord := convertStringToRecord(completeSample)
	applIDText := ExtractApplID(&testRecord)
	fmt.Println(applIDText)
	if applIDText != "15861673" {
		t.Error("Title failed")
	}
}
