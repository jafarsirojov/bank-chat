package app

import (
	"bank-chat/pkg/core/auth"
	"bank-chat/pkg/core/chat"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jafarsirojov/mux/pkg/mux"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
	"time"
)

type MainServer struct {
	exactMux *mux.ExactMux
	cardsSvc *chat.Service
}

func NewMainServer(exactMux *mux.ExactMux, cardsSvc *chat.Service) *MainServer {
	return &MainServer{exactMux: exactMux, cardsSvc: cardsSvc}
}

func (m *MainServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.exactMux.ServeHTTP(writer, request)
}

func (m *MainServer) HandleGetMessageAll(writer http.ResponseWriter, request *http.Request) {
	log.Print("start handler chat")
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	//if authentication.Id == 0 {
	//	log.Print("admin")
	//	log.Print("all chat")
	//	models, err := m.cardsSvc.All()
	//	if err != nil {
	//		log.Print("can't get all chat")
	//		writer.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	log.Print(models)
	//	err = rest.WriteJSONBody(writer, models)
	//	if err != nil {
	//		log.Print("can't write json get all chat")
	//		writer.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//}
	log.Print("all chat cards user")
	models, err := m.cardsSvc.GetMessageAll(authentication.Id)
	if err != nil {
		log.Print("can't get all chat")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print(models)
	token := request.Header.Get("Authorization")
	log.Printf("Bearer token %s", token)
	token = token[7:]
	log.Printf("token %s", token)

	for _, model := range models {
		us, err := getUserToSvcAuth(model.RecipientID, token)
		if err != nil {
			log.Println("err getUserToSvcAuth", err)
			continue
		}
		model.RecipientName = us.Name
	}

	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Print("can't write json get all chat")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("finish chat handler")
}

func (m *MainServer) HandleGetMessageByRecipientID(writer http.ResponseWriter, request *http.Request) {
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("transfer chat by id card")
	value, ok := mux.FromContext(request.Context(), "id")
	if !ok {
		log.Print("can't chat by id card")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	recipientID, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("can't strconv atoi: %d", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	//if authentication.Id == 0 {
	//	log.Print("admin")
	//	models, err := m.cardsSvc.GetMessageByRecipientID(id)
	//	if err != nil {
	//		log.Printf("can't chat %d", err)
	//		writer.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	log.Print(models)
	//	err = rest.WriteJSONBody(writer, models)
	//	if err != nil {
	//		log.Printf("can't write json get all chat %d", err)
	//		writer.WriteHeader(http.StatusInternalServerError)
	//		return
	//	}
	//	return
	//}

	models, err := m.cardsSvc.GetMessageByRecipientID(authentication.Id, recipientID)
	if err != nil {
		log.Printf("can't chat is not owner %d", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print(models)
	err = rest.WriteJSONBody(writer, models)
	if err != nil {
		log.Printf("can't write json get all chat %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (m *MainServer) HandlePostAddMassage(writer http.ResponseWriter, request *http.Request) {
	log.Print("starting save new chat")
	authentication, ok := jwt.FromContext(request.Context()).(*auth.Auth)
	if !ok {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is not ok")
		return
	}
	if authentication == nil {
		writer.WriteHeader(http.StatusBadRequest)
		log.Print("can't authentication is nil")
		return
	}
	log.Print(authentication)
	log.Print("post add chat transfer")
	model := chat.ModelMassage{}

	log.Print("start read json body is save new chat")
	err := rest.ReadJSONBody(request, &model)
	if err != nil {
		log.Printf("can't READ json POST model: %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("finish read json body is save new chat")
	log.Println(model)
	if model.ID != 0 {
		log.Printf("id card not 0!")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Print("start func add new chat is handler")
	err = m.cardsSvc.AddMassage(model)
	if err != nil {
		log.Printf("can't add (save) chat %d", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("finish add new chat")
}

func getUserToSvcAuth( /*data cards.ModelOperationsLog,*/ id int, token string) (us User, err error) {
	log.Print("starting sender request to history Svc")
	//requestBody, err := json.Marshal(data)
	//if err != nil {
	//	return fmt.Errorf("can't encode requestBody %v: %w", data, err)
	//}
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/users/%d", addrAuthSvc, id),
		bytes.NewBuffer(nil),
	)
	if err != nil {
		return User{}, fmt.Errorf("can't create request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	log.Print("started sender request to auth Svc")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return User{}, fmt.Errorf("can't send request: %w", err)
	}
	defer response.Body.Close()

	log.Print("finish sender request to history Svc")

	err = json.NewDecoder(response.Body).Decode(&us)
	if err != nil {
		log.Println("Can't getUserToSvcAuth, json.NewDecoder(response.Body).Decode(&us)", err)
		return User{}, err
	}

	switch response.StatusCode {
	case 200:
		log.Print("200 request to auth Svc")
		return us, nil
	case 400:
		log.Print("400 request to auth Svc")
		return User{}, fmt.Errorf("bad request is server: %s", addrAuthSvc)
	case 401:
		log.Print("401 unauthorized to auth Svc")
		return User{}, fmt.Errorf("unauthorized is server: %s", addrAuthSvc)
	case 500:
		log.Print("500 request to auth Svc")
		return User{}, fmt.Errorf("internel server error is server: %s", addrAuthSvc)
	default:
		log.Printf("response status code: %s", err)
		return User{}, fmt.Errorf("err: %s", addrAuthSvc)
	}
}

const addrAuthSvc = "http://localhost:9011"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}
