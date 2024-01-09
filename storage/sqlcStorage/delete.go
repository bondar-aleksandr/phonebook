package sqlcstorage

import (
	"context"
	"fmt"
	"github.com/bondar-aleksandr/phonebook/storage"
)

func(s *Storage) Delete(ctx context.Context, cd *storage.CrudData) (int64, error) {
	switch cd.CrudComb {
	case storage.CrudFname:
		count, err := s.queries.DeletePersonByFname(ctx, cd.FirstName)
		if err != nil {
			return count, fmt.Errorf("can't delete person by fname: %w", err)
		}
		return count, nil
	case storage.CrudPhone:
		count, err := s.queries.DeletePhoneByNumber(ctx, cd.Phone)
		if err != nil {
			return count, fmt.Errorf("can't delete phone by number: %w", err)
		}
		return count, nil
	}
	return 0, nil
}