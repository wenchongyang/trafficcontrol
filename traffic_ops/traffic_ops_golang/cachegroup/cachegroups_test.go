package cachegroup

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/apache/trafficcontrol/lib/go-tc"
	"github.com/apache/trafficcontrol/lib/go-tc/v13"
	"github.com/apache/trafficcontrol/lib/go-util"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/api"
	"github.com/apache/trafficcontrol/traffic_ops/traffic_ops_golang/test"
	"github.com/jmoiron/sqlx"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func getTestCacheGroups() []v13.CacheGroup {
	cgs := []v13.CacheGroup{}
	testCG1 := v13.CacheGroup{
		ID:                          1,
		Name:                        "cachegroup1",
		ShortName:                   "cg1",
		Latitude:                    38.7,
		Longitude:                   90.7,
		ParentCachegroupID:          2,
		SecondaryParentCachegroupID: 2,
		Type:        "EDGE_LOC",
		TypeID:      6,
		LastUpdated: tc.TimeNoMod{Time: time.Now()},
	}
	cgs = append(cgs, testCG1)

	testCG2 := v13.CacheGroup{
		ID:                          1,
		Name:                        "parentCacheGroup",
		ShortName:                   "pg1",
		Latitude:                    38.7,
		Longitude:                   90.7,
		ParentCachegroupID:          1,
		SecondaryParentCachegroupID: 1,
		Type:        "MID_LOC",
		TypeID:      7,
		LastUpdated: tc.TimeNoMod{Time: time.Now()},
	}
	cgs = append(cgs, testCG2)

	return cgs
}

func TestReadCacheGroups(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	testCGs := getTestCacheGroups()
	cols := test.ColsFromStructByTag("db", v13.CacheGroup{})
	rows := sqlmock.NewRows(cols)

	for _, ts := range testCGs {
		rows = rows.AddRow(
			ts.ID,
			ts.Name,
			ts.ShortName,
			ts.Latitude,
			ts.Longitude,
			ts.ParentName,
			ts.ParentCachegroupID,
			ts.SecondaryParentName,
			ts.SecondaryParentCachegroupID,
			ts.FallbackToClosest,
			ts.Type,
			ts.TypeID,
			ts.LastUpdated,
		)
	}
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	mock.ExpectCommit()
	v := map[string]string{"id": "1"}

	reqInfo := api.APIInfo{Tx: db.MustBegin(), CommitTx: util.BoolPtr(false)}
	cachegroups, errs, _ := GetTypeSingleton()(&reqInfo).Read(v)
	if len(errs) > 0 {
		t.Errorf("cdn.Read expected: no errors, actual: %v", errs)
	}

	if len(cachegroups) != 2 {
		t.Errorf("cdn.Read expected: len(cachegroups) == 2, actual: %v", len(cachegroups))
	}
}

func TestFuncs(t *testing.T) {
	if strings.Index(selectQuery(), "SELECT") != 0 {
		t.Errorf("expected selectQuery to start with SELECT")
	}
	if strings.Index(insertQuery(), "INSERT") != 0 {
		t.Errorf("expected insertQuery to start with INSERT")
	}
	if strings.Index(updateQuery(), "UPDATE") != 0 {
		t.Errorf("expected updateQuery to start with UPDATE")
	}
	if strings.Index(deleteQuery(), "DELETE") != 0 {
		t.Errorf("expected deleteQuery to start with DELETE")
	}
}

func TestInterfaces(t *testing.T) {
	var i interface{}
	i = &TOCacheGroup{}

	if _, ok := i.(api.Creator); !ok {
		t.Errorf("cachegroup must be creator")
	}
	if _, ok := i.(api.Reader); !ok {
		t.Errorf("cachegroup must be reader")
	}
	if _, ok := i.(api.Updater); !ok {
		t.Errorf("cachegroup must be updater")
	}
	if _, ok := i.(api.Deleter); !ok {
		t.Errorf("cachegroup must be deleter")
	}
	if _, ok := i.(api.Identifier); !ok {
		t.Errorf("cachegroup must be Identifier")
	}
}

func TestValidate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	defer db.Close()

	rows := sqlmock.NewRows([]string{"name", "use_in_table"})
	rows.AddRow("EDGE_LOC", "cachegroup")

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	tx := db.MustBegin()

	reqInfo := api.APIInfo{Tx: tx, CommitTx: util.BoolPtr(false)}

	// invalid name, shortname, loattude, and longitude
	id := 1
	nm := "not!a!valid!cachegroup"
	sn := "not!a!valid!shortname"
	la := -190.0
	lo := -190.0
	ty := "EDGE_LOC"
	ti := 6
	lu := tc.TimeNoMod{Time: time.Now()}
	c := TOCacheGroup{ReqInfo: &reqInfo, CacheGroupNullable: v13.CacheGroupNullable{ID: &id,
		Name:        &nm,
		ShortName:   &sn,
		Latitude:    &la,
		Longitude:   &lo,
		Type:        &ty,
		TypeID:      &ti,
		LastUpdated: &lu,
	}}
	errs := util.JoinErrsStr(test.SortErrors(test.SplitErrors(c.Validate())))

	expectedErrs := util.JoinErrsStr([]error{
		errors.New(`'latitude' Must be a floating point number within the range +-90`),
		errors.New(`'longitude' Must be a floating point number within the range +-180`),
		errors.New(`'name' invalid characters found - Use alphanumeric . or - or _ .`),
		errors.New(`'shortName' invalid characters found - Use alphanumeric . or - or _ .`),
	})

	if !reflect.DeepEqual(expectedErrs, errs) {
		t.Errorf("expected %s, got %s", expectedErrs, errs)
	}

	rows.AddRow("EDGE_LOC", "cachegroup")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	//  valid name, shortName latitude, longitude
	nm = "This.is.2.a-Valid---Cachegroup."
	sn = `awesome-cachegroup`
	la = 90.0
	lo = 90.0
	c = TOCacheGroup{ReqInfo: &reqInfo, CacheGroupNullable: v13.CacheGroupNullable{ID: &id,
		Name:        &nm,
		ShortName:   &sn,
		Latitude:    &la,
		Longitude:   &lo,
		Type:        &ty,
		TypeID:      &ti,
		LastUpdated: &lu,
	}}
	err = c.Validate()
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}
}
