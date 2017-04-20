package framework

import "github.com/Adirelle/blocksgo/lib/ipc"

// ClickableBlock handles click events
type ClickableBlock interface {
	AcceptClick(ipc.ClickEvent) bool
	Clicked(ipc.ClickEvent)
}
