package repositories

import configcat "github.com/configcat/go-sdk/v9"

type FeatureFlagRepository interface {
	HasFlagEnabled(email, flagName string) bool
}

func NewFeatureFlagRepository(configcatClient *configcat.Client) FeatureFlagRepository {
	return &featureFlagRepository{configcatClient: configcatClient}
}

type featureFlagRepository struct {
	configcatClient *configcat.Client
}

func (r *featureFlagRepository) HasFlagEnabled(email, flagName string) bool {
	user := &configcat.UserData{Email: email}
	return r.configcatClient.GetBoolValue(flagName, false, user)
}
