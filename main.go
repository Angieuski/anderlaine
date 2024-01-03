package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Tarefa struct { // Estrutura das tarefas
	Id        int       `json:"id"`
	Titulo    string    `json:"titulo"`
	Descricao string    `json:"descricao"`
	Data      time.Time `json:"data"`
	Status    string    `json:"status"`
}

var Tarefas []Tarefa = []Tarefa{
	Tarefa{
		Id:        1,
		Titulo:    "Estudar",
		Descricao: "Estudar por no mínimo 2 horas nas férias",
		Data:      time.Now(),
		Status:    "Concluído",
	},
	Tarefa{
		Id:        2,
		Titulo:    "Academia",
		Descricao: "Ir para academia",
		Data:      time.Now(),
		Status:    "Concluído",
	},
	Tarefa{
		Id:        3,
		Titulo:    "Dormir",
		Descricao: "Dormir às 23:00h",
		Data:      time.Now(),
		Status:    "Inconcluído",
	},
}

func rotaPrincipal(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Bem vindo!")
}

func listarTarefas(w http.ResponseWriter, r *http.Request) { // apresentar tarefas armazenadas
	json.NewEncoder(w).Encode(Tarefas)
}

func criarTarefa(w http.ResponseWriter, r *http.Request) { // criar nova tarefa
	w.WriteHeader(http.StatusCreated)

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		log.Fatal(erro)
		return
	}

	var novaTarefa Tarefa
	json.Unmarshal(body, &novaTarefa)
	novaTarefa.Id = len(Tarefas) + 1
	Tarefas = append(Tarefas, novaTarefa)

	encoder := json.NewEncoder(w)
	encoder.Encode(novaTarefa)
}

func modificarTarefa(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		fmt.Println("Erro ao converter a string para inteiro:", err)
		return
	}

	body, erro := io.ReadAll(r.Body)
	if erro != nil {
		fmt.Println("Erro ao ler o corpo da requisição:", erro)
		return
	}

	var tarefaModificada Tarefa // modificação da tarefa
	json.Unmarshal(body, &tarefaModificada)

	idxTarefa := -1
	for indice, Tarefa := range Tarefas {
		if Tarefa.Id == id {
			idxTarefa = indice

			break
		}
	}
	if idxTarefa < 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// ... código existente ...
	Tarefas[idxTarefa] = tarefaModificada
	json.NewEncoder(w).Encode(tarefaModificada)
	fmt.Printf("ID da Tarefa: %d\n", id)
	fmt.Printf("Corpo da Requisição: %s\n", body)
	// ... código existente ...
}

func acaoTarefas(w http.ResponseWriter, r *http.Request) { // selecionar ação de CRUD de tarefas
	// /tarefas
	w.Header().Set("Content-Type", "application/json") // seta a resposta do header pra Json

	partes := strings.Split(r.URL.Path, "/")
	if len(partes) == 2 || len(partes) == 3 && partes[2] == "" {
		if r.Method == "GET" {
			listarTarefas(w, r)
		} else if r.Method == "POST" {
			criarTarefa(w, r)
		}
	} else if len(partes) == 3 || len(partes) == 4 && partes[3] == "" {
		if r.Method == "GET" {
			buscarTarefa(w, r)
		} else if r.Method == "DELETE" {
			excluirTarefa(w, r)
		} else if r.Method == "PUT" {
			modificarTarefa(w, r)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func buscarTarefa(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	// Quebrar a URL em partes nas barras,  [...] / tarefas / XXXXX
	partes := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(partes[2])
	if err != nil {
		fmt.Println("Erro ao converter a string para inteiro:", err)
		return
	}

	for _, Tarefa := range Tarefas {
		if Tarefa.Id == id {
			json.NewEncoder(w).Encode(Tarefa)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func excluirTarefa(w http.ResponseWriter, r *http.Request) {
	// DELETE /tarefas/XXXX
	partes := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		fmt.Println("Erro ao converter a string para inteiro:", err)
		return
	}

	for _, Tarefa := range Tarefas {
		if Tarefa.Id == id {
			esquerdaArray := Tarefas[0:id]
			direitaArray := Tarefas[id+1 : len(Tarefas)]
			Tarefas = append(esquerdaArray, direitaArray...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)

}

func svConfig() {
	roteador := mux.NewRouter()
	roteador.HandleFunc("/", rotaPrincipal)
	roteador.HandleFunc("/tarefas", acaoTarefas)
	roteador.HandleFunc("/tarefas/", acaoTarefas)

	log.Fatal(http.ListenAndServe(":1337", roteador)) // com o  Nil pra mux é utilizado o ServerMux padrão com servidor rodando na porta 1337
}

func main() {
	svConfig()
}
