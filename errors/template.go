package errors

// TemplateT is a structure that defines a template for creating and propagating errors.
// It includes a kind, a message, and optional parameters (opts) that can be used to
// customize the error. The Template function initializes a new TemplateT instance.
// The New method creates a new error based on the template, while the Propagate method
// propagates an existing error with additional context provided by the template.
type TemplateT struct {
	kind    *Kind
	message string
	opts    []Opter
}

// Template initializes a new TemplateT instance, which is used to define a template
// for creating and propagating errors. It takes a kind, a message, and optional parameters
// (opts) that can be used to customize the error. This function returns a pointer to the
// TemplateT struct, allowing for the creation of new errors or the propagation of existing
// errors with additional context provided by the template.
func Template(kind *Kind, message string, opts ...Opter) *TemplateT {
	return &TemplateT{
		kind:    kind,
		message: message,
		opts:    opts,
	}
}

// New creates a new error based on the template, using the specified kind, message,
// and optional parameters (opts). It combines the template's options with any additional
// options provided, allowing for customization of the error. This method returns a pointer
// to an ErrorT instance, representing the newly created error.
func (t *TemplateT) New(opts ...Opter) *ErrorT {
	opts = append(opts, t.opts...)
	return New(t.kind, t.message, opts...)
}

// Propagate propagates an existing error with additional context provided by the template.
// It combines the template's options with any additional options provided, allowing for
// customization of the error propagation. This method returns a pointer to an ErrorT instance,
// representing the propagated error with the specified kind, message, and options.
func (t *TemplateT) Propagate(err error, opts ...Opter) *ErrorT {
	opts = append(opts, t.opts...)
	return PropagateAs(t.kind, err, t.message, opts...)
}
