// Copyright © 2015 Jason Smith <jasonrichardsmith@gmail.com>.
//
// Use of this source code is governed by the LGPL-3
// license that can be found in the LICENSE file.

/*

The consul package implements the Pool and Declare interfaces for the Herald package

The package takes one flag "consul-ttl" to set the time until the service expires from consul using a heartbeat.

This service must be registered with herald.

	consul.Register()

The package utilizies the Default client returned from consul api, but this can be overriden by edited

	ConsulConfig.Config.Address = "192.168.1.56"



*/
