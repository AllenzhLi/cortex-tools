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

const rulerAllAPIPath = "/ruler/rules"

// CreateRuleGroup creates a new rule group
func (r *CortexClient) CreateRuleGroup(ctx context.Context, namespace string, rg rwrulefmt.RuleGroup) error {
	payload, err := yaml.Marshal(&rg)
	if err != nil {
		return err
	}

	escapedNamespace := url.PathEscape(namespace)
	path := r.apiPath + "/" + escapedNamespace

	res, err := r.doRequest(path, "POST", payload)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

// DeleteRuleGroup creates a new rule group
func (r *CortexClient) DeleteRuleGroup(ctx context.Context, namespace, groupName string) error {
	escapedNamespace := url.PathEscape(namespace)
	escapedGroupName := url.PathEscape(groupName)
	path := r.apiPath + "/" + escapedNamespace + "/" + escapedGroupName

	_, err := r.doRequest(path, "DELETE", nil)
	return err
}

// GetRuleGroup retrieves a rule group
func (r *CortexClient) GetRuleGroup(ctx context.Context, namespace, groupName string) (*rwrulefmt.RuleGroup, error) {
	escapedNamespace := url.PathEscape(namespace)
	escapedGroupName := url.PathEscape(groupName)
	path := r.apiPath + "/" + escapedNamespace + "/" + escapedGroupName

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
	path := r.apiPath
	if namespace != "" {
		path = path + "/" + namespace
	}

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
