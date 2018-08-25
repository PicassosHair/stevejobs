package util

// RawPatentRecords defines raw PEDS bulk data.
type RawPatentRecords []struct {
	PatentCaseMetadata                     map[string]interface{}
	ProsecutionHistoryDataOrPatentTermData []struct {
		RecordedDate, CaseActionDescriptionText string
	}
}