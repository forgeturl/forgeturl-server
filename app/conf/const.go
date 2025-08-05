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
	AdminPage    PageType = 4 // start With A, for admin pages, not used yet

	OwnerPrefix    = uint8('O')
	ReadonlyPrefix = uint8('R')
	EditPrefix     = uint8('E')
	TempPrefix     = uint8('T') // start With T, for temporary pages
	AdminPrefix    = uint8('A')

	ReadOnlyStr = "readonly"
	EditStr     = "edit"
	AdminStr    = "admin"
)

func ParseIdType(pageId string) PageType {
	switch pageId[0] {
	case OwnerPrefix:
		return OwnerPage
	case ReadonlyPrefix:
		return ReadOnlyPage
	case EditPrefix:
		return EditPage
	case TempPrefix:
		return TempPage
	case AdminPrefix:
		return AdminPage
	}
	return OwnerPage
}

func (p PageType) IsOwnPage() bool {
	return p == OwnerPage
}

func (p PageType) IsReadOnlyPage() bool {
	return p == ReadOnlyPage
}

func (p PageType) IsEditPage() bool {
	return p == EditPage
}

func (p PageType) IsTempPage() bool {
	return p == TempPage
}

func (p PageType) IsAdminPage() bool {
	return p == AdminPage
}
