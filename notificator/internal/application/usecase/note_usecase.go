package usecase

type noteUsecase struct {
	//repo            repository.NoteRepository
}

type NoteUsecase interface {
}

func NewNoteUsecase() NoteUsecase {
	return &noteUsecase{}
}
