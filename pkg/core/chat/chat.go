package chat

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

func (service *Service) Start() {
	//conn, err := service.pool.Acquire(context.Background())
	//defer conn.Release()
	_, err := service.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS messages (
		id 			 BIGSERIAL unique,
		sender_id    INTEGER NOT NULL,
		recipient_id INTEGER NOT NULL,
   		message		 TEXT NOT NULL,
   		time    	 TIMESTAMP DEFAULT now() NOT NULL
	); `)
	if err != nil {
		log.Print("Can't repo Start()", err)
	}
	log.Print("repo Start()")

	_, err = service.pool.Exec(context.Background(), `
	INSERT INTO messages(id, sender_id, recipient_id, message)
	VALUES (0, 0, 0, 'hello world');`)
	if err != nil {
		log.Println("can't exec, insert message")
	}

	log.Print("Has hello world")
}

// created timestamp
// modified timestamp

// CRUD
func (service *Service) All() (models []ModelOperationsLog, err error) {
	rows, err := service.pool.Query(context.Background(), `SELECT id, name, number,recipientSender,count, balanceold, balancenew, time, owner_id FROM historyoperationslog;`)
	if err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := ModelOperationsLog{}
		err = rows.Scan(
			&model.Id,
			&model.Name,
			&model.Number,
			&model.RecipientSender,
			&model.Count,
			&model.BalanceOld,
			&model.BalanceNew,
			&model.Time,
			&model.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't get chat from db: %w", err)
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	return models, nil
}

func (service *Service) UserShowTransferLogByIdCard(idCard int, idUser int) (model []ModelOperationsLog, err error) {
	modHistoryLog := ModelOperationsLog{}
	err = service.pool.QueryRow(context.Background(), `
SELECT id, name, number, recipientSender, count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog WHERE id=$1 and owner_id=$2`, idCard, idUser).Scan(
		&modHistoryLog.Id,
		&modHistoryLog.Name,
		&modHistoryLog.Number,
		&modHistoryLog.RecipientSender,
		&modHistoryLog.Count,
		&modHistoryLog.BalanceOld,
		&modHistoryLog.BalanceNew,
		&modHistoryLog.Time,
		&modHistoryLog.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	model = append(model, modHistoryLog)
	return model, nil
}

func (service *Service) AdminShowTransferLogByIdCadr(id int) (model []ModelOperationsLog, err error) {
	modHistoryLog := ModelOperationsLog{}
	err = service.pool.QueryRow(context.Background(), `
SELECT id, name, number, recipientSender, count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog WHERE id=$1`, id).Scan(
		&modHistoryLog.Id,
		&modHistoryLog.Name,
		&modHistoryLog.Number,
		&modHistoryLog.RecipientSender,
		&modHistoryLog.Count,
		&modHistoryLog.BalanceOld,
		&modHistoryLog.BalanceNew,
		&modHistoryLog.Time,
		&modHistoryLog.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	model = append(model, modHistoryLog)
	return model, nil
}

func (service *Service) ShowOperationsLogByOwnerId(id int) (models []ModelOperationsLog, err error) {
	rows, err := service.pool.Query(context.Background(), `
SELECT id, name, number,recipientSender,count, balanceold, balancenew, time, owner_id 
FROM historyoperationslog  
WHERE owner_id= $1`, id)
	if err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := ModelOperationsLog{}
		err = rows.Scan(
			&model.Id,
			&model.Name,
			&model.Number,
			&model.RecipientSender,
			&model.Count,
			&model.BalanceOld,
			&model.BalanceNew,
			&model.Time,
			&model.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("can't get chat from db: %w", err)
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("can't get chat from db: %w", err)
	}
	return models, nil
}

func (service *Service) AddNewHistory(model ModelOperationsLog) (err error) {
	log.Print("started add new chat")
	log.Print("add model to db")
	_, err = service.pool.Exec(context.Background(), `
	INSERT INTO historyoperationslog(name, number,recipientSender,count, balanceold, balancenew, time, owner_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		model.Name,
		model.Number,
		model.RecipientSender,
		model.Count,
		model.BalanceOld,
		model.BalanceNew,
		model.Time,
		model.OwnerID,
	)
	if err != nil {
		log.Printf("can't exec insert add chat card: %d", err)
		return err
	}
	log.Print("saved model to db")
	log.Print("finish add model to db")
	return nil
}

func (service *Service) AddMassage(model ModelMassage) (err error) {
	log.Print("started add new Massage")
	log.Print("add model to db")
	_, err = service.pool.Exec(context.Background(), `
	INSERT INTO messages(sender_id, recipient_id, message)
	VALUES ($1, $2, $3)`,
		model.SenderID,
		model.RecipientID,
		model.Message,
	)
	if err != nil {
		log.Printf("can't exec insert add message: %v", err)
		return err
	}
	log.Print("saved model to db")
	log.Print("finish add model to db")
	return nil
}

func (service *Service) GetMessageByRecipientID(senderID, recipientID int) (models []ModelMassage, err error) {
	log.Print("started get Massage")
	log.Print("get model to db")
	rows, err := service.pool.Query(context.Background(), `
	SELECT id, message, sender_id, recipient_id FROM messages WHERE sender_id = $1 AND recipient_id = $2; `,
		senderID,
		recipientID,
	)
	if err != nil {
		log.Printf("can't query select message: %v", err)
		return nil, err
	}

	for rows.Next() {
		model := ModelMassage{}
		err = rows.Scan(
			&model.ID,
			&model.SenderID,
			&model.RecipientID,
			&model.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("can't get message from db: %w", err)
		}
		models = append(models, model)
	}
	log.Print("get model to db")
	log.Print("finish get model to db")
	return nil, nil
}

func (service *Service) GetMessageAll(senderID int) (models []ModelMassage, err error) {
	log.Print("started get Massage")
	log.Print("get model to db")
	rows, err := service.pool.Query(context.Background(), `
	SELECT sender_id, recipient_id FROM messages WHERE sender_id = $1 OR recipient_id = $1 GROUP BY sender_id, recipient_id; `,
		senderID,
	)
	if err != nil {
		log.Printf("can't query select message: %v", err)
		return nil, err
	}

	for rows.Next() {
		model := ModelMassage{}
		err = rows.Scan(
			//&model.ID,
			&model.SenderID,
			&model.RecipientID,
			//&model.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("can't get message from db: %w", err)
		}
		models = append(models, model)
	}
	log.Print("get model to db")
	log.Print("finish get model to db")
	return models, nil
}

type ModelOperationsLog struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Number          string `json:"number"`
	RecipientSender string `json:"recipientsender"`
	Count           int64  `json:"count"`
	BalanceOld      int64  `json:"balanceold"`
	BalanceNew      int64  `json:"balancenew"`
	Time            int64  `json:"time"`
	OwnerID         int64  `json:"ownerid"`
}

type ModelMassage struct {
	ID            int       `json:"id"`
	SenderID      int       `json:"sender_id"`
	RecipientID   int       `json:"recipient_id"`
	RecipientName string    `json:"recipient_name"`
	Message       string    `json:"message"`
	Time          time.Time `json:"time"`
}

type ModelTransferMoneyCardToCard struct {
	IdCardSender        int    `json:"id_card_sender"`
	NumberCardRecipient string `json:"number_card_recipient"`
	Count               int64  `json:"count"`
}
