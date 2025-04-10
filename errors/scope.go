package errors

func userScope(entries []Entry) *DictT {
	return Dict("user", entries, Public)
}

func systemScope(entries []Entry) *DictT {
	caller := "unknown"
	if _, _, name, ok := findCaller(); ok {
		caller = name
	}

	call := Dict(caller, entries, Restrict)
	return Dict("system", []Entry{call}, Restrict)
}
