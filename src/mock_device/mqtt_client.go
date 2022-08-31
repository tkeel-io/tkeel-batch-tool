// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"bytes"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"net/url"
	"text/template"
	"time"
)

//const (
//	host       = "http://preview.tkeel.io:30080/"
//	brokerUrl  = "45.120.216.31"
//	brokerPort = 31883
//	username   = "iotd-808c9e2e-0188-4eea-8106-d2c9247264a8"
//	password   = "ZDA3MmU1MGEtYWI1ZC0zNmYwLWE4M2MtMmE5ZDU3ZWRlN2Qx"
//	interval   = time.Second * 3
//)

type ClientOptions struct {
	Host        string
	DeviceID    string
	DeviceToken string
	Interval    int64
	Template    *template.Template
	Mode        string
}

// runDataSender use to to generate random numbers and send them into the device service as if a sensor
// was sending the data. Requires the Device Service along with Mongo, Core Data, and Metadata to be running
func RunDataSender(cfg ClientOptions) error {
	var mqttClientId = "IncomingDataPublisher"
	var qos = byte(0)
	var topic = fmt.Sprintf("v1/devices/me/%s", cfg.Mode)

	uri := &url.URL{
		Scheme: "tcp",
		Host:   cfg.Host,
		User:   url.UserPassword(cfg.DeviceID, cfg.DeviceToken),
	}

	client, err := createMqttClient(mqttClientId, uri)
	defer client.Disconnect(5000)
	if err != nil {
		fmt.Println(err)
	}

	var data = make(map[string]interface{})
	data["name"] = "MQTT test device"
	data["cmd"] = "randnum"
	data["method"] = "get"

	index := 1
	for {
		var buf bytes.Buffer //直接定义一个 Buffer 变量，而不用初始化
		err = cfg.Template.Execute(&buf, nil)
		if err != nil {
			return err
		}
		client.Publish(topic, qos, false, buf.Bytes())

		fmt.Printf("\n[%d]Send\n[%s]: \n %v", index, topic, buf.String())
		index++

		time.Sleep(time.Duration(cfg.Interval) * time.Millisecond)
	}
}

func createMqttClient(clientID string, uri *url.URL) (mqtt.Client, error) {
	fmt.Printf("Create MQTT client and connection: uri=%v clientID=%v ", uri.String(), clientID)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	opts.SetClientID(clientID)
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetConnectionLostHandler(func(client mqtt.Client, e error) {
		fmt.Printf("Connection lost : %v", e)
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			fmt.Printf("Reconnection failed : %v", e)
		} else {
			fmt.Printf("Reconnection sucessful : %v", e)
		}
	})

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return client, token.Error()
	}

	return client, nil
}
