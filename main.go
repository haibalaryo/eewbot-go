package main

import (
	"fmt"
	"github.com/nexryai/eewbot-go/notify"
	"github.com/nexryai/eewbot-go/quake"
	"github.com/nexryai/eewbot-go/xvfb"
	"os"
	"strconv"
	"sync"
)

func main() {
	targetHost := os.Getenv("INSTANCE_HOST")
	if targetHost == "" {
		fmt.Println("Error: INSTANCE_HOST is not set.")
		return
	}

	eventType := os.Getenv("KEVI_EVENT_TYPE")

	// --- 共通処理: スクショを撮ってドライブに上げる ---
	// EEW受信時 または 地震情報(Equake)受信時 に実行
	if eventType == "EEW_RECEIVED" || eventType == "EQUAKE" || os.Getenv("EEW_DEBUGMODE") == "1" {
		
		imageData, err := xvfb.TakeScreenshotOfXvfb()
		if err != nil {
			fmt.Println("Screenshot Error:", err)
			os.Exit(1)
		}

		data := notify.MisskeyDriveUploadForm{
			InstanceHost: os.Getenv("MISSKEY_HOST"),
			Token:        os.Getenv("MISSKEY_TOKEN"),
			Data:         *imageData,
		}

		driveApiResp, err := notify.UploadToMisskeyDrive(data)
		if err != nil {
			fmt.Println("Upload Error:", err)
			os.Exit(1)
		}
		fmt.Println("File Uploaded:", driveApiResp.FileID)

		// --- 投稿内容の作成 ---
		var text string
		
		if eventType == "EEW_RECEIVED" {
			// === 緊急地震速報 (EEW) ===
			reportNumInt, _ := strconv.ParseInt(os.Getenv("EEW_COUNT"), 10, 16)
			place := os.Getenv("EEW_PLACE")
			intensity := os.Getenv("EEW_INTENSITY")

			if quake.IsEmergency(intensity) {
				text = fmt.Sprintf("<center>$[bg.color=ff0000 ⚠️$[fg.color=fff **緊急地震速報(警報)** 第%d報]⚠️]</center>\n\n震源: **%s**\n$[fg 最大震度: 震度**%s**]",
					reportNumInt, place, intensity)
			} else {
				text = fmt.Sprintf("<center>**⚠️緊急地震速報(EEW) 第%d報⚠️**</center>\n\n震源: **%s**\n最大震度: 震度**%s**",
					reportNumInt, place, intensity)
			}

		} else if eventType == "EQUAKE" {
			// === 地震情報 (結果報告) ===
			// ※Config.jsonで設定した環境変数 EQ_PLACE, EQ_INTENSITY を受け取る
			place := os.Getenv("EQ_PLACE")
			intensity := os.Getenv("EQ_INTENSITY")

			text = fmt.Sprintf("<center>**【地震情報】**</center>\n\n震源: **%s**\n最大震度: 震度**%s**\n\n(詳細情報は気象庁HP等を確認してください)",
				place, intensity)
		
		} else if os.Getenv("EEW_DEBUGMODE") == "1" {
			// === デバッグモード ===
			text = "Botは正常に起動しています。(Debug Mode)"
		}

		// --- Misskeyに投稿 ---
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			note := notify.MisskeyNote{
				InstanceHost: targetHost,
				Token:        os.Getenv("MISSKEY_TOKEN"),
				Text:         text,
				LocalOnly:    false,
				Visibility:   "public",
				FileIds:      []string{driveApiResp.FileID},
			}

			err = notify.PostToMisskey(note)
			if err != nil {
				fmt.Println("Post Error:", err)
			} else {
				fmt.Println("Note Posted Successfully")
			}
		}()
		wg.Wait()
	}
}