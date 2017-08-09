package grabber

import (
	"github.com/cavaliercoder/grab"
	"github.com/samkrew/media_stats/logger"
	"time"
	"encoding/hex"
	"crypto"
	"github.com/imkira/go-libav/avformat"
	"github.com/imkira/go-libav/avutil"
	"github.com/dustin/go-humanize"
	"os"
	"github.com/samkrew/media_stats/queue"
	"github.com/tsuru/monsterqueue"
	"github.com/samkrew/media_stats/db"
	"fmt"
	"strconv"
)

var tmpDir = os.TempDir() + "/media_stats/"

type GrabberTask struct {
}

func getVideoStream(fmtCtx *avformat.Context) *avformat.Stream {
	for _, stream := range fmtCtx.Streams() {
		switch stream.CodecContext().CodecType() {
		case avutil.MediaTypeVideo:
			return stream
		}
	}
	return nil
}

func getStats(file string) (resolution string, bitrate string, err error) {
	fmtCtx, err := avformat.NewContextForInput()
	if err != nil {
		fmt.Errorf("Failed to open input context: %v", err)
		return
	}

	if err = fmtCtx.OpenInput(file, nil, nil); err != nil {
		fmt.Errorf("Failed to open input file: %v", err)
		return
	}

	stream := getVideoStream(fmtCtx)
	codecCtx := stream.CodecContext()

	bitrate = humanize.Bytes(uint64(codecCtx.BitRate())) + "/sec"
	resolution = strconv.Itoa(codecCtx.Width()) + "x" + strconv.Itoa(codecCtx.Height())

	fmtCtx.CloseInput()
	fmtCtx.Free()
	return
}

func (t *GrabberTask) Run(j monsterqueue.Job) {
	params := j.Parameters()
	url := params["url"].(string)
	hash := params["hash"].(string)

	client := grab.NewClient()
	req, _ := grab.NewRequest(tmpDir + j.ID(), url)

	md5 := make([]byte, 16)
	hex.Decode(md5, []byte(hash))
	req.SetChecksum(crypto.MD5.New(), md5, true)

	resp := client.Do(req)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

Loop:
	for {
		select {
		case <-ticker.C:
			db.UpdateStatus(j.ID(), fmt.Sprintf("Transferred %v / %v bytes (%.2f%%)",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress()))

		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		logger.L.Debugf("[Task #%v] Download failed: %v", j.ID(), err) // Not real error
		db.FinishStatsWithError(j.ID(), err.Error())
		if err == grab.ErrBadChecksum {
			j.Success("Success")
		} else {
			j.Error(err)
		}
		return
	}

	resolution, bitrate, err := getStats(resp.Filename)
	if err != nil {
		db.FinishStatsWithError(j.ID(), err.Error())
		logger.L.Debugf("[Task #%v] Get stats failed:", j.ID(), err)
		j.Error(err)
		return
	}
	db.FinishStats(j.ID(), resolution, bitrate)
	logger.L.Debugf("[Task #%v] File %s bitrate %s, resulution %s", j.ID(), resp.Filename, bitrate, resolution)
	os.Remove(resp.Filename)
	j.Success("Success")
}

func (t *GrabberTask) Name() string {
	return "Grabber"
}

func AddTask(url string, hash string) (string, error) {
	queue, err := queue.GetClientQueue()
	if err != nil {
		return "", err
	}

	job, err := queue.Instance.Enqueue("Grabber", monsterqueue.JobParams{
		"url":   url,
		"hash":  hash,
	})
	if err != nil {
		return "", err
	}


	return job.ID(), nil
}

func StartGrabber() {
	queue, err := queue.GetServerQueue()
	if err != nil {
		logger.L.Fatalf("Get queue error:", err)
		return
	}
	queue.Instance.RegisterTask(&GrabberTask{})
	queue.Run()
}

func init() {
	os.RemoveAll(tmpDir)
	os.MkdirAll(tmpDir, os.ModePerm);
}