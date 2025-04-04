package persistence

import (
	"github.com/alimitedgroup/MVP/srv/order/business/model"
	"github.com/alimitedgroup/MVP/srv/order/business/port"
)

type TransferPersistanceAdapter struct {
	transferRepo ITransferRepository
}

func NewTransferPersistanceAdapter(transferRepo ITransferRepository) *TransferPersistanceAdapter {
	return &TransferPersistanceAdapter{transferRepo}
}

func (s *TransferPersistanceAdapter) SetComplete(transferId model.TransferID) error {
	err := s.transferRepo.SetComplete(string(transferId))
	if err != nil {
		return err
	}

	return nil
}

func (s *TransferPersistanceAdapter) IncrementLinkedStockUpdate(transferId model.TransferID) error {
	err := s.transferRepo.IncrementLinkedStockUpdate(string(transferId))
	if err != nil {
		return err
	}
	return nil
}

func (s *TransferPersistanceAdapter) ApplyTransferUpdate(cmd port.ApplyTransferUpdateCmd) {
	status := cmd.Status
	linkedStockUpdate := 0

	if old, err := s.transferRepo.GetTransfer(cmd.ID); err == nil {
		linkedStockUpdate = old.LinkedStockUpdate
		if old.Status == "Completed" {
			status = old.Status
		}

	}

	goods := make([]TransferUpdateGood, 0, len(cmd.Goods))
	for _, good := range cmd.Goods {
		goods = append(goods, TransferUpdateGood{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	transfer := Transfer{
		ID:                cmd.ID,
		Status:            status,
		SenderID:          cmd.SenderID,
		ReceiverID:        cmd.ReceiverID,
		LinkedStockUpdate: linkedStockUpdate,
		Goods:             goods,
		ReservationID:     cmd.ReservationID,
		UpdateTime:        cmd.UpdateTime,
		CreationTime:      cmd.CreationTime,
	}

	s.transferRepo.SetTransfer(cmd.ID, transfer)
}

func (s *TransferPersistanceAdapter) GetTransfer(transferId model.TransferID) (model.Transfer, error) {
	transfer, err := s.transferRepo.GetTransfer(string(transferId))
	if err != nil {
		return model.Transfer{}, err
	}

	modelTransfer := repoTransferToModelTransfer(transfer)
	return modelTransfer, nil
}

func (s *TransferPersistanceAdapter) GetAllTransfer() []model.Transfer {
	transfers := s.transferRepo.GetTransfers()
	modelTransfer := repoTransfersToModelTransfers(transfers)
	return modelTransfer
}

func repoTransferToModelTransfer(transfer Transfer) model.Transfer {
	goods := make([]model.GoodStock, 0, len(transfer.Goods))
	for _, good := range transfer.Goods {
		goods = append(goods, model.GoodStock{
			GoodID:   good.GoodID,
			Quantity: good.Quantity,
		})
	}

	return model.Transfer{
		ID:                transfer.ID,
		Status:            transfer.Status,
		SenderID:          transfer.SenderID,
		ReceiverID:        transfer.ReceiverID,
		Goods:             goods,
		LinkedStockUpdate: transfer.LinkedStockUpdate,
		ReservationID:     transfer.ReservationID,
		UpdateTime:        transfer.UpdateTime,
		CreationTime:      transfer.CreationTime,
	}
}

func repoTransfersToModelTransfers(transfers []Transfer) []model.Transfer {
	list := make([]model.Transfer, 0, len(transfers))
	for _, transfer := range transfers {
		list = append(list, repoTransferToModelTransfer(transfer))
	}

	return list
}
