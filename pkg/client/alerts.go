package client

import (
	"context"
	"io/ioutil"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	alertmanagerAPIPath    = "/api/v1/alerts"
	alertmanagerAllAPIPath = "/multitenant_alertmanager/configs"
)

type configCompat struct {
	TemplateFiles      map[string]string `yaml:"template_files"`
	AlertmanagerConfig string            `yaml:"alertmanager_config"`
}

// CreateAlertmanagerConfig creates a new alertmanager config
func (r *CortexClient) CreateAlertmanagerConfig(ctx context.Context, cfg string, templates map[string]string) error {
	payload, err := yaml.Marshal(&configCompat{
		TemplateFiles:      templates,
		AlertmanagerConfig: cfg,
	})
	if err != nil {
		return err
	}

	_, err = r.doRequest(alertmanagerAPIPath, "POST", payload)
	return err
}

// DeleteAlermanagerConfig deletes the users alertmanager config
func (r *CortexClient) DeleteAlermanagerConfig(ctx context.Context) error {
	_, err := r.doRequest(alertmanagerAPIPath, "DELETE", nil)
	return err
}

// GetAlertmanagerConfig retrieves a user alertmanager config
func (r *CortexClient) GetAlertmanagerConfig(ctx context.Context) (string, map[string]string, error) {
	res, err := r.doRequest(alertmanagerAPIPath, "GET", nil)
	if err != nil {
		log.Debugln("no alert config present in response")
		return "", nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", nil, err
	}

	compat := configCompat{}
	err = yaml.Unmarshal(body, &compat)
	if err != nil {
		log.WithFields(log.Fields{
			"body": string(body),
		}).Debugln("failed to unmarshal rule group from response")

		return "", nil, errors.Wrap(err, "unable to unmarshal response")
	}

	return compat.AlertmanagerConfig, compat.TemplateFiles, nil
}

// ListAlertmanagerConfig retrieves all user alertmanager config
func (r *CortexClient) ListAlertmanagerConfig(ctx context.Context) (map[string]configCompat, error) {
	res, err := r.doRequest(alertmanagerAllAPIPath, "GET", nil)
	if err != nil {
		log.Debugln("no alert config present in response")
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	compatMap := make(map[string]configCompat)
	err = yaml.Unmarshal(body, &compatMap)
	if err != nil {
		log.WithFields(log.Fields{
			"body": string(body),
		}).Debugln("failed to unmarshal rule group from response")

		return nil, errors.Wrap(err, "unable to unmarshal response")
	}

	return compatMap, nil
}
