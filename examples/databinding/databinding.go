// Copyright 2013 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

func main() {
	walk.AppendToWalkInit(func() {
		walk.FocusEffect, _ = walk.NewBorderGlowEffect(walk.RGB(0, 63, 255))
		walk.InteractionEffect, _ = walk.NewDropShadowEffect(walk.RGB(63, 63, 63))
		walk.ValidationErrorEffect, _ = walk.NewBorderGlowEffect(walk.RGB(255, 0, 0))
	})

	var mw *walk.MainWindow
	var outTE *walk.TextEdit

	animal := new(Animal)

	if _, err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk Data Binding Example",
		MinSize:  Size{300, 200},
		Layout:   VBox{},
		Children: []Widget{
			PushButton{
				Text: "Edit Animal",
				OnClicked: func() {
					if cmd, err := RunAnimalDialog(mw, animal); err != nil {
						log.Print(err)
					} else if cmd == walk.DlgCmdOK {
						outTE.SetText(fmt.Sprintf("%+v", animal))
					}
				},
			},
			Label{
				Text: "animal:",
			},
			// 编辑框
			TextEdit{
				AssignTo: &outTE,
				ReadOnly: true,
				Text:     fmt.Sprintf("%+v", animal),
			},
		},
	}.Run()); err != nil {
		log.Fatal(err)
	}
}

// 动物结构体
type Animal struct {
	Name          string
	ArrivalDate   time.Time
	SpeciesId     int
	Speed         int
	Sex           Sex
	Weight        float64
	PreferredFood string
	Domesticated  bool  // 是否是家养的
	Remarks       string
	Patience      time.Duration
}

func (a *Animal) PatienceField() *DurationField {
	return &DurationField{&a.Patience}
}

type Species struct {
	Id   int
	Name string
}

func KnownSpecies() []*Species {
	// 已知物种
	return []*Species{
		{1, "Dog"},
		{2, "Cat"},
		{3, "Bird"},
		{4, "Fish"},
		{5, "Elephant"},
	}
}

type DurationField struct {
	p *time.Duration
}

func (*DurationField) CanSet() bool       { return true }
func (f *DurationField) Get() interface{} { return f.p.String() }
func (f *DurationField) Set(v interface{}) error {
	x, err := time.ParseDuration(v.(string))
	if err == nil {
		*f.p = x
	}
	return err
}
func (f *DurationField) Zero() interface{} { return "" }

type Sex byte

const (
	SexMale Sex = 1 + iota
	SexFemale
	SexHermaphrodite
)

/*
 * 弹窗
 */
func RunAnimalDialog(owner walk.Form, animal *Animal) (int, error) {
	var dlg *walk.Dialog
	var db *walk.DataBinder
	var acceptPB, cancelPB *walk.PushButton

	return Dialog{
		AssignTo:      &dlg,
		// 标题
		Title:         Bind("'Animal Details' + (animal.Name == '' ? '' : ' - ' + animal.Name)"),
		// 确认按钮
		DefaultButton: &acceptPB,
		// 取消按钮
		CancelButton:  &cancelPB,
		// 数据绑定
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "animal",
			DataSource:     animal,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		MinSize: Size{300, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					// Name
					Label{
						Text: "Name:",
					},
					LineEdit{
						// 绑定的数据变量必须使用驼峰命名
						Text: Bind("Name"),
					},

					// Arrival Date
					Label{
						Text: "Arrival Date:",
					},
					// 日期编辑器
					DateEdit{
						Date: Bind("ArrivalDate"),
					},

					// Species
					Label{
						Text: "Species:",
					},
					ComboBox{
						// SelRequired 需要选择
						Value:         Bind("SpeciesId", SelRequired{}),
						BindingMember: "Id",
						DisplayMember: "Name",
						// 可选的模型
						Model:         KnownSpecies(),
					},

					// Speed
					Label{
						Text: "Speed:",
					},
					Slider{
						Value: Bind("Speed"),
					},

					// 单选组
					RadioButtonGroupBox{
						ColumnSpan: 2,
						Title:      "Sex",
						Layout:     HBox{},
						DataMember: "Sex",
						// 单选按钮组
						Buttons: []RadioButton{
							{Text: "Male", Value: SexMale},
							{Text: "Female", Value: SexFemale},
							{Text: "Hermaphrodite", Value: SexHermaphrodite},
						},
					},

					// Weight
					Label{
						Text: "Weight:",
					},
					NumberEdit{
						Value:    Bind("Weight", Range{0.01, 9999.99}),
						Suffix:   " kg",  // 后缀
						Decimals: 2,
					},

					// Preferred Food
					Label{
						Text: "Preferred Food:",
					},
					ComboBox{
						Editable: true,
						Value:    Bind("PreferredFood"),
						Model:    []string{"Fruit", "Grass", "Fish", "Meat"},
					},

					// Domesticated: 单选框
					Label{
						Text: "Domesticated:",
					},
					CheckBox{
						Checked: Bind("Domesticated"),
					},

					VSpacer{
						ColumnSpan: 2,
						Size:       8,
					},

					// Remarks：
					Label{
						ColumnSpan: 2,
						Text:       "Remarks:",
					},
					TextEdit{
						ColumnSpan: 2,
						MinSize:    Size{100, 50},
						Text:       Bind("Remarks"),
					},

					// Patience：
					Label{
						ColumnSpan: 2,
						Text:       "Patience:",
					},
					LineEdit{
						ColumnSpan: 2,
						Text:       Bind("PatienceField"),
					},
				},
			},

			// 按钮组
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							// db.Submit(): Widget -> Data
							// 可以使用 db.Reset() 从数据返回到 Widget
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}

							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
