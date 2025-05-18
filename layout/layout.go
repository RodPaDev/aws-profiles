package layout

type LayoutSizes struct {
	ScreenWidth  int
	ScreenHeight int

	LeftPaneWidth  int
	RightPaneWidth int

	MainHeight int

	FooterHeight  int
	FooterVisible int

	HeightAvailable int
}

func ComputeLayout(screenWidth, screenHeight int) LayoutSizes {
	const border = 2
	const footer = 4
	const left = 32
	const gap = 6 // space between panes (borders, margins, etc)

	heightAvailable := screenHeight - border
	footerVisible := footer - 1
	mainHeight := heightAvailable - footerVisible - 2 // for spacing/padding

	return LayoutSizes{
		ScreenWidth:     screenWidth,
		ScreenHeight:    screenHeight,
		LeftPaneWidth:   left,
		RightPaneWidth:  screenWidth - left - gap,
		MainHeight:      mainHeight,
		FooterHeight:    footer,
		FooterVisible:   footerVisible,
		HeightAvailable: heightAvailable,
	}
}
