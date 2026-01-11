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
    // 空っぽだったらエラーで落とす
    if targetHost == "" {
        fmt.Println("Error: INSTANCE_HOST is not set.")
        return
    }

	if os.Getenv("KEVI_EVENT_TYPE") == "EEW_RECEIVED" {
		imageData, err := xvfb.TakeScreenshotOfXvfb()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data := notify.MisskeyDriveUploadForm{
			InstanceHost: os.Getenv("MISSKEY_HOST"),
			Token:        os.Getenv("MISSKEY_TOKEN"),
			Data:         *imageData,
		}

		driveApiResp, err := notify.UploadToMisskeyDrive(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(driveApiResp.FileID)

		// 通知の準備
		reportNumInt, err := strconv.ParseInt(os.Getenv("EEW_COUNT"), 10, 16)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		q := quake.Equake{
			ReportNum:     uint16(reportNumInt),
			DispIntensity: os.Getenv("EEW_INTENSITY"),
			Place:         os.Getenv("EEW_PLACE"),
		}

		var wg sync.WaitGroup

		// Misskeyに通知
		wg.Add(1)
		go func() {
			defer wg.Done()
			var text string
			if quake.IsEmergency(q.DispIntensity) {
				text = fmt.Sprintf("<center>$[bg.color=ff0000 ⚠️$[fg.color=fff **緊急地震速報(警報)** 第%d報]⚠️]</center>\n\n震源: **%s**\n$[fg 最大震度: 震度**%s**]",
					q.ReportNum,
					q.Place,
					q.DispIntensity)
			} else {
				text = fmt.Sprintf("<center>**⚠️緊急地震速報(EEW) 第%d報⚠️**</center>\n\n震源: **%s**\n最大震度: 震度**%s**",
					q.ReportNum,
					q.Place,
					q.DispIntensity)
			}

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
				fmt.Println(err)
			}
		}()

		// // Discordに通知
		// wg.Add(1)
		// go func() {
		// 	defer wg.Done()
		// 	var color int
		// 	var discordEmbedTitle string
		// 	var discordEmbedMessage string

		// 	if quake.IsEmergency(q.DispIntensity) {
		// 		color = 0xD50000

		// 		// 本文
		// 		discordEmbedTitle = fmt.Sprintf("最大震度%s 緊急地震速報第%d報 強い揺れに警戒",
		// 			q.DispIntensity,
		// 			q.ReportNum)

		// 		discordEmbedMessage = fmt.Sprintf("強い揺れに警戒。%sを震源とする最大震度%sの地震が発生しました。",
		// 			q.Place,
		// 			q.DispIntensity)

		// 	} else {
		// 		if os.Getenv("DISCORD_SILENT") == "1" {
		// 			return
		// 		}

		// 		color = 0xffc000

		// 		// 本文
		// 		discordEmbedTitle = fmt.Sprintf("最大震度%s 緊急地震速報第%d報",
		// 			q.DispIntensity,
		// 			q.ReportNum)

		// 		discordEmbedMessage = fmt.Sprintf("%sを震源とする最大震度%sの地震が発生しました。",
		// 			q.Place,
		// 			q.DispIntensity)
		// 	}

		// 	var discordNotify notify.DiscordHook
		// 	discordNotify.Username = "EEW Bot"
		// 	discordNotify.Content = fmt.Sprintf("最大震度%s %s震源 EEW第%d報",
		// 		q.DispIntensity,
		// 		q.Place,
		// 		q.ReportNum)
		// 	discordNotify.Embeds = []notify.DiscordEmbed{
		// 		notify.DiscordEmbed{
		// 			Title:  discordEmbedTitle,
		// 			Desc:   discordEmbedMessage,
		// 			Color:  color,
		// 			Author: notify.DiscordAuthor{Name: "eewbot-go"},
		// 			Image:  notify.DiscordImg{URL: driveApiResp.Url},
		// 			Fields: []notify.DiscordField{
		// 				notify.DiscordField{Name: "震源", Value: q.Place, Inline: true},
		// 				notify.DiscordField{Name: "最大震度", Value: q.DispIntensity, Inline: true},
		// 			},
		// 		},
		// 	}

		// 	err = notify.NotifyToDiscord(discordNotify)
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}
		// }()

		wg.Wait()

	} else if os.Getenv("EEW_DEBUGMODE") == "1" {
		imageData, err := xvfb.TakeScreenshotOfXvfb()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data := notify.MisskeyDriveUploadForm{
			InstanceHost: targetHost,
			Token:        os.Getenv("MISSKEY_TOKEN"),
			Data:         *imageData,
		}

		driveApiResp, err := notify.UploadToMisskeyDrive(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(driveApiResp.FileID)

		text := "Botは起動しました。"

		note := notify.MisskeyNote{
			InstanceHost: targetHost,
			Token:        os.Getenv("MISSKEY_TOKEN"),
			Text:         text,
			LocalOnly:    true,
			Visibility:   "home",
			FileIds:      []string{driveApiResp.FileID},
		}

		err = notify.PostToMisskey(note)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
