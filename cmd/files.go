package cmd

import (
	"github.com/Masterminds/cookoo"
	fsnotify "gopkg.in/fsnotify.v1"
	"path/filepath"
	"path"
	"time"
	"fmt"
	"os"
)

// This file contains commands related to the file system.

// FindCodl finds all CODL files (*.codl) in a given directory.
func FindCodl(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", ".", p)

	where := filepath.Join(dir, "*.codl")
	files, err := filepath.Glob(where)

	return files, err
}

// Watch uses fsnotify to watch for changes to .codl files.
func Watch(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", ".", p)
	route := cookoo.GetString("update", "@update", p)

	r, ok := c.Has("router")
	if !ok {
		return time.Now(), fmt.Errorf("Could not find 'router' in context.")
	}

	router := r.(*cookoo.Router)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()
	watcher.Add(dir)

	fmt.Printf("[INFO] Watching %s for changes to .codl files.\n", dir)

	// Watch for updates to files.
	for {
		select {
		case good := <-watcher.Events:

			// Look for create, write, and rename events.
			switch good.Op {
			//case fsnotify.Create, fsnotify.Write, fsnotify.Rename:
			case fsnotify.Write, fsnotify.Create:
				if path.Ext(good.Name) != ".codl" {
					continue
				}
				fmt.Printf("[INFO] %s has changed. Updating. (%s)\n", good.Name, good.String())
				c.Put("files", []string{good.Name})
				err := router.HandleRequest(route, c, false)
				if err != nil {
					return time.Now(), err
				}
				c.Put("lastUpdated", time.Now())

			// Log but otherwise ignore Remove.
			case fsnotify.Remove:
				fmt.Printf("[INFO] %s has been removed.\n", good.Name)
			}
		case bad := <-watcher.Errors:
			c.Logf("warn", "Error watching: %s", bad.Error())
		}
	}
}


// FilterUnchanged takes a list of files and a timestamp and returns only those changed since the time.
func FilterUnchanged(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	files := p.Get("files", []string{}).([]string)
	since := p.Get("since", time.Now().Add(time.Second * -10)).(time.Time)

	modified := []string{}
	for _, file := range files {
		stat, err := os.Stat(file)
		if err == nil  && stat.ModTime().After(since) {
			modified = append(modified, file)
		}
	}

	return modified, nil
}

// Repeats a route at intervals of 'period'.
//
// Params
// 	- route: The route to repeat
// 	- period: The time.Duration to wait between executions.
//
// Note: This expects `contex.Get("router"...)` to return a *cookoo.Router.
func Repeat(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	route := p.Get("route", "default").(string)
	period := p.Get("period", time.Second).(time.Duration)

	r, ok := c.Has("router")
	if !ok {
		return time.Now(), fmt.Errorf("Could not find 'router' in context.")
	}

	router := r.(*cookoo.Router)

	for {
		c.Logf("info", "Repeat\n")
		err := router.HandleRequest(route, c, false)
		if err != nil {
			return time.Now(), err
		}
		c.Put("lastUpdated", time.Now())
		time.Sleep(period)
	}

	return time.Now(), nil
}
