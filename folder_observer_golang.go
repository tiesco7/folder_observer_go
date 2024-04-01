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

	// Arquivo de log
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

	// mudanças
	registrarMudanca := func(event fsnotify.Event) {
		info, err := os.Stat(event.Name)
		if err != nil {
			logger.Printf("Erro ao obter informações sobre o arquivo %s: %s", event.Name, err)
			return
		}

		// operação realizada
		operacao := ""
		if event.Op&fsnotify.Create == fsnotify.Create {
			operacao = "Criado"
		}
		if event.Op&fsnotify.Remove == fsnotify.Remove {
			operacao = "Removido"
		}
		if event.Op&fsnotify.Rename == fsnotify.Rename {
			operacao = "Renomeado"
		}
		if event.Op&fsnotify.Write == fsnotify.Write {
			operacao = "Modificado"
		}

		// log
		logger.Printf("%s - Arquivo: %s - Data/Hora: %s - Operação: %s", event.Name, info.ModTime().Format(time.RFC3339), time.Now().Format(time.RFC3339), operacao)
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
