package config

import (
	"context"
	"errors"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
)

// developmentSecretPrefix marks every secret committed in .env.example for local
// use; such a value is public, so production must never run on one.
const developmentSecretPrefix = "dev-insecure-"

// validate is the one validator: it is safe for concurrent use and caches the
// struct metadata, so it is built once. Rules sit in `validate` tags next to the
// `env` tag of the setting they constrain. A rule that tags express badly (one
// across settings) is a crossChecks method on the root; a section that matters
// only when its flag is on (OAuth, MCP) is excluded from the root's tags
// (`validate:"-"`) and checked by the root's crossChecks via section.
var validate = newValidator()

// appEnvKey carries the APP_ENV a check runs under into the custom rules.
type appEnvKey struct{}

// newValidator registers the two rules the library has no tag for.
func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// origin: scheme and host only, as in https://code-basics.com. A URL with a
	// path or query would be baked into tokens and links built from it.
	mustRegister(v, "origin", func(_ context.Context, fl validator.FieldLevel) bool {
		u, err := url.Parse(fl.Field().String())
		if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
			return false
		}
		return u.Path == "" && u.RawQuery == "" && u.Fragment == "" && u.User == nil
	})
	// notdevsecret: in production a secret must not be a public development value.
	mustRegister(v, "notdevsecret", func(ctx context.Context, fl validator.FieldLevel) bool {
		env, _ := ctx.Value(appEnvKey{}).(AppEnv)
		return env != Production || !strings.HasPrefix(fl.Field().String(), developmentSecretPrefix)
	})
	return v
}

func mustRegister(v *validator.Validate, tag string, fn validator.FuncCtx) {
	if err := v.RegisterValidationCtx(tag, fn); err != nil {
		panic("register validator " + tag + ": " + err.Error())
	}
}

// crossChecker is implemented by a root that has rules across settings.
type crossChecker interface {
	crossChecks() []string
}

// check validates a parsed root and returns one error naming every violated
// setting by its variable, or nil. A secret's value never appears in it: a
// message states the rule, not the input.
func check(root crossChecker, appEnv AppEnv) error {
	problems := violations(root, "", appEnv)
	problems = append(problems, root.crossChecks()...)
	if len(problems) == 0 {
		return nil
	}
	return oops.Errorf("invalid configuration:\n  %s", strings.Join(problems, "\n  "))
}

// section reports the violations of a nested config section whose env prefix is
// prefix; a root calls it for the sections it validates conditionally.
func section(value any, prefix string, appEnv AppEnv) []string {
	return violations(value, prefix, appEnv)
}

func violations(value any, prefix string, appEnv AppEnv) []string {
	err := validate.StructCtx(context.WithValue(context.Background(), appEnvKey{}, appEnv), value)
	if err == nil {
		return nil
	}
	var invalid validator.ValidationErrors
	if !errors.As(err, &invalid) {
		return []string{"config cannot be validated: " + err.Error()}
	}
	problems := make([]string, 0, len(invalid))
	for _, fe := range invalid {
		problems = append(problems, prefix+envName(reflect.TypeOf(value), fe.StructNamespace())+" "+describe(fe))
	}
	return problems
}

// describe states the violated rule without the offending value.
func describe(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "oneof":
		return "must be one of: " + fe.Param()
	case "origin":
		return "must be a bare origin such as https://code-basics.com (scheme and host, no path)"
	case "notdevsecret":
		return "must differ from the public development value"
	case "min", "gte":
		return "must be at least " + fe.Param()
	default:
		return "is invalid (" + fe.Tag() + ")"
	}
}

// envName turns a validator namespace such as "Server.OAuth.Issuer" into the
// variable it came from (OAUTH_ISSUER), following the env and envPrefix tags,
// so a message names what the operator sets, not a Go field.
func envName(root reflect.Type, namespace string) string {
	segments := strings.Split(namespace, ".")[1:] // the first is the root type
	t := root
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	prefix := ""
	for i, seg := range segments {
		field, ok := t.FieldByName(seg)
		if !ok {
			return namespace
		}
		if name, _, _ := strings.Cut(field.Tag.Get("env"), ","); name != "" && i == len(segments)-1 {
			return prefix + name
		}
		prefix += field.Tag.Get("envPrefix")
		t = field.Type
	}
	return namespace
}
