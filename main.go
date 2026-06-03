package main

import (
	"fmt"
	"net/http"
	"socialnetwork/handlers" // ← this must match your module name in go.mod
	"github.com/gorilla/websocket"
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


	fmt.Println("socialnetwork server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
