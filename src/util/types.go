package util

// RawPatentRecords defines raw PEDS bulk data.
// View sample_record to see its structure.
type RawPatentRecords []struct {
	PatentCaseMetadata        map[string]interface{}
	ProsecutionHistoryDataBag struct {
		ProsecutionHistoryData []struct {
			EventDate, EventDescriptionText string
		}
	}
}
