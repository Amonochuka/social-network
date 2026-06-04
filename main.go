package main

import (
	"fmt"
	"net/http"
	"log"
	"socialnetwork/handlers" // ← this must match your module name in go.mod
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "socialnetwork server running on port 8080\n")
	})
		///Groups endpoints  handlers 

	http.HandleFunc("/groups", handlers.Groups)
	http.HandleFunc("/groupdetails", handlers.GroupDetails)
	http.HandleFunc("/creategroup", handlers.CreateGroup)
	http.HandleFunc("/inviteusers", handlers.InviteUsers)
	http.HandleFunc("/acceptinvitation", handlers.AcceptInvitation)
	http.HandleFunc("/rejectinvitation", handlers.RejectInvitation)
	http.HandleFunc("/invitationpermission", handlers.InvitationPermission)


	//Group Content  endpoint  handlers
	http.HandleFunc("/posts", handlers.Posts)
	http.HandleFunc("/comments", handlers.Comments)
	http.HandleFunc("/membership", handlers.Membership)

	///Events  endpoints handlers 

		http.HandleFunc("/createevents", handlers.CreateEvents)
	http.HandleFunc("/eventslistening", handlers.EventsListening)
	http.HandleFunc("/goingresponce", handlers.GoingResponce)
	http.HandleFunc("/notgoingresponce", handlers.NotGoingResponce)


	/*//WEBSOCKET HANDLER 
	//TESTING WEBSOCKETS 

	http.HandleFunc("/ws", handlers.HandleWebSocket)
    */

	//WEBSOCKET HUB 

   	broadcastHub := handlers.NewHub()
	go broadcastHub.Run()

	// Private chat hub (new)
	privateHub := handlers.NewPrivateHub()
	go privateHub.Run()

	roomManager := handlers.NewRoomManager()

	// Broadcast endpoint
	http.HandleFunc("/broadcastchat", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServeWs(broadcastHub, w, r)
	})

	// Private chat endpoint
	http.HandleFunc("/privatechat", func(w http.ResponseWriter, r *http.Request) {
		handlers.ServePrivateWs(privateHub, w, r)
	})

///wscat -c "ws://localhost:8080/groupchat/general?username=alice" TESTING  GROUPCHAT
		http.HandleFunc("/groupchat/", func(w http.ResponseWriter, r *http.Request) {
		roomName := handlers.ExtractRoomName(r.URL.Path)
		if roomName == "" || roomName == "groupchat" {
			http.Error(w, "Room name required. Use /groupchat/roomname", http.StatusBadRequest)
			return
		}
		room := roomManager.GetOrCreateRoom(roomName)
		handlers.ServeGroupWs(room, w, r)
	})

	//////////////////////////////
	log.Println("Broadcast: /broadcastchat")
	log.Println("Private: /privatechat")
	log.Println("Group: /groupchat/roomname")
	fmt.Println("socialnetwork server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
