package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	myCommands "ReadyGoBot/commands"

	"github.com/bwmarrin/discordgo"
	dotenv "github.com/joho/godotenv"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() {
	dotenv.Load()
	flag.Parse()
}

func init() {
	var err error
	token := *BotToken
	if token == "" {
		token = os.Getenv("SECRET_KEY")
	}
	authString := fmt.Sprintf("Bot %s", token)
	s, err = discordgo.New(authString)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	// integerOptionMinValue          = 1.0
	// dmPermission                   = false
	// defaultMemberPermissions int64 = discordgo.PermissionManageServer

	/** Prefix every command with a string */
	commandPrefix = "rg"
	/** Add every command in this array */
	commands = []*myCommands.RGCommand{
		&myCommands.AddTrackCommand,
		&myCommands.ListTracksCommand,
		&myCommands.AddHotlapCommand,
		&myCommands.GetLeaderboardCommand,
	}
	commandHandlers   = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), len(commands))
	componentHandlers = make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate), len(commands))
)

func bindCommand(cmd *myCommands.RGCommand) {
	cmd.Command.Name = commandPrefix + cmd.Command.Name

	commandHandlers[cmd.Command.Name] = cmd.Handler
	for componentId, componentHandler := range cmd.ComponentHandlers {
		componentHandlers[componentId] = componentHandler
	}
}

func init() {
	for _, v := range commands {
		bindCommand(v)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, &v.Command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Command.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
