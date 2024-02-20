//go:build js && wasm
// +build js,wasm

//go:generate cp $GOROOT/misc/wasm/wasm_exec.js .
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"
	"text/template"

	"github.com/euforic/templit"
)

type App struct {
	templateEditor js.Value
	jsonEditor     js.Value
	outputEditor   js.Value
}

func main() {
	app := &App{}
	app.init()

	// Block forever to keep the Go program running
	select {}
}

func (app *App) init() {
	app.initExports()
	app.initCodeMirror()
	app.initVimToggle()
	app.initRender()
}

func (app *App) initExports() {
	js.Global().Set("renderTemplate", js.FuncOf(app.renderTemplate))
}

func (app *App) initCodeMirror() {
	cm := js.Global().Get("CodeMirror")
	app.templateEditor = cm.New(js.Global().Get("document").Call("getElementById", "template"), js.ValueOf(map[string]interface{}{
		"mode":        "go",
		"theme":       "material-darker",
		"lineNumbers": true,
		"value":       "Hello, {{.Name}}",
	}))

	app.jsonEditor = cm.New(js.Global().Get("document").Call("getElementById", "json"), js.ValueOf(map[string]interface{}{
		"mode":        "javascript",
		"theme":       "material-darker",
		"lineNumbers": true,
		"value": `{
  "Name": "World"
}`,
	}))

	app.outputEditor = cm.New(js.Global().Get("document").Call("getElementById", "output"), js.ValueOf(map[string]interface{}{
		"lineNumbers": true,
		"readOnly":    true,
		"theme":       "material-darker",
		"value":       "Hello, World",
	}))
}

func (app *App) initVimToggle() {
	vimToggle := js.Global().Get("document").Call("getElementById", "vim-toggle")
	isVimEnabled := js.Global().Get("localStorage").Call("getItem", "vimEnabled").String() == "true"
	app.setKeyMap(isVimEnabled)
	vimToggle.Set("checked", isVimEnabled)
	vimToggle.Call("addEventListener", "change", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		app.setKeyMap(vimToggle.Get("checked").Bool())
		js.Global().Get("localStorage").Call("setItem", "vimEnabled", vimToggle.Get("checked").Bool())
		return nil
	}))
}

func (app *App) setKeyMap(vimEnabled bool) {
	keyMap := "default"
	if vimEnabled {
		keyMap = "vim"
	}
	app.templateEditor.Call("setOption", "keyMap", keyMap)
	app.jsonEditor.Call("setOption", "keyMap", keyMap)
	app.outputEditor.Call("setOption", "keyMap", keyMap)
}

func (app *App) initRender() {
	saveFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		app.saveCurrentData()
		return nil
	})
	renderFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go app.render()
		return nil
	})
	app.templateEditor.Call("on", "change", saveFunc)
	app.jsonEditor.Call("on", "change", saveFunc)
	app.templateEditor.Call("on", "change", renderFunc)
	app.jsonEditor.Call("on", "change", renderFunc)
	app.loadPersistedData()
}

func (app *App) render() {
	template := app.templateEditor.Call("getValue").String()
	json := app.jsonEditor.Call("getValue").String()
	result := app.renderTemplate(js.Null(), []js.Value{js.ValueOf(template), js.ValueOf(json)}).(string)
	app.outputEditor.Call("setValue", result)
}

func (app *App) renderTemplate(this js.Value, args []js.Value) interface{} {
	tmplStr := args[0].String()
	jsonStr := args[1].String()

	tmpl, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return fmt.Sprintf("Template error: %v", err)
	}

	tmpl.Funcs(templit.DefaultFuncMap)

	var data interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return fmt.Sprintf("JSON error: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Sprintf("Execution error: %v", err)
	}

	return buf.String()
}

func (app *App) saveCurrentData() {
	js.Global().Get("localStorage").Call("setItem", "templateData", app.templateEditor.Call("getValue").String())
	js.Global().Get("localStorage").Call("setItem", "jsonData", app.jsonEditor.Call("getValue").String())
}

func (app *App) loadPersistedData() {
	persistedTemplateData := js.Global().Get("localStorage").Call("getItem", "templateData").String()
	persistedJsonData := js.Global().Get("localStorage").Call("getItem", "jsonData").String()

	if persistedTemplateData != "<null>" {
		app.templateEditor.Call("setValue", persistedTemplateData)
	}

	if persistedJsonData != "<null>" {
		app.jsonEditor.Call("setValue", persistedJsonData)
	}

	if persistedTemplateData != "<null>" || persistedJsonData != "<null>" {
		app.render()
	}
}
