package tui

const (
	layoutVerticalChrome      = 6
	layoutMinBodyHeight       = 1
	layoutNarrowThreshold     = 100
	layoutNarrowHorizontalPad = 2
	layoutNavRatioDivisor     = 4
	layoutDetailRatioDivisor  = 3
	layoutInnerGap            = 2
	layoutMinListWidth        = 30
	layoutMinNarrowListWidth  = 1
)

type layoutSpec struct {
	narrow      bool
	navWidth    int
	listWidth   int
	detailWidth int
	bodyHeight  int
}

func calculateLayout(width int, height int) layoutSpec {
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	bodyHeight := height - layoutVerticalChrome
	bodyHeight = max(bodyHeight, layoutMinBodyHeight)

	if width < layoutNarrowThreshold {
		listWidth := max(layoutMinNarrowListWidth, width-layoutNarrowHorizontalPad)
		return layoutSpec{narrow: true, listWidth: listWidth, bodyHeight: bodyHeight}
	}

	navWidth := width / layoutNavRatioDivisor
	detailWidth := width / layoutDetailRatioDivisor
	listWidth := width - navWidth - detailWidth - layoutInnerGap
	listWidth = max(listWidth, layoutMinListWidth)

	return layoutSpec{
		narrow:      false,
		navWidth:    navWidth,
		listWidth:   listWidth,
		detailWidth: detailWidth,
		bodyHeight:  bodyHeight,
	}
}
