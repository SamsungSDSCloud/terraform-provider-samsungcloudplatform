package dns

type CreateDnsDomainRequest struct {
	// DNS Domain Name
	DnsDomainName string
	// DNS Domain Address Type(DOMESTIC|OVERSEA)
	CountryType string
	// English Name of Company
	EnCompanyName string
	// Korean Name of Company
	KoCompanyName string
	// Email
	RegisteredByEmail string
	// Tel
	RegisteredByTelno string
	// English Detail Address
	EnDetailAddress string
	// Korean Detail Address
	KoDetailAddress string
	// English First Address
	FirstEnAddress string
	// Korean First Address
	FirstKoAddress string
	// English Second Address (Complete only for overseas address (OVERSEA))
	SecondEnAddress string
	// Korean Second Address (Complete only for overseas address (OVERSEA))
	SecondKoAddress string
	// Postal Code
	PostalCode string
	// DNS Domain Description
	DnsDescription string
	Tags           []TagRequest
}

type CreateDnsRecordRequest struct {
	// DNS Record Type
	DnsRecordType string
	// DNS Record Name
	DnsRecordName string
	// DNS Record TTL
	Ttl int32
	// DNS Record Mapping
	DnsRecordMapping []DnsRecordMappingRequest
}

type ChangeDnsRecordRequest struct {
	// DNS Record Type
	DnsRecordType string
	// DNS Record TTL
	Ttl int32
	// DNS Record Mapping
	DnsRecordMapping []DnsRecordMappingRequest
}

type DnsRecordMappingRequest struct {
	// DNS Record Mapping Preference
	Preference int32
	// DNS Record Destination
	RecordDestination string
}

type TagRequest struct {
	TagKey   string
	TagValue string
}
