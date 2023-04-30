package shell

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	clientErr "github.com/blokhinnv/gophkeeper/internal/client/errors"
	"github.com/blokhinnv/gophkeeper/internal/client/service"
	"github.com/blokhinnv/gophkeeper/internal/server/models"
)

// ShellController defines the interface for the shell controller.
type ShellController interface {
	// ShowMenu displays the available options based on user authentication.
	ShowMenu()
	// SetAuth sets the authentication status and token for the shell controller.
	SetAuth(token string)
}

// ShellOption is a type for available shell options
type ShellOption string

// Constants for available shell options
const (
	Login    ShellOption = "login"
	Register ShellOption = "register"
	Sync     ShellOption = "sync"
	Show     ShellOption = "show"
	Add      ShellOption = "add"
	Update   ShellOption = "update"
	Delete   ShellOption = "delete"
	Quit     ShellOption = "quit"
)

var (
	// unauthrorizedOpts is a slice of ShellOption for unauthorized users
	unauthrorizedOpts = []ShellOption{Login, Register, Quit}
	// authorizedOpts is a slice of ShellOption for authorized users
	authorizedOpts = []ShellOption{Sync, Show, Add, Update, Delete, Quit}
)

// shellController implements the ShellController interface.
type shellController struct {
	Auth       bool
	Options    []ShellOption
	Token      string
	DataString string
	Data       any

	authService    service.AuthService
	syncService    service.SyncService
	storageService service.StorageService

	listener net.Listener
}

// NewShellController creates a new shell controller with the specified base URL.
func NewShellController(serverBaseURL string) ShellController {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("Error while creating a listener: %v", err)
	}
	ctrl := &shellController{
		authService:    service.NewAuthService(serverBaseURL),
		syncService:    service.NewSyncService(serverBaseURL),
		storageService: service.NewStorageService(serverBaseURL),
		listener:       listener,
	}
	go ctrl.listenerLoop()
	return ctrl
}

// listenerLoop accepts connections from server and syncs the data for all the
// available clients of the same user.
func (s *shellController) listenerLoop() {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			log.Fatalln("Error accepting: ", err.Error())
		}
		s.sync()
		connection.Close()
	}
}

func (s *shellController) registerClient() {
	r, err := s.syncService.Register(s.Token, s.listener.Addr().String())
	if err != nil {
		log.Fatalf("unable to register client")
		return
	}
	fmt.Println(r)
}

func (s *shellController) unregisterClient() {
	r, err := s.syncService.Unregister(s.Token, s.listener.Addr().String())
	if err != nil {
		log.Fatalf("unable to unregister client")
		return
	}
	fmt.Println(r)
}

// SetAuth sets the authentication status and token for the shell controller.
func (s *shellController) SetAuth(token string) {
	s.Token = token
	s.Auth = true
	s.registerClient()
}

// ShowMenu displays the available options based on user authentication.
func (s *shellController) ShowMenu() {
	opts := unauthrorizedOpts
	if s.Auth {
		opts = authorizedOpts
	}
	selectedOpt := ShellOption(selectItem("Select action: ", opts))
	switch selectedOpt {
	case Quit:
		s.quit()
	case Login:
		s.login()
	case Register:
		s.register()
	case Sync:
		s.sync()
	case Show:
		s.show()
	case Add:
		s.add()
	case Update:
		s.update()
	case Delete:
		s.delete()
	}
}

// quit quits the shell session.
func (s *shellController) quit() {
	fmt.Println("Bye!")
	s.listener.Close()
	s.unregisterClient()
	os.Exit(0)
}

// login prompts the user for login credentials and authenticates the user.
func (s *shellController) login() {
	username := promptText("Username: ")
	password := promptText("Password: ")
	tok, err := s.authService.Auth(username, password)
	if err != nil {
		if errors.Is(err, clientErr.ErrServerUnavailable) {
			log.Fatalln(err)
		}
		fmt.Println(err)
		return
	}
	fmt.Println("Token: ", tok)
	s.Token = tok
	s.Auth = true
	s.registerClient()
}

// register prompts the user for registration credentials and registers the user.
func (s *shellController) register() {
	username := promptText("Username: ")
	password := promptText("Password: ")
	err := s.authService.Register(username, password)
	if err != nil {
		if errors.Is(err, clientErr.ErrServerUnavailable) {
			log.Fatalln(err)
		}
		fmt.Println(err)
		return
	}
}

// sync retrieves the data from the server and stores it in the shell controller.
func (s *shellController) sync() {
	fmt.Println("sync....")
	syncResp, err := s.syncService.Sync(s.Token, models.AllowedCollection)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.Data = syncResp
}

// show displays the data received from the last sync operation
// in a human-readable format, with indentation and formatting applied.
func (s *shellController) show() {
	resJSON, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(resJSON))
}

// add prompts the user for the required fields to add a new record
// to the selected collection
func (s *shellController) add() {
	selectedCollection := models.Collection(
		selectItem("Select collection: ", models.AllowedCollection),
	)
	body, err := getBody(selectedCollection, false)
	if err != nil {
		fmt.Println("unable to create request body: ", err)
		return
	}
	msg, err := s.storageService.Add(body, selectedCollection, s.Token)
	if err != nil {
		if errors.Is(err, clientErr.ErrServerUnavailable) {
			log.Fatalln(err)
		}
		fmt.Println("unable to add the record: ", err)
		return
	}
	fmt.Println(msg)
}

// update prompts the user for the required fields to update an existing
// record in the selected collection.
func (s *shellController) update() {
	selectedCollection := models.Collection(
		selectItem("Select collection: ", models.AllowedCollection),
	)
	body, err := getBody(selectedCollection, true)
	if err != nil {
		fmt.Println("unable to create request body: ", err)
		return
	}
	msg, err := s.storageService.Update(body, selectedCollection, s.Token)
	if err != nil {
		if errors.Is(err, clientErr.ErrServerUnavailable) {
			log.Fatalln(err)
		}
		fmt.Println("unable to update the record: ", err)
		return
	}
	fmt.Println(msg)
}

// delete prompts the user for the record ID of the record to delete
// from the selected collection.
func (s *shellController) delete() {
	selectedCollection := models.Collection(
		selectItem("Select collection: ", models.AllowedCollection),
	)
	recordID := promptText("Record ID: ")
	body := fmt.Sprintf(`{"record_id": "%v"}`, recordID)
	msg, err := s.storageService.Delete(body, selectedCollection, s.Token)
	if err != nil {
		if errors.Is(err, clientErr.ErrServerUnavailable) {
			log.Fatalln(err)
		}
		fmt.Println("unable to delete the record: ", err)
		return
	}
	fmt.Println(msg)
}
