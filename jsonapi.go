package teamtailorgo

type MarshalIdentifier interface {
	GetID() string
}

type UnmarshalIdentifier interface {
	SetID(string) error
}

type UnmarshalToOneRelations interface {
	SetToOneReferenceID(name, ID string) error
}
