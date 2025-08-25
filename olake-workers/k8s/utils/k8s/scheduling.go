package k8s

import (
	"fmt"
	"strings"

	appConfig "github.com/datazip-inc/olake-ui/olake-workers/k8s/config"
	"github.com/datazip-inc/olake-ui/olake-workers/k8s/logger"

	"k8s.io/apimachinery/pkg/util/validation"
)

// validateLabelPair validates a single key-value label pair
func validateLabelPair(key, value string) error {
	if key == "" {
		return fmt.Errorf("empty label key")
	}

	if value == "" {
		return fmt.Errorf("empty label value for key '%s'", key)
	}

	if errs := validation.IsQualifiedName(key); len(errs) > 0 {
		return fmt.Errorf("invalid label key '%s': %v", key, errs)
	}

	if errs := validation.IsValidLabelValue(value); len(errs) > 0 {
		return fmt.Errorf("invalid label value '%s' for key '%s': %v", value, key, errs)
	}

	return nil
}

// GetValidJobMapping loads the JobID to node mapping configuration from config
// with enhanced error handling and detailed validation
func GetValidJobMapping(cfg *appConfig.Config) map[int]map[string]string {
	result := make(map[int]map[string]string)
	for jobID, nodeLabels := range cfg.Kubernetes.JobMapping {
		if jobID <= 0 {
			logger.Warnf("Invalid JobID: %d", jobID)
			continue
		}

		// Handle null/empty mappings
		if nodeLabels == nil {
			logger.Warnf("jobID[%d]: null mapping", jobID)
			continue
		}

		validMapping := make(map[string]string)
		for key, value := range nodeLabels {
			if err := validateLabelPair(strings.TrimSpace(key), strings.TrimSpace(value)); err != nil {
				logger.Warnf("jobID[%d]: found invalid label pair: %s", jobID, err)
				continue
			}

			validMapping[key] = value
		}
	}

	return result
}
