package service

import (
	"errors"
	"log"
	"tic-tac-toe/pkg/common"
)

type GameService struct {
	queue *common.PlayerQueue
	games *common.GameList
}

func NewGameService() *GameService {
	return &GameService{
		queue: common.NewPlayerQueue(),
		games: common.NewGameList(),
	}
}

func (g *GameService) AddPlayerToQueue(player common.Player) error {
	playerInQueue := g.queue.InQueue(player)

	if playerInQueue {
		return errors.New("[gs] player already in queue")
	}

	q := g.queue

	if q.Size() == 0 {
		q.Enqueue(player)
		log.Printf("[gs] player added to queue %v\n", player)
	} else {
		log.Printf("[gs] player will be added to gameroom instead of queue %v", player)

		p, err := q.Dequeue()

		if err != nil {
			return errors.New("dequeue player error")
		}

		log.Printf("[gs] player dequeued %v", p)

		gr := common.NewGameRoom()
		gr.AddPlayer(p.User, p.Conn)
		gr.AddPlayer(player.User, player.Conn)
		g.games.Add(gr)
		gr.StartGame()

		log.Printf("[gs] new game room added with players:\n\tp1 %v\n\tp2 %v", p, player)
	}

	log.Println("[gs] queue size: ", q.Size())

	return nil
}

func (g *GameService) RemovePlayerFromQueue(player common.Player) (common.Player, error) {
	p, err := g.queue.Remove(player)

	if err != nil {
		return common.Player{}, err
	}

	log.Printf("[gs] player removed from queue %v", p)

	return p, nil
}

func (g *GameService) PlayerInQueue(player common.Player) bool {
	r := g.queue.InQueue(player)

	log.Printf("[gs] player %v in queue: %t", player, r)

	return r
}

func (g *GameService) CreateGameRoom() *common.GameRoom {
	gr := common.NewGameRoom()
	g.games.Add(gr)

	log.Print("[gs] new game room added")

	return gr
}

func (g *GameService) GameRoomList() *common.GameList {
	return g.games
}
