var all_users = make(map[int]*User)

type User struct {
	Nick   string
	ircObj map[string]*irc.Connection
	ws     *websocket.Conn
}