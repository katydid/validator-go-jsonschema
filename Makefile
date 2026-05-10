# Copyright 2013 Walter Schulze
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: nuke regenerate gofmt build test

all: nuke regenerate build test

test:
	go clean -testcache
	TESTSUITE=MUST go test -v ./...

suite_draft4:
	go clean -testcache
	go test -run=TestSuiteDraft4 -v ./jsonschema

suite_202012:
	go clean -testcache
	go test -run=TestSuite202012 -v ./jsonschema

debug:
	go clean -testcache
	go test -run=TestDebug -v ./jsonschema

build:
	go build ./...

install:
	go install ./...

bench:
	TESTSUITE=MUST go test -test.v -test.run=XXX -test.bench=. ./...

regenerate:
	(cd jsonschema/funcs/email && make regenerate)

clean:
	go clean ./...

nuke: clean
	go clean -i ./...

gofmt:
	gofmt -l -s -w .

diff:
	git diff --exit-code .
