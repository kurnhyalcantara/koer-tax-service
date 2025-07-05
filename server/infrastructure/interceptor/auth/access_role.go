package interceptor

const apiServicePath = "/koer_tax.service.v3.KoerTaxService/"

var accessibleRoles = map[string][]string{
	apiServicePath + "": {"data_entry:maker"},
}
