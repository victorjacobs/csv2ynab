package importer

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fsnotify/fsnotify"
	cfg "github.com/victorjacobs/csv2ynab/config"
)

func Watch(config cfg.Config) {
	done := make(chan bool)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					continue
				}

				fileName := event.Name

				ynabConfig, err := config.YNABConfigForFile(fileName)
				if err != nil {
					log.Errorf("Failed to get YNAB config: %v", err)
				}

				log.Printf("Processing new file %v", fileName)

				if err := ProcessFile(ynabConfig, fileName, ""); err != nil {
					log.Errorf("Processing failed: %v", err)
					continue
				}

				if err := os.Remove(fileName); err != nil {
					log.Errorf("Removing processed file failed: %v", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Errorf("Watching failed: %v", err)
			}
		}
	}()

	log.Printf("Watching for new files in %q", config.Watch.Directory)
	err = watcher.Add(config.Watch.Directory)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
