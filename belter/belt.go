package Belt

import (
	"encoding/json"
	fmt "fmt"
	"io/ioutil"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Version struct {
	Major               byte
	Minor               byte
	Build               byte
	Experimental        bool
	ExperimentalVersion byte
}

type BlueBot struct {
	dg      *discordgo.Session
	Debug   bool
	version Version
	OwnID   string
	OwnAV   string
	OwnName string
	Stop    bool
}

// Vars after this

var bbb *BlueBot
var BotFunc string

// Functions after this

func BBReady(s *discordgo.Session, r *discordgo.Ready) {
	bbb.OwnID = r.User.ID
	bbb.OwnAV = r.User.Avatar
	bbb.OwnName = r.User.Username
	fmt.Println("Discord: Ready message received\nBB: I am '" + bbb.OwnName + "'!\nBB: My User ID: " + bbb.OwnID)
	for _, G := range r.Guilds {
		AttachToGuild(G)
		fmt.Println("Server " + G.Name + "'s Owner is " + GetownerName(G))
	}
}

func GetownerName(g *discordgo.Guild) string {
	Owner := g.OwnerID
	for _, M := range g.Members {
		if M.User.ID == Owner {
			OwnerName := M.User.Username
			return OwnerName
		}
	}
	return "Error"
}

func AttachToGuild(g *discordgo.Guild) {
	GID, GIDerr := json.Marshal(g)
	if GIDerr == nil {
		ioutil.WriteFile(g.ID+".GLD", GID, 0777)
	} else {
		fmt.Println("GLD writing error: " + GIDerr.Error())
	}
}

func BBGuildCreate(s *discordgo.Session, m *discordgo.GuildCreate) {
	fmt.Println("I have joined a new guild!")
	ProcessGuildCreate(m.Guild)
}

func ProcessGuildCreate(g *discordgo.Guild) {
	AttachToGuild(g)
}

func Initialize(Token string) {
	isdebug, err := ioutil.ReadFile("debugtoggle")
	bbb = &BlueBot{
		version: Version{0, 0, 1, true, 1},
		Debug:   (err == nil && len(isdebug) > 0),
		Stop:    false,
	}
	bbb.dg, err = discordgo.New(Token)
	if err != nil {
		fmt.Println("Discord Session error, check token, error message: " + err.Error())
		return
	} else {
		fmt.Println("Discord session initiated, token: " + Token)
	}
	// handlers
	bbb.dg.AddHandler(BBReady)
	bbb.dg.AddHandler(BBGuildCreate)

	fmt.Println("BB: Handlers installed")

	err = bbb.dg.Open()
	if err == nil {
		fmt.Println("Discord: Connection established")
		for !bbb.Stop {
			time.Sleep(400 * time.Millisecond)
		}
	} else {
		fmt.Println("Error opening websocket connection: ", err.Error())
	}
	fmt.Println("BBB: Blue belt stopping...")
	bbb.dg.Close()
}
