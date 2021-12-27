package command

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/victorjacobs/csv2ynab/config"
)

func WatchDirectories(config config.Config) {
	log.Println("Watching for new files")

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

					if event.Op&fsnotify.Create != fsnotify.Create {
						continue
					}

					log.Printf("Processing new file %q", event.Name)

					err = ProcessFile(watchDir.Merge(config.Ynab), event.Name, "")

					if err != nil {
						log.Printf("Processing failed: %v", err)
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					log.Printf("Watching failed: %v", err)
				}
			}
		}()

		err = watcher.Add(watchDir.Path)
		if err != nil {
			log.Fatal(err)
		}
	}

	<-done
}
