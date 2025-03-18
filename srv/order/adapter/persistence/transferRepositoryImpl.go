package persistence

import "sync"

type TransferRepositoryImpl struct {
	m           sync.Mutex
	transferMap map[string]Transfer
}

func NewTransferRepositoryImpl() *TransferRepositoryImpl {
	return &TransferRepositoryImpl{transferMap: make(map[string]Transfer)}
}

func (s *TransferRepositoryImpl) GetTransfer(transferId string) (Transfer, error) {
	s.m.Lock()
	defer s.m.Unlock()

	transfer, exist := s.transferMap[transferId]
	if !exist {
		return Transfer{}, ErrTransferNotFound
	}

	return transfer, nil
}

func (s *TransferRepositoryImpl) GetTransfers() ([]Transfer, error) {
	s.m.Lock()
	defer s.m.Unlock()

	transfers := make([]Transfer, 0, len(s.transferMap))
	for _, transfer := range s.transferMap {
		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (s *TransferRepositoryImpl) SetTransfer(transferId string, transfer Transfer) bool {
	s.m.Lock()
	defer s.m.Unlock()

	_, exist := s.transferMap[transferId]
	s.transferMap[transferId] = transfer

	return exist
}

func (s *TransferRepositoryImpl) SetComplete(transferId string) error {
	s.m.Lock()
	defer s.m.Unlock()

	transfer, exist := s.transferMap[transferId]
	if !exist {
		return ErrTransferNotFound
	}

	transfer.Status = "Completed"
	s.transferMap[transferId] = transfer

	return nil
}

func (s *TransferRepositoryImpl) IncrementLinkedStockUpdate(transferId string) error {
	s.m.Lock()
	defer s.m.Unlock()

	transfer, exist := s.transferMap[transferId]
	if !exist {
		return ErrTransferNotFound
	}

	transfer.LinkedStockUpdate += 1
	s.transferMap[transferId] = transfer

	return nil
}
