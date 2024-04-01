package main

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	// Pasta 
	pasta := "C:\\Users\\tiesc\\PROJETOS\\folder_observer_go\\caminho_teste"

	// Arquivo log
	arquivoLog := "log.txt"

	// Inicializa o observer
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Adiciona a pasta ao observer
	err = watcher.Add(pasta)
	if err != nil {
		log.Fatal(err)
	}

	// Criaa ou abre o arquivo de logs
	logFile, err := os.OpenFile(arquivoLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Cria um logger para o log
	logger := log.New(logFile, "", log.LstdFlags)

	// observa mudanças

	// função alterada para registrar mudanças de exclusão e alterações de nome do arquivo
	registrarMudanca := func(event fsnotify.Event) {
		// Verifica se o arquivo ainda existe
		_, err := os.Stat(event.Name)
		existe := err == nil

		// operação realizada
		operacao := ""
		switch {
		case event.Op&fsnotify.Create == fsnotify.Create:
			operacao = "Criado"
		case event.Op&fsnotify.Remove == fsnotify.Remove:
			if !existe {
				operacao = "Removido"
			}
		case event.Op&fsnotify.Rename == fsnotify.Rename:
			if !existe {
				operacao = "Renomeado"
			}
		case event.Op&fsnotify.Write == fsnotify.Write:
			operacao = "Modificado"
		}

		// llog
		if existe {
			info, err := os.Stat(event.Name)
			if err != nil {
				logger.Printf("Erro ao obter informações sobre o arquivo %s: %s", event.Name, err)
				return
			}
			logger.Printf("%s - Arquivo: %s - Data/Hora: %s - Operação: %s", event.Name, info.ModTime().Format(time.RFC3339), time.Now().Format(time.RFC3339), operacao)
		} else {
			logger.Printf("%s - Data/Hora: %s - Operação: %s", event.Name, time.Now().Format(time.RFC3339), operacao)
		}
	}


	// Loop observer
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			registrarMudanca(event)
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Println("Erro:", err)
		}
	}
}
