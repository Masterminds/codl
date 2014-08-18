package cmd

import (
	"github.com/Masterminds/cookoo"
	"path/filepath"
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
