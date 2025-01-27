package importer

import (
	"log"
	"os"
	"path/filepath"

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
				baseFileName := filepath.Base(fileName)
				var matchedPattern *cfg.WatchPattern

				for _, p := range config.Watch.Patterns {
					match, err := filepath.Match(p.Pattern, baseFileName)
					if err != nil {
						log.Printf("Matching failed: %v", err)
						continue
					}

					if match {
						matchedPattern = &p
						break
					}
				}

				if matchedPattern == nil {
					continue
				}

				log.Printf("Processing new file %q", fileName)

				err = ProcessFile(matchedPattern.Merge(config.Ynab), fileName, "")

				if err != nil {
					log.Printf("Processing failed: %v", err)
					continue
				}

				err = os.Remove(fileName)
				if err != nil {
					log.Printf("Removing processed file failed: %v", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Printf("Watching failed: %v", err)
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
