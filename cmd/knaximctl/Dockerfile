# Copyright August 2020 Maxset Worldwide Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:1.13.4 AS backend

WORKDIR /go/src/git.maxset.io/web/knaxim

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd

RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo -v ./cmd/knaximctl

FROM alpine:3.12 AS production

COPY --from=backend /go/bin/knaximctl /bin/knaximctl
COPY container.config.json /etc/knaxim/conf.json
COPY resource /resource
COPY init.sh /bin/run.sh
RUN chmod 0755 /bin/run.sh
CMD ["/bin/run.sh"]
