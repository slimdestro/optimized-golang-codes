/**
	@ High performance WebRTC server in Golang
	@ slimdestro
*/ 

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/pion/webrtc/v2"
)

func main() {
	// Create a new WebRTC configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new WebRTC peer connection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new audio track
	audioTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeOpus, 12345, "audio", "pion")
	if err != nil {
		log.Fatal(err)
	}

	// Create a new video track
	videoTrack, err := peerConnection.NewTrack(webrtc.DefaultPayloadTypeVP8, 12346, "video", "pion")
	if err != nil {
		log.Fatal(err)
	}

	// Add the audio and video tracks to the peer connection
	if _, err = peerConnection.AddTrack(audioTrack); err != nil {
		log.Fatal(err)
	}
	if _, err = peerConnection.AddTrack(videoTrack); err != nil {
		log.Fatal(err)
	}

	// Create an offer to send to the remote peer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the local description
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		log.Fatal(err)
	}

	// Create an HTTP server to receive the remote peer's answer
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the remote description
		if err = peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  r.FormValue("sdp"),
		}); err != nil {
			log.Fatal(err)
		}
	})

	// Start the HTTP server
	go http.ListenAndServe("localhost:8080", nil)

	// Wait for the remote peer's answer
	select {
	case <-peerConnection.SignalingStateChange:
	case <-peerConnection.ICEConnectionStateChange:
	case <-peerConnection.ICEGatheringStateChange:
	case <-peerConnection.ICECandidateChange:
	case <-peerConnection.TrackChange:
	case <-context.Done():
	}

	// Close the peer connection
	peerConnection.Close()
}