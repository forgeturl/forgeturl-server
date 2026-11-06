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
	InvalidPage  PageType = -1

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
	if pageId == "" {
		return InvalidPage
	}
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
	return InvalidPage
}

// IDs are a type prefix followed by the generator's Base58 token. Accept the
// full int64 token length so existing and future generated links remain valid.
func ValidPageID(id string) bool {
	if len(id) < 2 || len(id) > 13 || ParseIdType(id) == InvalidPage {
		return false
	}
	const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	for _, char := range id[1:] {
		found := false
		for _, allowed := range alphabet {
			if char == allowed {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
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
