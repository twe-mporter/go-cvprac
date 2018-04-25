//
// Copyright (c) 2016-2017, Arista Networks, Inc. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//   * Redistributions of source code must retain the above copyright notice,
//   this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//   notice, this list of conditions and the following disclaimer in the
//   documentation and/or other materials provided with the distribution.
//
//   * Neither the name of Arista Networks nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL ARISTA NETWORKS
// BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
// BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
// OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN
// IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package cvpapi

import (
    "fmt"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// ConfigletList represents a list of configlets
type ConfigletList struct {
	Total int         `json:"total"`
	Data  []Configlet `json:"data"`

	ErrorResponse
}

// Custom structs for GetAppliedDevices()
type Device struct {
    IPAddress           string  `json:"ipAddress"`
    AppliedBy           string  `json:"appliedBy"`
    AppliedDate         int64   `json:"appliedDate"`
    ContainerName       string  `json:"containerName"`
    HostName            string  `json:"hostName"`
    TotalDevicesCount   int     `json:"totalDevicesCount"`

}

type DeviceList struct {
    Total   int         `json:"total"`
    Data    []Device    `json:"data"`

    ErrorResponse
}

// Configlet represents a Configlet
type Configlet struct {
	IsDefault            string `json:"isDefault"`
	DateTimeInLongFormat int64  `json:"dateTimeInLongFormat"`
	ContainerCount       int    `json:"containerCount"`
	NetElementCount      int    `json:"netElementCount"`
	IsAutoBuilder        string `json:"isAutoBuilder"`
	Reconciled           bool   `json:"reconciled"`
	FactoryID            int    `json:"factoryId"`
	Config               string `json:"config"`
	User                 string `json:"user"`
	Note                 string `json:"note"`
	Name                 string `json:"name"`
	Key                  string `json:"key"`
	ID                   int    `json:"id"`
	Type                 string `json:"type"`

	ErrorResponse
}



func (c Configlet) String() string {
	return c.Name
}

// ConfigletHistoryEntry represents a configlet history entry
type ConfigletHistoryEntry struct {
	ConfigletID                 string `json:"configletId"`
	OldUserID                   string `json:"oldUserId"`
	NewUserID                   string `json:"newUserId"`
	OldConfig                   string `json:"oldConfig"`
	NewConfig                   string `json:"newConfig"`
	OldDate                     string `json:"oldDate"`
	NewDate                     string `json:"newDate"`
	OldDateTimeInLongFormat     int64  `json:"oldDateTimeInLongFormat"`
	UpdatedDateTimeInLongFormat int64  `json:"updatedDateTimeInLongFormat"`
	FactoryID                   int    `json:"factoryId"`
	Key                         string `json:"key"`
	ID                          int    `json:"id"`
}

// ConfigletHistoryList represents a list of ConfigletHistoryEntry's
type ConfigletHistoryList struct {
	Total       int                     `json:"total"`
	HistoryList []ConfigletHistoryEntry `json:"configletHistory"`

	ErrorResponse
}

// ConfigletOpReturn represents the
type ConfigletOpReturn struct {
	Data Configlet `json:"data"`

	ErrorResponse
}

// UpdateResponse is a custom response type for synchronous UpdateConfiglet
type UpdateReturn struct {
    Data		    string		`json:"data"`
    TaskIDs 		[]string	`json:"taskIds"`
    ErrorResponse
}

// GetConfigletByName returns the configlet with the specified name
func (c CvpRestAPI) GetConfigletByName(name string) (*Configlet, error) {
	var info Configlet

	query := &url.Values{"name": {name}}

	resp, err := c.client.Get("/configlet/getConfigletByName.do", query)
	if err != nil {
		return nil, errors.Errorf("GetConfigletByName: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return nil, errors.Errorf("GetConfigletByName: %s", err)
	}

	if err := info.Error(); err != nil {
		// Entity does not exist
		if info.ErrorCode == "132801" {
			return nil, nil
		}
		return nil, errors.Errorf("GetConfigletByName: %s", err)
	}
	return &info, nil
}

// GetConfigletHistory returns the history for a configlet provided the key, and a range.
func (c CvpRestAPI) GetConfigletHistory(key string, start int,
	end int) (*ConfigletHistoryList, error) {
	var info ConfigletHistoryList

	//queryparam := url.Values{"name": {key},}

	query := &url.Values{
		"configletId": {key},
		"queryparam":  {""},
		"startIndex":  {strconv.Itoa(start)},
		"endIndex":    {strconv.Itoa(end)},
	}

	resp, err := c.client.Get("/configlet/getConfigletHistory.do", query)
	if err != nil {
		return nil, errors.Errorf("GetConfigletHistory: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return nil, errors.Errorf("GetConfigletHistory: %s", err)
	}

	if err := info.Error(); err != nil {
		return nil, errors.Errorf("GetConfigletHistory: %s", err)
	}

	return &info, nil
}

// GetAllConfigletHistory returns all the history for a given configlet
func (c CvpRestAPI) GetAllConfigletHistory(key string) (*ConfigletHistoryList, error) {
	return c.GetConfigletHistory(key, 0, 0)
}

// AddConfiglet creates/adds a configlet
func (c CvpRestAPI) AddConfiglet(name string, config string) (*Configlet, error) {
	var info ConfigletOpReturn

	data := map[string]string{
		"name":   name,
		"config": config,
	}

	resp, err := c.client.Post("/configlet/addConfiglet.do", nil, data)
	if err != nil {
		return nil, errors.Errorf("AddConfiglet: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return nil, errors.Errorf("AddConfiglet: %s", err)
	}

	if err := info.Error(); err != nil {
		return nil, errors.Errorf("AddConfiglet: %s", err)
	}

	return &info.Data, nil
}

// DeleteConfiglet deletes a configlet.
func (c CvpRestAPI) DeleteConfiglet(name string, key string) error {
	var info ErrorResponse

	data := []map[string]string{
		map[string]string{
			"name": name,
			"key":  key,
		},
	}
	resp, err := c.client.Post("/configlet/deleteConfiglet.do", nil, data)
	if err != nil {
		return errors.Errorf("DeleteConfiglet: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return errors.Errorf("DeleteConfiglet: %s", err)
	}

	if err := info.Error(); err != nil {
		return errors.Errorf("DeleteConfiglet: %s", err)
	}

	return nil
}

// UpdateConfiglet updates a configlet. This is the synchronous version that waits
// for any tasks to be created and returns those task IDs.
func (c CvpRestAPI) UpdateConfiglet(config, name, key string) ([]string, error) {

	// Response includes task information
    var info UpdateReturn

    data := struct {
        Config          string  `json:"config"`
        Key             string  `json:"key"`
        Name            string  `json:"name"`
        WaitForTaskIDs  bool    `json:"waitForTaskIds"`
    }{
        Config: config,
        Key:    key,
        Name:   name,
        WaitForTaskIDs: true,
    }

    jsonData, jsonErr := json.Marshal(data)

    if jsonErr != nil {
        return nil, errors.Errorf("UpdateConfiglet: %s", jsonErr)
    }

	resp, err := c.client.Post("/configlet/updateConfiglet.do", nil, jsonData)

	if err != nil {
		return nil, errors.Errorf("UpdateConfiglet: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return nil, errors.Errorf("UpdateConfiglet: %s", err)
	}

	if err := info.Error(); err != nil {
		return nil, errors.Errorf("UpdateConfiglet: %s", err)
	}

	return info.TaskIDs, nil
}


// UpdateConfigletAsync updates a configlet. This is the asynchronous version and does not 
// wait for any tasks to be created to return them. 
func (c CvpRestAPI) UpdateConfigletAsync(config string, name string, key string) error {
	var info ErrorResponse

	data := map[string]string{
		"config": config,
		"key":    key,
		"name":   name,
	}

	resp, err := c.client.Post("/configlet/updateConfiglet.do", nil, data)

	if err != nil {
		return errors.Errorf("UpdateConfiglet: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return errors.Errorf("UpdateConfiglet: %s", err)
	}

	if err := info.Error(); err != nil {
		return errors.Errorf("UpdateConfiglet: %s", err)
	}

	return nil
}


func (c CvpRestAPI) ValidateConfig(serialNumber, configlet string) {

    data := map[string]string{
        "config": configlet,
        "netElementId": serialNumber,
    }

    resp, err := c.client.Post("/configlet/validateConfig.do", nil, data)

    if err != nil {
        fmt.Println(err)
        //return errors.Errorf("ValidateConfig: %s", err)
    }


    fmt.Println(string(resp))



}


func (c CvpRestAPI) GetAppliedDevices(name string) (*DeviceList, error) {

    var info DeviceList

	query := &url.Values{
		"configletName": {name},
		"queryparam":  {""},
		"startIndex":  {strconv.Itoa(0)},
		"endIndex":    {strconv.Itoa(0)},
    }

	resp, err := c.client.Get("/configlet/getAppliedDevices.do", query)

    if err != nil {
        return nil, errors.Errorf("GetAppliedDevices: %s", err)
    }

    if err = json.Unmarshal(resp, &info); err != nil {
        return nil, errors.Errorf("GetAppliedDevices: %s", err)
    }

    if err := info.Error(); err != nil {
        return nil, errors.Errorf("GetAppliedDevices: %s", err)
    }

    return &info, nil

}



// SearchConfigletsWithRange search function for configlets.
func (c CvpRestAPI) SearchConfigletsWithRange(searchStr string, start int,
	end int) (*ConfigletList, error) {
	var info ConfigletList

	//queryparam := url.Values{"name": {key},}

	query := &url.Values{
		"queryparam": {searchStr},
		"startIndex": {strconv.Itoa(start)},
		"endIndex":   {strconv.Itoa(end)},
	}

	resp, err := c.client.Get("/configlet/searchConfiglets.do", query)
	if err != nil {
		return nil, errors.Errorf("SearchConfiglets: %s", err)
	}

	if err = json.Unmarshal(resp, &info); err != nil {
		return nil, errors.Errorf("SearchConfiglets: %s", err)
	}

	if err := info.Error(); err != nil {
		return nil, errors.Errorf("SearchConfiglets: %s", err)
	}

	return &info, nil
}

// SearchConfiglets search function for configlets.
func (c CvpRestAPI) SearchConfiglets(searchStr string) (*ConfigletList, error) {
	return c.SearchConfigletsWithRange(searchStr, 0, 0)
}
