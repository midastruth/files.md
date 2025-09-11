//go:build js && wasm

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"syscall/js"
	_ "time/tzdata" // for was env we need timezone database

	"zakirullin/stuffbot/internal"
	"zakirullin/stuffbot/internal/db"
	"zakirullin/stuffbot/internal/fs"
	"zakirullin/stuffbot/internal/userconfig"
	"zakirullin/stuffbot/pkg/tg"
)

var (
	reply func(u internal.Update)
	chat  *tg.FakeTG
)

type Update struct {
	Message string
	Command *tg.Cmd
}

type Response struct {
	Messages []tg.Message
}

func Reply(_ js.Value, args []js.Value) any {
	logToJS("Wasm: called reply")
	upd := tg.NewUpd(-1, args[0].String())
	go reply(upd)

	return nil
}

func ReplyCmd(_ js.Value, args []js.Value) any {
	fmt.Println("CMD", args[0].String())
	var cmd *tg.Cmd
	err := json.Unmarshal([]byte(args[0].String()), &cmd)
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	upd := tg.NewUpdCmd(-1, *cmd)
	go reply(upd)

	return nil
}

func main() {
	js.Global().Set("wasmReply", js.FuncOf(Reply))
	js.Global().Set("wasmReplyCmd", js.FuncOf(ReplyCmd))
	fs.Exists = exists
	fs.ReadFile = readFile
	fs.WriteFile = writeFile
	fs.ReadDir = readDir
	initBot()
	js.Global().Call("wasmReady")

	select {}
}

func callAsync(funcName string, callback func(js.Value, error), args ...any) {
	promise := js.Global().Call(funcName, args)

	var successFunc, errorFunc js.Func

	successFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer successFunc.Release()
		defer errorFunc.Release()
		callback(args[0], nil)
		return nil
	})

	errorFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer successFunc.Release()
		defer errorFunc.Release()

		var err error
		if len(args) > 0 && !args[0].IsUndefined() {
			errorMsg := args[0].Get("message").String()
			if strings.Contains(errorMsg, "could not be found") {
				err = os.ErrNotExist
			} else {
				err = errors.New(errorMsg)
			}
		} else {
			err = errors.New("unknown error")
		}

		callback(js.Undefined(), err)
		return nil
	})

	promise.Call("then", successFunc).Call("catch", errorFunc)
}

func sendDueResponsesToJS() {
	var r Response
	r.Messages = chat.Messages
	if chat.EditedMessages != nil {
		r.Messages = append(r.Messages, chat.EditedMessages...)
	}

	chat.Messages = nil
	chat.EditedMessages = nil

	response, err := json.Marshal(r)
	if err != nil {
		// TODO handle err
	}
	logToJS("Sending from WASM to JS:", string(response))
	sendToJS(string(response))
}

func sendToJS(vals ...any) {
	var paths []any
	for _, path := range writtenFiles {
		paths = append(paths, "/"+path)
	}
	writtenFiles = nil
	js.Global().Call("receive", paths)
}

func logToJS(vals ...any) {
	js.Global().Call("logWasm", vals...)
}

func initBot() {
	//opts := &tint.Options{
	//	Level: slog.LevelDebug,
	//}
	//logger := slog.New(tint.NewHandler(os.Stderr, opts))
	//slog.SetDefault(logger)

	// For GUI app we don't have required .env params
	//_ = godotenv.Load()
	//err := config.LoadGUIConfig()
	//if err != nil {
	//	panic(fmt.Sprintf("Error loading cfg: %s\n", err))
	//}

	// TODO move to embed
	//err = i18n.LoadLangFile("i18n/ru.json")
	//if err != nil {
	//	panic(fmt.Sprintf("Error loading i18n: %s\n", err))
	//}

	reply = func(u internal.Update) {
		defer func() {
			err := recover()
			if err != nil {
				debug.PrintStack()
				slog.Error("Bot panic", "err", err)
				logToJS("Wasm panic:", err)
			}
		}()

		userID := u.UserID()

		userPath := ""
		userFS, err := fs.NewFS(userPath, NewJSFS())
		if err != nil {
			fmt.Printf("Bot error: can't create fs: %v", err)
		}

		confFilename := "config.json"
		userconf := userconfig.NewConfig(userFS, userID, confFilename)
		err = userconf.CreateDefaultIfNotExists()
		if err != nil {
			fmt.Printf("Bot error: can't create default user config")
		}

		if chat == nil {
			chat = tg.NewFakeTG()
		}
		bot := internal.NewBot(userID, chat, userFS, db.NewDB(userID), userconf)
		if err := bot.Reply(u); err != nil {
			logToJS("Wasm error:", err)
			fmt.Printf("Bot error: %v", err)
		}

		sendDueResponsesToJS()
	}
}

//func send(update Update) Response {
//	if update.Command != nil {
//		_ = reply(tg.NewUpdCmd(1, *update.Command))
//	} else {
//		_ = reply(tg.NewUpd(1, update.Message))
//	}
//
//	var r Response
//	r.Messages = chat.Messages
//	if chat.EditedMessages != nil {
//		r.Messages = append(r.Messages, chat.EditedMessages...)
//	}
//
//	chat.Messages = nil
//	chat.EditedMessages = nil
//
//	return r
//}

func newUpdate(message string, cmd *tg.Cmd) Update {
	return Update{
		Message: message,
		Command: cmd,
	}
}
