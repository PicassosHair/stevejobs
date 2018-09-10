package util

// RawPatentRecords defines raw PEDS bulk data.
type RawPatentRecords []struct {
	PatentCaseMetadata        map[string]interface{}
	ProsecutionHistoryDataBag struct {
		ProsecutionHistoryData []struct {
			EventDate, EventDescriptionText string
		}
	}
}
