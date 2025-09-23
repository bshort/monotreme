package license

import (
	v1pb "github.com/bshort/monotreme/proto/gen/api/v1"
)

type FeatureType string

const (
	// Enterprise features.

	// FeatureTypeSSO allows the user to use SSO.
	FeatureTypeSSO FeatureType = "ysh.monotreme.sso"
	// FeatureTypeAdvancedAnalytics allows the user to use advanced analytics.
	FeatureTypeAdvancedAnalytics FeatureType = "ysh.monotreme.advanced-analytics"

	// Usages.

	// FeatureTypeUnlimitedAccounts allows the user to create unlimited accounts.
	FeatureTypeUnlimitedAccounts FeatureType = "ysh.monotreme.unlimited-accounts"
	// FeatureTypeUnlimitedShortcuts allows the user to create unlimited shortcuts.
	FeatureTypeUnlimitedShortcuts FeatureType = "ysh.monotreme.unlimited-shortcuts"
	// FeatureTypeUnlimitedAccounts allows the user to create unlimited collections.
	FeatureTypeUnlimitedCollections FeatureType = "ysh.monotreme.unlimited-collections"

	// Customization.

	// FeatureTypeCustomeBranding allows the user to customize the branding.
	FeatureTypeCustomeBranding FeatureType = "ysh.monotreme.custom-branding"
)

func (f FeatureType) String() string {
	return string(f)
}

// FeatureMatrix is a matrix of features in [Free, Pro, Enterprise].
var FeatureMatrix = map[FeatureType][3]bool{
	FeatureTypeUnlimitedAccounts:    {false, true, true},
	FeatureTypeUnlimitedShortcuts:   {false, true, true},
	FeatureTypeUnlimitedCollections: {false, true, true},
	FeatureTypeCustomeBranding:      {false, false, true},
	FeatureTypeSSO:                  {false, false, true},
	FeatureTypeAdvancedAnalytics:    {false, false, true},
}

func getDefaultFeatures(plan v1pb.PlanType) []FeatureType {
	var features []FeatureType
	for feature, enabled := range FeatureMatrix {
		if enabled[plan-1] {
			features = append(features, feature)
		}
	}
	return features
}

func validateFeatureString(feature string) (FeatureType, bool) {
	switch feature {
	case "ysh.monotreme.unlimited-accounts":
		return FeatureTypeUnlimitedAccounts, true
	case "ysh.monotreme.unlimited-shortcuts":
		return FeatureTypeUnlimitedShortcuts, true
	case "ysh.monotreme.unlimited-collections":
		return FeatureTypeUnlimitedCollections, true
	case "ysh.monotreme.custom-branding":
		return FeatureTypeCustomeBranding, true
	case "ysh.monotreme.sso":
		return FeatureTypeSSO, true
	case "ysh.monotreme.advanced-analytics":
		return FeatureTypeAdvancedAnalytics, true
	default:
		return "", false
	}
}
