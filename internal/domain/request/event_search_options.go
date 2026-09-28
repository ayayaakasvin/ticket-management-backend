package request

type EventSearchOptions struct {
	Status 			string 		// default "Active"
	Tags 			[]string 	// {music, dance, batchata}
	OrganizerID		uint64		// in case if user wants to find specific person events
	
	Category		uint
	
	Latitude		float64
	Longitude		float64
}