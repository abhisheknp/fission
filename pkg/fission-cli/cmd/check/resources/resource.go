package resources

// Resource interface that every check should implement
type Resource interface {
	Check() Results
	GetLabel() string
}

// Result of check
type Result struct {
	Description string
	Ok          bool
}

// Results is used store multiple result
type Results []*Result

func getResults(desc string, ok bool) Results {
	return Results{
		{
			Description: desc,
			Ok:          ok,
		},
	}
}

func appendResult(results Results, desc string, ok bool) Results {
	res := &Result{
		Description: desc,
		Ok:          ok,
	}
	return append(results, res)
}
