// Package pluginmgr loads plugins at startup and dispatches events to them.
// A panicking plugin does not take down ixr; panics are caught and logged.
package pluginmgr
