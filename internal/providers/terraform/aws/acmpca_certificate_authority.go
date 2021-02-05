package aws

import (
	"github.com/infracost/infracost/internal/schema"
	"github.com/infracost/infracost/internal/usage"

	"github.com/shopspring/decimal"
)

func GetACMPCACertificateAuthorityRegistryItem() *schema.RegistryItem {
	return &schema.RegistryItem{
		Name:  "aws_acmpca_certificate_authority",
		RFunc: NewACMPCACertificateAuthority,
	}
}

func NewACMPCACertificateAuthority(d *schema.ResourceData, u *schema.UsageData) *schema.Resource {
	region := d.Get("region").String()

	costComponents := []*schema.CostComponent{
		{
			Name:            "Private certificate authority",
			Unit:            "months",
			UnitMultiplier:  1,
			MonthlyQuantity: decimalPtr(decimal.NewFromInt(1)),
			ProductFilter: &schema.ProductFilter{
				VendorName:    strPtr("aws"),
				Region:        strPtr(region),
				Service:       strPtr("AWSCertificateManager"),
				ProductFamily: strPtr("AWS Certificate Manager"),
				AttributeFilters: []*schema.AttributeFilter{
					{Key: "usagetype", ValueRegex: strPtr("/PaidPrivateCA/")},
				},
			},
		},
	}

	certificateTierLimits := []int{1000, 9000, 10000}
	if u != nil && u.Get("monthly_requests").Exists() {
		monthlyCertificatesRequests := decimal.NewFromInt(u.Get("monthly_requests").Int())
		privateCertificateTier := usage.CalculateTierRequests(monthlyCertificatesRequests, certificateTierLimits)
		tierOne := privateCertificateTier["1"]
		tierTwo := privateCertificateTier["2"]
		tierThree := privateCertificateTier["3"]

		if tierOne.GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (first 1K)", "0", &tierOne))
		}

		if tierTwo.GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (next 9K)", "1000", &tierTwo))
		}

		if tierThree.GreaterThan(decimal.NewFromInt(0)) {
			costComponents = append(costComponents, certificateCostComponent(region, "Certificates (over 10K)", "10000", &tierThree))
		}
	} else {
		var unknown *decimal.Decimal
		costComponents = append(costComponents, certificateCostComponent(region, "Certificates (first 1K)", "0", unknown))
	}

	return &schema.Resource{
		Name:           d.Address,
		CostComponents: costComponents,
	}
}

func certificateCostComponent(region string, displayName string, usageTier string, monthlyQuantity *decimal.Decimal) *schema.CostComponent {
	return &schema.CostComponent{
		Name:            displayName,
		Unit:            "requests",
		UnitMultiplier:  1,
		MonthlyQuantity: monthlyQuantity,
		ProductFilter: &schema.ProductFilter{
			VendorName:    strPtr("aws"),
			Region:        strPtr(region),
			Service:       strPtr("AWSCertificateManager"),
			ProductFamily: strPtr("AWS Certificate Manager"),
			AttributeFilters: []*schema.AttributeFilter{
				{Key: "usagetype", ValueRegex: strPtr("/PrivateCertificatesIssued/")},
			},
		},
		PriceFilter: &schema.PriceFilter{
			StartUsageAmount: strPtr(usageTier),
		},
	}
}