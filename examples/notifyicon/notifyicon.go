// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
)

import (
	"github.com/lxn/walk"
)

func main() {
	// We need either a walk.MainWindow or a walk.Dialog for their message loop.
	// We will not make it visible in this example, though.
	mw, err := walk.NewMainWindow()
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个 icon
	// We load our icon from a file.
	icon, err := walk.Resources.Icon("../img/stop.ico")
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个 notify icon ，并确保在退出时清除它
	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	// 设置图标和提示文本
	// Set the icon and a tool tip text.
	if err := ni.SetIcon(icon); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}

	// 当鼠标左键按下或弹起
	// When the left mouse button is pressed, bring up our balloon.
	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button != walk.LeftButton {
			return
		}

		if err := ni.ShowCustom(
			"Walk NotifyIcon Example",
			"There are multiple ShowX methods sporting different icons.",
			icon); err != nil {

			log.Fatal(err)
		}
	})

	// 放置一个 exit 操作到 menu
	// We put an exit action into the context menu.
	// Action: https://github.com/lxn/walk/blob/master/action.go#L23
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err)
	}
	// 添加触发事件
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err)
	}

	// notify icon初始化时候是隐藏的，我们让它显示
	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Now that the icon is visible, we can bring up an info balloon.
	if err := ni.ShowInfo("Walk NotifyIcon Example", "Click the icon to show again."); err != nil {
		log.Fatal(err)
	}

	// Run the message loop.
	mw.Run()
}
