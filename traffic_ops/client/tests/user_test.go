/*
   Copyright 2015 Comcast Cable Communications Management, LLC

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package test

import (
	"net/http"
	"testing"

	"github.com/jheitz200/traffic_control/traffic_ops/client"
	"github.com/jheitz200/traffic_control/traffic_ops/client/fixtures"
)

func TestUsers(t *testing.T) {
	resp := fixtures.Users()
	server := validServer(resp)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a successful Traffic Ops request for Users")

	users, err := to.Users()
	if err != nil {
		Error(t, "Should be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should be able to make a request to Traffic Ops")
	}

	for _, u := range users {
		if u.FullName != "Bob Smith" {
			Error(t, "Should get back \"Bob Smith\" for \"FullName\", got %s", u.FullName)
		} else {
			Success(t, "Should get back \"Bob Smith\" for \"FullName\"")
		}

		if u.PublicSSHKey != "some-ssh-key" {
			Error(t, "Should get back \"some-ssh-key\" for \"PublicSSHKey\", got %s", u.PublicSSHKey)
		} else {
			Success(t, "Should get back \"some-ssh-key\" for \"PublicSSHKey\"")
		}

		if u.Role != "3" {
			Error(t, "Should get back \"3\" for \"Role\", got %s", u.Role)
		} else {
			Success(t, "Should get back \"3\" for \"Role\"")
		}

		if u.Email != "bobsmith@email.com" {
			Error(t, "Should get back \"bobsmith@email.com\" for \"Email\", got %s", u.Email)
		} else {
			Success(t, "Should get back \"bobsmith@email.com\" for \"Email\"")
		}

		if u.Username != "bsmith" {
			Error(t, "Should get back \"bsmith\" for \"Username\", got %s", u.Username)
		} else {
			Success(t, "Should get back \"bsmith\" for \"Username\"")
		}
	}
}

func TestUsersUnauthorized(t *testing.T) {
	server := invalidServer(http.StatusUnauthorized)
	defer server.Close()

	var httpClient http.Client
	to := client.Session{
		URL:       server.URL,
		UserAgent: &httpClient,
	}

	Context(t, "Given the need to test a failed Traffic Ops request for Users")

	_, err := to.Users()
	if err == nil {
		Error(t, "Should not be able to make a request to Traffic Ops")
	} else {
		Success(t, "Should not be able to make a request to Traffic Ops")
	}
}
