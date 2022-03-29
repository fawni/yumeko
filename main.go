package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/logrusorgru/aurora/v3"
	"github.com/muesli/termenv"
	"github.com/x6r/sip/internal/common"
	"github.com/x6r/sip/internal/mpv"
	"github.com/x6r/sip/internal/piped"
)

const videoPerPage = 20

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		common.Fatal("User interuptted.")
	}()
}

func main() {
	flag.Parse()
	baseURL, feedURL := piped.GetInstanceFromConfig()
	searchURL := baseURL + "/search?q="
	streamURL := baseURL + "/streams/"

	termenv.ClearScreen()

	var (
		body []byte
		err  error
		feed []common.Video
	)

	if flag.Arg(0) != "" {
		query := fmt.Sprintf("%s%s&filter=videos", searchURL, url.QueryEscape(flag.Arg(0)))
		_, body, err = piped.Get(query)
		if err != nil {
			common.Fatal(err)
		}
		var sr common.Search
		if err := json.Unmarshal(body, &sr); err != nil {
			common.Fatal(err)
		}

		feed = sr.Items
	} else {
		_, body, err = piped.Get(feedURL)
		if err != nil {
			common.Fatal(err)
		}

		if err := json.Unmarshal(body, &feed); err != nil {
			var video common.Video
			if err := json.Unmarshal(body, &video); err != nil {
				common.Fatal(err)
			}
			common.Fatal(video.Error)
		}
	}

	var (
		page int
		end  bool
	)

page:
	for i := 0; i < videoPerPage; i++ {
		if (i + (page * videoPerPage)) == len(feed) {
			end = true
			break
		}
		video := feed[i+(page*videoPerPage)]
		d, _ := time.ParseDuration(strconv.Itoa(video.Duration) + "s")
		duration := common.FmtDuration(d)
		title := video.Title
		if len(title) > 70 {
			title = fmt.Sprintf("%.70s...", title)
		}
		fmt.Printf("%2d. %s - %s (%s) | %s\n", i+(page*videoPerPage)+1, aurora.Green(video.Uploader), title, aurora.Cyan(duration), aurora.Magenta(humanize.Time(time.UnixMilli(video.Uploaded))))
	}

	var keys string
	if page == 0 {
		keys = "[Q] Exit\n[N] Next page"
	} else if end {
		keys = "[Q] Exit\n[P] Previous page"
	} else {
		keys = "[Q] Exit\n[P] Previous page\n[N] Next page"
	}
	fmt.Println(aurora.Bold(aurora.Green(keys)))

choice:
	var input string

	fmt.Printf("%s %s › ", aurora.Green("?"), aurora.Bold(aurora.Cyan("Enter choice")))
	fmt.Scanln(&input)

	switch input {
	case "q", "Q":
		os.Exit(0)
	case "n", "N":
		if end {
			fmt.Println(aurora.Red("You are on the last page"))
			goto choice
		}
		page++
		termenv.ClearScreen()
		goto page
	case "p", "P":
		if page == 0 {
			fmt.Println(aurora.Red("You are on the first page"))
			goto choice
		}
		end = false
		page--
		termenv.ClearScreen()
		goto page
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice-page > len(feed) || choice == 0 {
		fmt.Println(aurora.Red("Invalid choice"))
		goto choice
	}

	video := feed[choice-1]
	videoURL := "https://youtube.com" + video.URL
	termenv.ClearScreen()

	mpv.Play(videoURL, fmt.Sprintf("%s - %s", video.Uploader, video.Title))

	_, body, err = piped.Get(streamURL + strings.TrimPrefix(video.URL, "/watch?v="))
	if err != nil {
		common.Fatal(err)
	}
	if err := json.Unmarshal(body, &video.Details); err != nil {
		common.Fatal(err)
	}

	_, body, err = piped.Get(baseURL + video.Details.ChannelURL)
	if err != nil {
		common.Fatal(err)
	}
	if err := json.Unmarshal(body, &video.Channel); err != nil {
		common.Fatal(err)
	}

	var desc string
	if len(video.Details.Description) > 300 {
		desc = fmt.Sprintf("%.400s...", video.Details.Description)
	} else {
		desc = video.Details.Description
	}
	videoTime := time.UnixMilli(video.Uploaded)
	d, err := time.ParseDuration(strconv.Itoa(video.Duration) + "s")
	if err != nil {
		common.Fatal(err)
	}
	duration := common.FmtDuration(d)

	fmt.Printf("%s: %s (%s)\n", aurora.Bold(aurora.Magenta("Title")), video.Title, aurora.Italic(duration))
	fmt.Printf("%s: %s (%s)\n", aurora.Bold(aurora.Magenta("Uploader")), video.Uploader, aurora.Italic(common.SI(video.Channel.Subscibers, "Subscribers")))
	fmt.Printf("%s: %s (%s)\n", aurora.Bold(aurora.Magenta("Uploaded")), common.FmtTime(videoTime), aurora.Italic(humanize.Time(videoTime)))
	fmt.Printf("%s: %s\n", aurora.Bold(aurora.Magenta("Viewes")), humanize.Comma(video.Details.Views))
	fmt.Printf("%s: %s/%s\n", aurora.Bold(aurora.Magenta("Likes/Dislikes")), aurora.Green(humanize.Comma(video.Details.Likes)), aurora.Red(humanize.Comma(video.Details.Dislikes)))
	fmt.Printf("%s: %s\n", aurora.Bold(aurora.Magenta("Link")), aurora.Cyan(videoURL))
	fmt.Printf("%s:\n-----\n%s\n-----\n", aurora.Bold(aurora.Magenta("Description")), desc)
	fmt.Println(aurora.Bold(aurora.Green("[B] Go back\n[Q] Exit")))

	fmt.Printf("%s %s › ", aurora.Green("?"), aurora.Bold(aurora.Cyan("Enter choice")))
	fmt.Scanln(&input)

	switch input {
	case "q", "Q":
		os.Exit(0)
	case "b", "B":
		termenv.ClearScreen()
		goto page
	}
}
