package conf

const (
	ProjectName    = "forgeturl-server"
	ProjectVersion = "v1.0.0"
)

type PageType int

const (
	OwnerPage    PageType = 0 // start with O
	ReadOnlyPage PageType = 1 // start With R
	EditPage     PageType = 2 // start With E
	TempPage     PageType = 3 // start With T, for temporary pages

	OwnerPrefix    = uint8('O')
	ReadonlyPrefix = uint8('R')
	EditPrefix     = uint8('E')
	TempPrefix     = uint8('T') // start With T, for temporary pages
)
