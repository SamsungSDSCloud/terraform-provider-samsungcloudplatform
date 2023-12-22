package dns

import (
	"context"
	sdk "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/client"
	dns2 "github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/dns2"
)

type Client struct {
	config    *sdk.Configuration
	sdkClient *dns2.APIClient
}

func NewClient(config *sdk.Configuration) *Client {
	return &Client{
		config:    config,
		sdkClient: dns2.NewAPIClient(config),
	}
}

func (client *Client) CreateDnsDomain(ctx context.Context, request CreateDnsDomainRequest, tags map[string]interface{}) (dns2.AsyncResponse, int, error) {

	result, c, err := client.sdkClient.DnsOpenApiV3ControllerApi.CreateDnsDomain1(ctx, client.config.ProjectId, dns2.CreateDnsDomainOpenApiV3Request{
		CountryType:       "DOMESTIC",
		DnsDomainName:     request.DnsDomainName,
		DnsEnvUsage:       "PRIVATE",
		EnCompanyName:     "Samsung SDS",
		KoCompanyName:     "삼성에스디에스",
		RegisteredByEmail: "cccs@samsung.com",
		RegisteredByTelno: "070-7010-3000",
		EnDetailAddress:   "SDS Campus",
		KoDetailAddress:   "SDS 캠퍼스",
		FirstEnAddress:    "125 Olympic-ro 35-gil, Songpa-gu, Seoul",
		FirstKoAddress:    "서울특별시 송파구 올림픽로35길 125 (신천동)",
		SecondEnAddress:   "",
		SecondKoAddress:   "",
		PostalCode:        "05510",
		DnsDescription:    request.DnsDescription,
		Tags:              client.sdkClient.ToTagRequestList(tags),
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) CreateDnsRecord(ctx context.Context, dnsDomainId string, request CreateDnsRecordRequest) (dns2.AsyncResponse, int, error) {
	dnsRecordMapping := make([]dns2.DnsRecordMappingRequest, 0)
	for _, recordInfo := range request.DnsRecordMapping {
		dnsRecordMapping = append(dnsRecordMapping, dns2.DnsRecordMappingRequest{
			RecordDestination: recordInfo.RecordDestination,
			Preference:        recordInfo.Preference,
		})
	}

	result, c, err := client.sdkClient.DnsOpenApiV3ControllerApi.CreateDnsRecord(ctx, client.config.ProjectId, dnsDomainId, dns2.CreateDnsRecordOpenApiV3Request{
		DnsRecordType:    request.DnsRecordType,
		DnsRecordName:    request.DnsRecordName,
		Ttl:              request.Ttl,
		DnsRecordMapping: dnsRecordMapping,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetDnsDomainDetail(ctx context.Context, dnsDomainId string) (dns2.DnsDomainServiceDetailResponse, int, error) {
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.DetailDnsDomain(ctx, client.config.ProjectId, dnsDomainId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetDnsRecordDetail(ctx context.Context, dnsDomainId string, dnsRecordId string) (dns2.DnsDomainRecordDetailResponse, int, error) {
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.DetailDnsRecord(ctx, client.config.ProjectId, dnsDomainId, dnsRecordId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) UpdateDnsDomainAutoExtension(ctx context.Context, dnsDomainId string, autoExtension *bool) (dns2.DnsDomainServiceDetailResponse, error) {
	result, _, err := client.sdkClient.DnsOpenApiV2ControllerApi.ChangeAutoExtensionOfDomain(ctx, client.config.ProjectId, dnsDomainId, dns2.ChangeAutoExtensionRequest{
		AutoExtension: autoExtension,
	})
	return result, err
}

func (client *Client) UpdateDnsDomainDescription(ctx context.Context, dnsDomainId string, description string) (dns2.DnsDomainServiceDetailResponse, error) {
	result, _, err := client.sdkClient.DnsOpenApiV2ControllerApi.ChangeDescriptionOfDomain(ctx, client.config.ProjectId, dnsDomainId, dns2.ChangeDnsDomainDescriptionRequest{
		Description: description,
	})
	return result, err
}

func (client *Client) UpdateDnsRecord(ctx context.Context, dnsDomainId string, dnsRecordId string, request ChangeDnsRecordRequest) (dns2.AsyncResponse, int, error) {
	dnsRecordMapping := make([]dns2.DnsRecordMappingRequest, 0)
	for _, recordInfo := range request.DnsRecordMapping {
		dnsRecordMapping = append(dnsRecordMapping, dns2.DnsRecordMappingRequest{
			RecordDestination: recordInfo.RecordDestination,
			Preference:        recordInfo.Preference,
		})
	}

	result, c, err := client.sdkClient.DnsOpenApiV3ControllerApi.ChangeDnsRecord1(ctx, client.config.ProjectId, dnsDomainId, dnsRecordId, dns2.ChangeDnsRecordOpenApiV3Request{
		Ttl:              request.Ttl,
		DnsRecordMapping: dnsRecordMapping,
	})
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteDnsDomain(ctx context.Context, dnsDomainId string) (dns2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.DeleteDnsDomain(ctx, client.config.ProjectId, dnsDomainId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) DeleteDnsRecord(ctx context.Context, dnsDomainId string, dnsRecordId string) (dns2.AsyncResponse, int, error) {
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.DeleteDnsRecord(ctx, client.config.ProjectId, dnsDomainId, dnsRecordId)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetDnsDomainList(ctx context.Context, request *dns2.DnsOpenApiV2ControllerApiListDnsDomainOpts) (dns2.ListResponseOfDnsDomainServiceListItemResponse, int, error) {
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.ListDnsDomain(ctx, client.config.ProjectId, request)
	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}

func (client *Client) GetDnsRecordList(ctx context.Context, dnsDomainId string, request *dns2.DnsOpenApiV2ControllerApiListDnsRecordOpts) (dns2.ListResponseOfDnsDomainRecordListItemResponse, int, error) {
	//fmt.Printf("GetDnsRecordList-request: %+v\n", request)
	result, c, err := client.sdkClient.DnsOpenApiV2ControllerApi.ListDnsRecord(ctx, client.config.ProjectId, dnsDomainId, request)

	var statusCode int
	if c != nil {
		statusCode = c.StatusCode
	}
	return result, statusCode, err
}
