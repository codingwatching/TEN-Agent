//
// Copyright © 2025 Agora
// This file is part of TEN Framework, an open source project.
// Licensed under the Apache License, Version 2.0, with certain conditions.
// Refer to the "LICENSE" file in the root directory for more information.
//

package default_extension_go

import (
	ten "ten_framework/ten_runtime"
)

type serverExtension struct {
	ten.DefaultExtension

	registerCount int
	isStopping    bool
}

func (ext *serverExtension) OnConfigure(tenEnv ten.TenEnv) {
	tenEnv.Log(ten.LogLevelDebug, "OnConfigure")

	ext.registerCount = 0
	ext.isStopping = false
	tenEnv.OnConfigureDone()
}

func (ext *serverExtension) OnCmd(tenEnv ten.TenEnv, cmd ten.Cmd) {
	name, _ := cmd.GetName()
	if name == "register" {
		ext.registerCount++

		cmdResult, _ := ten.NewCmdResult(ten.StatusCodeOk, cmd)
		tenEnv.ReturnResult(cmdResult, nil)
	} else if name == "unregister" {
		ext.registerCount--

		cmdResult, _ := ten.NewCmdResult(ten.StatusCodeOk, cmd)
		tenEnv.ReturnResult(cmdResult, nil)

		if ext.registerCount == 0 && ext.isStopping {
			tenEnv.OnStopDone()
		}
	} else if name == "test" {
		cmdResult, _ := ten.NewCmdResult(ten.StatusCodeOk, cmd)
		cmdResult.SetPropertyString("detail", "ok")
		tenEnv.ReturnResult(cmdResult, nil)
	} else {
		panic("unknown cmd name: " + name)
	}
}

func (ext *serverExtension) OnStop(tenEnv ten.TenEnv) {
	ext.isStopping = true

	if ext.registerCount == 0 {
		tenEnv.OnStopDone()
		return
	}

	tenEnv.Log(ten.LogLevelInfo, "server extension is stopped")
}

type clientExtension struct {
	ten.DefaultExtension
}

func (ext *clientExtension) OnStart(tenEnv ten.TenEnv) {
	tenEnv.Log(ten.LogLevelDebug, "OnStart")

	cmd, _ := ten.NewCmd("register")
	tenEnv.SendCmd(
		cmd,
		func(tenEnv ten.TenEnv, cmdResult ten.CmdResult, err error) {
			if err != nil {
				panic("Failed to send cmd: " + err.Error())
			}

			statusCode, _ := cmdResult.GetStatusCode()
			if statusCode != ten.StatusCodeOk {
				panic("Failed to register")
			}

			tenEnv.OnStartDone()
		},
	)
}

func (ext *clientExtension) OnCmd(tenEnv ten.TenEnv, cmd ten.Cmd) {
	// bypass the cmd to the next extension
	tenEnv.SendCmd(
		cmd,
		func(tenEnv ten.TenEnv, cmdResult ten.CmdResult, err error) {
			if err != nil {
				panic("Failed to send cmd: " + err.Error())
			}

			tenEnv.ReturnResult(cmdResult, nil)
		},
	)
}

func (ext *clientExtension) OnDeinit(tenEnv ten.TenEnv) {
	tenEnv.Log(ten.LogLevelDebug, "OnDeinit")

	cmd, _ := ten.NewCmd("unregister")
	tenEnv.SendCmd(
		cmd,
		func(tenEnv ten.TenEnv, cmdResult ten.CmdResult, err error) {
			if err != nil {
				panic("Failed to send cmd: " + err.Error())
			}

			statusCode, _ := cmdResult.GetStatusCode()
			if statusCode != ten.StatusCodeOk {
				panic("Failed to unregister")
			}

			tenEnv.OnDeinitDone()
		},
	)
}

func newAExtension(name string) ten.Extension {
	if name == "server" {
		return &serverExtension{}
	} else if name == "client" {
		return &clientExtension{}
	}

	return nil
}

func init() {
	// Register addon.
	err := ten.RegisterAddonAsExtension(
		"default_extension_go",
		ten.NewDefaultExtensionAddon(newAExtension),
	)
	if err != nil {
		panic("Failed to register addon.")
	}
}
