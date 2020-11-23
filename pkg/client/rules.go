package client

import (
	"context"
	"io/ioutil"
	"net/url"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/grafana/cortex-tools/pkg/rules/rwrulefmt"
)

const (
	rulerAPIPath    = "/api/v1/rules"
	legacyAPIPath   = "/api/prom/rules"
	rulerAllAPIPath = "/api/v1/allrules"
)

// CreateRuleGroup creates a new rule group
func (r *CortexClient) CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error {
	payload, err := yaml.Marshal(&rg)
	if err != nil {
		return err
	}

	path := rulerAPIPath
	if r.legacy {
		path = legacyAPIPath
	}
	escapedNamespace := url.PathEscape(namespace)
	path += "/" + escapedNamespace

	log.WithFields(log.Fields{
		"url": path,
	}).Debugln("path built to request rule group")

	res, err := r.doRequest(path, "POST", payload)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

// DeleteRuleGroup creates a new rule group
func (r *CortexClient) DeleteRuleGroup(ctx context.Context, namespace, groupName string) error {
	path := rulerAPIPath
	if r.legacy {
		path = legacyAPIPath
	}
	escapedNamespace := url.PathEscape(namespace)
	escapedGroupName := url.PathEscape(groupName)
	path += "/" + escapedNamespace + "/" + escapedGroupName

	log.WithFields(log.Fields{
		"url": path,
	}).Debugln("path built to request rule group")

	_, err := r.doRequest(path, "DELETE", nil)
	return err
}

// GetRuleGroup retrieves a rule group
func (r *CortexClient) GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error) {
	path := rulerAPIPath
	if r.legacy {
		path = legacyAPIPath
	}
	escapedNamespace := url.PathEscape(namespace)
	escapedGroupName := url.PathEscape(groupName)
	path += "/" + escapedNamespace + "/" + escapedGroupName

	log.WithFields(log.Fields{
		"url": path,
	}).Debugln("path built to request rule group")

	res, err := r.doRequest(path, "GET", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	rg := rwrulefmt.RuleGroup{}
	err = yaml.Unmarshal(body, &rg)
	if err != nil {
		log.WithFields(log.Fields{
			"body": string(body),
		}).Debugln("failed to unmarshal rule group from response")

		return nil, errors.Wrap(err, "unable to unmarshal response")
	}

	return &rg, nil
}

// ListRules retrieves a rule group
func (r *CortexClient) ListRules(ctx context.Context, namespace string) (map[string][]rwrulefmt.RuleGroup, error) {
	path := rulerAPIPath
	if r.legacy {
		path = legacyAPIPath
	}
	if namespace != "" {
		path = path + "/" + namespace
	}

	log.WithFields(log.Fields{
		"url": path,
	}).Debugln("path built to request rule group")

	res, err := r.doRequest(path, "GET", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	ruleSet := map[string][]rwrulefmt.RuleGroup{}
	err = yaml.Unmarshal(body, &ruleSet)
	if err != nil {
		return nil, err
	}

	return ruleSet, nil
}

// ListAllRules retrieves all user rule group
func (r *CortexClient) ListAllRules(ctx context.Context) (map[string]map[string][]rwrulefmt.RuleGroup, error) {
	path := rulerAllAPIPath

	log.WithFields(log.Fields{
		"url": path,
	}).Debugln("path built to request rule group")

	res, err := r.doRequest(path, "GET", nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	ruleSet := map[string]map[string][]rwrulefmt.RuleGroup{}
	err = yaml.Unmarshal(body, &ruleSet)
	if err != nil {
		return nil, err
	}

	return ruleSet, nil
}
