package main

import (
	"log"
	"time"

	"github.com/pion/interceptor"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"

	"pion-learning/signal"
)

func main() {
	m := &webrtc.MediaEngine{}

	if err := m.RegisterCodec(webrtc.RTPCodecParameters{
		RTPCodecCapability: webrtc.RTPCodecCapability{
			MimeType:     "video/VP8",
			ClockRate:    90000,
			Channels:     0,
			SDPFmtpLine:  "",
			RTCPFeedback: nil,
		},
		PayloadType: 96,
	}, webrtc.RTPCodecTypeVideo); err != nil {
		log.Panic(err)
	}

	i := &interceptor.Registry{}

	if err := webrtc.RegisterDefaultInterceptors(m, i); err != nil {
		log.Panic(err)
	}

	api := webrtc.NewAPI(webrtc.WithMediaEngine(m), webrtc.WithInterceptorRegistry(i))

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := api.NewPeerConnection(config)
	if err != nil {
		log.Panic(err)
	}

	outputTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{
		MimeType: "video/vp8",
	}, "video", "stream-id")
	if err != nil {
		log.Panic(err)
	}

	rtpSender, err := pc.AddTrack(outputTrack)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		// RTP Confrol Protocol
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	offer := webrtc.SessionDescription{}
	signal.Decode("eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz0tIDM5OTg1OTM3NzY1NDI4MzI3MyAyIElOIElQNCAxMjcuMC4wLjFcclxucz0tXHJcbnQ9MCAwXHJcbmE9Z3JvdXA6QlVORExFIDBcclxuYT1leHRtYXAtYWxsb3ctbWl4ZWRcclxuYT1tc2lkLXNlbWFudGljOiBXTVMgeTVzRUZyMjlSQjYwN25yandjbTFTV2lzWjBqWkF5ZmF2T2tIXHJcbm09dmlkZW8gOSBVRFAvVExTL1JUUC9TQVZQRiA5NiA5NyA5OCA5OSAxMDAgMTAxIDEwMiAxMjEgMTI3IDEyMCAxMjUgMTA3IDEwOCAxMDkgMTI0IDExOSAxMjMgMTE4IDExNCAxMTUgMTE2XHJcbmM9SU4gSVA0IDAuMC4wLjBcclxuYT1ydGNwOjkgSU4gSVA0IDAuMC4wLjBcclxuYT1pY2UtdWZyYWc6NTZaclxyXG5hPWljZS1wd2Q6NmM5ZDdvODFHS2g4MHVOT2FLSXltZTNuXHJcbmE9aWNlLW9wdGlvbnM6dHJpY2tsZVxyXG5hPWZpbmdlcnByaW50OnNoYS0yNTYgODE6M0M6NkI6RTY6OUQ6Nzc6RkE6RDY6RTg6QTc6NkI6QTc6NTg6N0Q6REE6MDM6RUU6NjA6Q0M6ODc6MzE6MDQ6RDc6QjM6NzU6QTc6OTc6Qjc6Nzc6NkE6MzM6QTJcclxuYT1zZXR1cDphY3RwYXNzXHJcbmE9bWlkOjBcclxuYT1leHRtYXA6MSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDp0b2Zmc2V0XHJcbmE9ZXh0bWFwOjIgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvYWJzLXNlbmQtdGltZVxyXG5hPWV4dG1hcDozIHVybjozZ3BwOnZpZGVvLW9yaWVudGF0aW9uXHJcbmE9ZXh0bWFwOjQgaHR0cDovL3d3dy5pZXRmLm9yZy9pZC9kcmFmdC1ob2xtZXItcm1jYXQtdHJhbnNwb3J0LXdpZGUtY2MtZXh0ZW5zaW9ucy0wMVxyXG5hPWV4dG1hcDo1IGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L3BsYXlvdXQtZGVsYXlcclxuYT1leHRtYXA6NiBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC92aWRlby1jb250ZW50LXR5cGVcclxuYT1leHRtYXA6NyBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC92aWRlby10aW1pbmdcclxuYT1leHRtYXA6OCBodHRwOi8vd3d3LndlYnJ0Yy5vcmcvZXhwZXJpbWVudHMvcnRwLWhkcmV4dC9jb2xvci1zcGFjZVxyXG5hPWV4dG1hcDo5IHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6bWlkXHJcbmE9ZXh0bWFwOjEwIHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6cnRwLXN0cmVhbS1pZFxyXG5hPWV4dG1hcDoxMSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOnJlcGFpcmVkLXJ0cC1zdHJlYW0taWRcclxuYT1zZW5kcmVjdlxyXG5hPW1zaWQ6eTVzRUZyMjlSQjYwN25yandjbTFTV2lzWjBqWkF5ZmF2T2tIIDQyM2M3MTMzLWYwZjMtNDhkYS1hNTI4LTUzYWNlNzFhM2ZmMlxyXG5hPXJ0Y3AtbXV4XHJcbmE9cnRjcC1yc2l6ZVxyXG5hPXJ0cG1hcDo5NiBWUDgvOTAwMDBcclxuYT1ydGNwLWZiOjk2IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6OTYgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjo5NiBjY20gZmlyXHJcbmE9cnRjcC1mYjo5NiBuYWNrXHJcbmE9cnRjcC1mYjo5NiBuYWNrIHBsaVxyXG5hPXJ0cG1hcDo5NyBydHgvOTAwMDBcclxuYT1mbXRwOjk3IGFwdD05NlxyXG5hPXJ0cG1hcDo5OCBWUDkvOTAwMDBcclxuYT1ydGNwLWZiOjk4IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6OTggdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjo5OCBjY20gZmlyXHJcbmE9cnRjcC1mYjo5OCBuYWNrXHJcbmE9cnRjcC1mYjo5OCBuYWNrIHBsaVxyXG5hPWZtdHA6OTggcHJvZmlsZS1pZD0wXHJcbmE9cnRwbWFwOjk5IHJ0eC85MDAwMFxyXG5hPWZtdHA6OTkgYXB0PTk4XHJcbmE9cnRwbWFwOjEwMCBWUDkvOTAwMDBcclxuYT1ydGNwLWZiOjEwMCBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEwMCB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEwMCBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMDAgbmFja1xyXG5hPXJ0Y3AtZmI6MTAwIG5hY2sgcGxpXHJcbmE9Zm10cDoxMDAgcHJvZmlsZS1pZD0yXHJcbmE9cnRwbWFwOjEwMSBydHgvOTAwMDBcclxuYT1mbXRwOjEwMSBhcHQ9MTAwXHJcbmE9cnRwbWFwOjEwMiBIMjY0LzkwMDAwXHJcbmE9cnRjcC1mYjoxMDIgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMDIgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMDIgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTAyIG5hY2tcclxuYT1ydGNwLWZiOjEwMiBuYWNrIHBsaVxyXG5hPWZtdHA6MTAyIGxldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTE7cGFja2V0aXphdGlvbi1tb2RlPTE7cHJvZmlsZS1sZXZlbC1pZD00MjAwMWZcclxuYT1ydHBtYXA6MTIxIHJ0eC85MDAwMFxyXG5hPWZtdHA6MTIxIGFwdD0xMDJcclxuYT1ydHBtYXA6MTI3IEgyNjQvOTAwMDBcclxuYT1ydGNwLWZiOjEyNyBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEyNyB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEyNyBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMjcgbmFja1xyXG5hPXJ0Y3AtZmI6MTI3IG5hY2sgcGxpXHJcbmE9Zm10cDoxMjcgbGV2ZWwtYXN5bW1ldHJ5LWFsbG93ZWQ9MTtwYWNrZXRpemF0aW9uLW1vZGU9MDtwcm9maWxlLWxldmVsLWlkPTQyMDAxZlxyXG5hPXJ0cG1hcDoxMjAgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMjAgYXB0PTEyN1xyXG5hPXJ0cG1hcDoxMjUgSDI2NC85MDAwMFxyXG5hPXJ0Y3AtZmI6MTI1IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTI1IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTI1IGNjbSBmaXJcclxuYT1ydGNwLWZiOjEyNSBuYWNrXHJcbmE9cnRjcC1mYjoxMjUgbmFjayBwbGlcclxuYT1mbXRwOjEyNSBsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0xO3Byb2ZpbGUtbGV2ZWwtaWQ9NDJlMDFmXHJcbmE9cnRwbWFwOjEwNyBydHgvOTAwMDBcclxuYT1mbXRwOjEwNyBhcHQ9MTI1XHJcbmE9cnRwbWFwOjEwOCBIMjY0LzkwMDAwXHJcbmE9cnRjcC1mYjoxMDggZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMDggdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMDggY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTA4IG5hY2tcclxuYT1ydGNwLWZiOjEwOCBuYWNrIHBsaVxyXG5hPWZtdHA6MTA4IGxldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTE7cGFja2V0aXphdGlvbi1tb2RlPTA7cHJvZmlsZS1sZXZlbC1pZD00MmUwMWZcclxuYT1ydHBtYXA6MTA5IHJ0eC85MDAwMFxyXG5hPWZtdHA6MTA5IGFwdD0xMDhcclxuYT1ydHBtYXA6MTI0IEgyNjQvOTAwMDBcclxuYT1ydGNwLWZiOjEyNCBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEyNCB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEyNCBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMjQgbmFja1xyXG5hPXJ0Y3AtZmI6MTI0IG5hY2sgcGxpXHJcbmE9Zm10cDoxMjQgbGV2ZWwtYXN5bW1ldHJ5LWFsbG93ZWQ9MTtwYWNrZXRpemF0aW9uLW1vZGU9MTtwcm9maWxlLWxldmVsLWlkPTRkMDAzMlxyXG5hPXJ0cG1hcDoxMTkgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMTkgYXB0PTEyNFxyXG5hPXJ0cG1hcDoxMjMgSDI2NC85MDAwMFxyXG5hPXJ0Y3AtZmI6MTIzIGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTIzIHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTIzIGNjbSBmaXJcclxuYT1ydGNwLWZiOjEyMyBuYWNrXHJcbmE9cnRjcC1mYjoxMjMgbmFjayBwbGlcclxuYT1mbXRwOjEyMyBsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0xO3Byb2ZpbGUtbGV2ZWwtaWQ9NjQwMDMyXHJcbmE9cnRwbWFwOjExOCBydHgvOTAwMDBcclxuYT1mbXRwOjExOCBhcHQ9MTIzXHJcbmE9cnRwbWFwOjExNCByZWQvOTAwMDBcclxuYT1ydHBtYXA6MTE1IHJ0eC85MDAwMFxyXG5hPWZtdHA6MTE1IGFwdD0xMTRcclxuYT1ydHBtYXA6MTE2IHVscGZlYy85MDAwMFxyXG5hPXNzcmMtZ3JvdXA6RklEIDMyMjI1NjI4ODAgMjc4MDM0NjQ4MFxyXG5hPXNzcmM6MzIyMjU2Mjg4MCBjbmFtZTpmcnJvNSthZ3JwQmJOK0I0XHJcbmE9c3NyYzozMjIyNTYyODgwIG1zaWQ6eTVzRUZyMjlSQjYwN25yandjbTFTV2lzWjBqWkF5ZmF2T2tIIDQyM2M3MTMzLWYwZjMtNDhkYS1hNTI4LTUzYWNlNzFhM2ZmMlxyXG5hPXNzcmM6MzIyMjU2Mjg4MCBtc2xhYmVsOnk1c0VGcjI5UkI2MDducmp3Y20xU1dpc1owalpBeWZhdk9rSFxyXG5hPXNzcmM6MzIyMjU2Mjg4MCBsYWJlbDo0MjNjNzEzMy1mMGYzLTQ4ZGEtYTUyOC01M2FjZTcxYTNmZjJcclxuYT1zc3JjOjI3ODAzNDY0ODAgY25hbWU6ZnJybzUrYWdycEJiTitCNFxyXG5hPXNzcmM6Mjc4MDM0NjQ4MCBtc2lkOnk1c0VGcjI5UkI2MDducmp3Y20xU1dpc1owalpBeWZhdk9rSCA0MjNjNzEzMy1mMGYzLTQ4ZGEtYTUyOC01M2FjZTcxYTNmZjJcclxuYT1zc3JjOjI3ODAzNDY0ODAgbXNsYWJlbDp5NXNFRnIyOVJCNjA3bnJqd2NtMVNXaXNaMGpaQXlmYXZPa0hcclxuYT1zc3JjOjI3ODAzNDY0ODAgbGFiZWw6NDIzYzcxMzMtZjBmMy00OGRhLWE1MjgtNTNhY2U3MWEzZmYyXHJcbiJ9", &offer)

	log.Print("got me offer")
	err = pc.SetRemoteDescription(offer)
	if err != nil {
		log.Panic(err)
	}

	pc.OnTrack(func(tr *webrtc.TrackRemote, r *webrtc.RTPReceiver) {

		go func() {
			ticker := time.NewTicker(time.Second * 3)
			for range ticker.C {
				errSend := pc.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{
						MediaSSRC: uint32(tr.SSRC()),
					},
				})
				if errSend != nil {
					log.Println(errSend)
				}
			}
		}()

		log.Printf("Track has started, of type %d: %s \n", tr.PayloadType(), tr.Codec().MimeType)
		for {
			rtp, _, readErr := tr.ReadRTP()
			if readErr != nil {
				panic(readErr)
			}

			if writeErr := outputTrack.WriteRTP(rtp); writeErr != nil {
				panic(writeErr)
			}
		}
	})

	pc.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		log.Printf("Connection State has changed %s \n", is.String())
	})

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(pc)
	if err = pc.SetLocalDescription(answer); err != nil {
		log.Panic(err)
	}

	<-gatherComplete

	log.Println(signal.Encode(*pc.LocalDescription()))

	select {}
}
