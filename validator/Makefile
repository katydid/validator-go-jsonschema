# Copyright 2015 Walter Schulze
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

regenerate:
	(cd bnf && ./concat.sh)
	gocc -zip -o . -p github.com/katydid/validator-go-jsonschema/validator ./bnf/all.bnf
	rm ./util/litconv.go && rm ./util/rune.go && rmdir util
	gofmt -l -s -w ./parser/ ./errors/ ./lexer/ ./token/
	(cd funcs && make regenerate)
	(cd compose && make regenerate)

clean:
	rm *.txt || true

nuke: clean
	rm -rf ./errors/errors.go
	rm -rf ./lexer/acttab.go
	rm -rf ./lexer/lexer.go
	rm -rf ./lexer/transitiontable.go
	rm -rf ./parser/action.go
	rm -rf ./parser/actiontable.go
	rm -rf ./parser/gototable.go
	rm -rf ./parser/parser.go
	rm -rf ./parser/productionstable.go
	rm -rf ./token/token.go
	rm -rf ./util
	(cd funcs && make nuke)