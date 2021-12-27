package command

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/victorjacobs/csv2ynab/config"
)

func WatchDirectories(config config.Config) {
	done := make(chan bool)

	for _, w := range config.WatchDirectories {
		watchDir := w

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

					if strings.Contains(fileName, ".processed") {
						continue
					}

					log.Printf("Processing new file %q", fileName)

					err = ProcessFile(watchDir.Merge(config.Ynab), fileName, "")

					if err != nil {
						log.Printf("Processing failed: %v", err)
						continue
					}

					err = os.Rename(fileName, fmt.Sprintf("%v.processed", fileName))
					if err != nil {
						log.Printf("Moving processed file failed: %v", err)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					log.Printf("Watching failed: %v", err)
				}
			}
		}()

		log.Printf("Watching for new files in %q", watchDir.Path)
		err = watcher.Add(watchDir.Path)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
}
