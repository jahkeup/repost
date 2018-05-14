#!/usr/bin/env bash
#
# This script generates convenience functions for use in a template
# func map.
#
# usage: ./gentemplateapi.sh Name Tag User Domain
#

echo -e "// This is a generated file, do not edit.\npackage emailfuncs\n" > ./funcs.go

function writefn() {
    local name=$1
    cat >> ./funcs.go <<EOF

// $name returns the $name from the parsed Email address provided.
func $name(str string) (string, error) {
	a, err := Parse(str)
	if err != nil {
		return "", err
	}
	return a.$name()
}

EOF
}

for addrfn in $@; do
    writefn $addrfn
done
