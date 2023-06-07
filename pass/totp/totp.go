package totp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xlzd/gotp"
)

const (
	DefaultDigits = "6"
	DefaultPeriod = "30"
)

var Cmd = &cli.Command{
	Name:        "totp",
	Usage:       "generate Time-based One Time Password(TOTP)",
	Subcommands: []*cli.Command{showCmd, copyCmd},
}

func GenOTP(uri string) (string, error) {
	uri = strings.TrimSpace(uri)
	u, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("failed to parse URI: %w", err)
	}

	if !isTOTP(uri) {
		return "", fmt.Errorf(
			"unsupported or missing TYPE field in uri",
		)
	}

	if u.Scheme != "otpauth" {
		return "", fmt.Errorf("not an otpauth uri")
	}

	val := u.Query()

	secret := val.Get("secret")
	if secret == "" {
		return "", fmt.Errorf("missing secret in uri")
	}
	secret = strings.ToUpper(secret)

	algo := val.Get("algorithm")

	_digits := val.Get("digits")
	if _digits == "" {
		_digits = DefaultDigits
	}

	_period := val.Get("period")
	if _period == "" {
		_period = DefaultPeriod
	}

	digits, err := strconv.Atoi(_digits)
	if err != nil {
		return "", err
	}

	period, err := strconv.Atoi(_period)
	if err != nil {
		return "", err
	}

	h := &gotp.Hasher{HashName: algo}

	switch algo {
	case "":
		h = nil
	case "SHA1", "sha1":
		h.Digest = sha1.New
	case "SHA256", "sha256":
		h.Digest = sha256.New
	case "SHA512", "sha512":
		h.Digest = sha512.New
	default:
		return "", fmt.Errorf("hashing algorithm not supported %q", algo)
	}

	if !gotp.IsSecretValid(secret) {
		return "", fmt.Errorf("secret %q not valid", secret)
	}

	totp := gotp.NewTOTP(secret, digits, period, h)
	return totp.Now(), nil
}

func isTOTP(uri string) bool {
	reg := regexp.MustCompile(`otpauth://totp/.*`)
	return reg.MatchString(uri)
}
