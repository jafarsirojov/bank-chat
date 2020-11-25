package app

import (
	"bank-chat/pkg/core/auth"
	"bank-chat/pkg/core/chat"
	"github.com/jafarsirojov/mux/pkg/mux"
	"github.com/jafarsirojov/mux/pkg/mux/middleware/jwt"
	"github.com/jafarsirojov/rest/pkg/rest"
	"log"
	"net/http"
	"strconv"
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
