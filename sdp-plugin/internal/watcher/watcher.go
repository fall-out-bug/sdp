package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatcherConfig holds configuration for the file watcher
type WatcherConfig struct {
	// IncludePatterns specifies glob patterns for files to include
	IncludePatterns []string

	// ExcludePatterns specifies glob patterns for files to exclude
	ExcludePatterns []string

	// DebounceInterval specifies how long to wait before processing changes
	DebounceInterval time.Duration

	// OnChange is called when a matching file changes
	OnChange func(path string)

	// OnError is called when an error occurs
	OnError func(error)
}

// Watcher monitors files for changes
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	config    *WatcherConfig
	watchPath string
	done      chan struct{}
	stopped   chan struct{}
	timer     *time.Timer
	lastFile  string
	mu        sync.Mutex // Protects timer and lastFile
}

// NewWatcher creates a new file watcher
func NewWatcher(watchPath string, config *WatcherConfig) (*Watcher, error) {
	if config == nil {
		config = &WatcherConfig{}
	}

	// Set default debounce interval
	if config.DebounceInterval == 0 {
		config.DebounceInterval = 100 * time.Millisecond
	}

	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	w := &Watcher{
		fsWatcher: fsWatcher,
		config:    config,
		watchPath: watchPath,
		done:      make(chan struct{}),
		stopped:   make(chan struct{}),
	}

	return w, nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	// Add watch path
	err := w.addWatch(w.watchPath)
	if err != nil {
		return fmt.Errorf("failed to add watch: %w", err)
	}

	// Process events
	for {
		select {
		case <-w.done:
			close(w.stopped)
			return nil

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return nil
			}
			w.handleEvent(event)

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return nil
			}
			if w.config.OnError != nil {
				w.config.OnError(err)
			}
		}
	}
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.done)
	w.fsWatcher.Close()
	<-w.stopped
}

// Close cleans up resources
func (w *Watcher) Close() {
	w.fsWatcher.Close()
}

func (w *Watcher) addWatch(path string) error {
	// Check if path exists
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Add directory to watcher
	if info.IsDir() {
		return w.fsWatcher.Add(path)
	}

	// If it's a file, watch its parent directory
	return w.fsWatcher.Add(filepath.Dir(path))
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only care about create and write events
	if event.Op&fsnotify.Create == 0 && event.Op&fsnotify.Write == 0 {
		return
	}

	// Skip directories
	info, err := os.Stat(event.Name)
	if err == nil && info.IsDir() {
		// Watch new directories
		if event.Op&fsnotify.Create != 0 {
			w.addWatch(event.Name)
		}
		return
	}

	// Check if file matches patterns
	if !w.matchesPatterns(event.Name) {
		return
	}

	// Debounce rapid changes
	w.debounce(event.Name)
}

func (w *Watcher) matchesPatterns(path string) bool {
	filename := filepath.Base(path)

	// Check exclude patterns
	for _, pattern := range w.config.ExcludePatterns {
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return false
		}
	}

	// If no include patterns, match all
	if len(w.config.IncludePatterns) == 0 {
		return true
	}

	// Check include patterns
	for _, pattern := range w.config.IncludePatterns {
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}

	return false
}

func (w *Watcher) debounce(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Reset timer if same file
	if w.timer != nil && w.lastFile == path {
		w.timer.Stop()
	}

	w.lastFile = path
	w.timer = time.AfterFunc(w.config.DebounceInterval, func() {
		if w.config.OnChange != nil {
			w.config.OnChange(path)
		}
	})
}
