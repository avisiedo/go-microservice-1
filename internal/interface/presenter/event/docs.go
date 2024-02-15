// Package event holds the presenter component for the asynchronous kafka events
//
// An initial granularity could be one file per topic, but if the
// events are only a few would make sense to use only one file.
package event
