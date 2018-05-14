//go:generate ./gentemplateapi.sh Name Tag User Domain
package emailfuncs

import (
	"net/mail"
	"strings"

	"github.com/pkg/errors"
)

const (
	tagOrDomainSeparators = tagSeparator + domainSeparator
	domainSeparator       = "@"
	tagSeparator          = "+"
)

type address struct {
	raw *mail.Address
}

func Parse(inaddr string) (*address, error) {
	mailaddr, err := mail.ParseAddress(inaddr)
	if err != nil {
		return nil, errors.Wrapf(err, "could not parse provided address %q", inaddr)
	}

	addr := &address{mailaddr}

	return addr, nil
}

func (a *address) Name() (string, error) {
	return a.raw.Name, nil
}

func (a *address) Domain() (string, error) {
	afterAt := strings.Index(a.raw.Address, domainSeparator)
	if afterAt < 0 || afterAt > len(a.raw.Address)-1 {
		return "", errors.Errorf("no domain part to email address %q", a.raw.Address)
	}
	domainPart := a.raw.Address[afterAt+1 : len(a.raw.Address)]
	return domainPart, nil
}

func (a *address) User() (string, error) {
	cutAt := strings.IndexAny(a.raw.Address, tagOrDomainSeparators)
	if cutAt < 0 {
		return "", errors.Errorf("could not find appropriate separators in address %q", a.raw.Address)
	}
	userPart := a.raw.Address[:cutAt]
	return userPart, nil
}

func (a *address) Tag() (string, error) {
	beginAt := strings.Index(a.raw.Address, tagSeparator)
	if beginAt < 0 {
		return "", nil
	}
	finishAt := strings.Index(a.raw.Address, domainSeparator)
	if finishAt < 0 {
		return "", nil
	}

	// defensive, not sure when this'll happen though.
	if finishAt-beginAt < 0 {
		return "", nil
	}

	tagPart := a.raw.Address[beginAt+1 : finishAt]
	return tagPart, nil
}
