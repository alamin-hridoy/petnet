package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultTimeout = 10 * time.Second

// pollOnError will continuously call f, sleeping for sleepTime after
// calls, until the error returned is nil, or the context is cancelled.
func pollOnError(ctx context.Context, logger *logrus.Entry, f func() error, sleepTime time.Duration) (err error) {
	for ctx.Err() == nil {
		if err = f(); err == nil {
			return nil
		}

		logger.Errorf("pollOnError inner func failed: %v\nsleeping...\n", err)
		time.Sleep(sleepTime)
	}

	logger.Errorf("ctx.Err: %v", ctx.Err())

	return err
}

func PostToSlack(logger *logrus.Entry, hookURL, text, channel, username, iconURL string) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return PostToSlackWithContext(ctx, logger, hookURL, text, channel, username, iconURL)
}

func PostToSlackWithContext(ctx context.Context, logger *logrus.Entry, hookURL, text, channel, username, iconURL string) error {
	return PostWithMrkDownWithContext(ctx, logger, hookURL, map[string]interface{}{
		"text":     text,
		"channel":  channel,
		"username": username,
		"icon_url": iconURL,
	})
}

func PostWithMrkDown(logger *logrus.Entry, hookURL string, fields map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	return PostWithMrkDownWithContext(ctx, logger, hookURL, fields)
}

func PostWithMrkDownWithContext(ctx context.Context, logger *logrus.Entry, hookURL string, fields map[string]interface{}) error {
	logger.Infof("slack post fields: %+v", fields)

	req, err := httpPost(ctx, hookURL, fields)
	if err != nil {
		return fmt.Errorf("httpPost err: %w", err)
	}
	respString, err := httpDo(req)
	if err != nil {
		return fmt.Errorf("slack post err: %w", err)
	}
	logger.Infof("slack post resp: %q", respString)
	return nil
}

// PostMrkDownWithRetries posts formatted content to slack with few retries.
func PostMrkDownWithRetries(logger *logrus.Entry, hookURL string, payload map[string]interface{}) {
	const (
		sleepTime    = 10 * time.Second
		slackTimeout = 1 * time.Minute
	)

	ctx, cancel := context.WithTimeout(context.Background(), slackTimeout)
	defer cancel()

	if err := pollOnError(ctx, logger, func() error {
		return PostWithMrkDown(logger, hookURL, payload)
	}, sleepTime); err != nil {
		logger.WithError(err).Error("failed to post to slack")
	}
}

// PostToSlackWithRetries posts text to slack with few retries.
func PostToSlackWithRetries(logger *logrus.Entry, hookURL, text, channel, username, iconURL string) {
	const (
		sleepTime    = 10 * time.Second
		slackTimeout = 1 * time.Minute
	)

	ctx, cancel := context.WithTimeout(context.Background(), slackTimeout)
	defer cancel()

	if err := pollOnError(ctx, logger, func() error {
		return PostToSlack(logger, hookURL, text, channel, username, iconURL)
	}, sleepTime); err != nil {
		logger.WithError(err).Error("failed to post to slack")
	}
}

func httpPost(ctx context.Context, url string, data map[string]interface{}) (*http.Request, error) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext err: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Brankas suite")
	return req, nil
}

func httpDo(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	strRespBody := string(respBody)

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("%s: %s", http.StatusText(resp.StatusCode), strRespBody)
	}
	return strRespBody, nil
}
