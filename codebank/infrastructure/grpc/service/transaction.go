package service

import (
	"codebank/dto"
	"codebank/infrastructure/grpc/pb"
	"codebank/usecase"
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {

	transactionDto := dto.Transaction{
		Name:            in.GetCreditCard().GetName(),
		Number:          in.GetCreditCard().GetNumber(),
		ExpirationMonth: in.GetCreditCard().GetExpirationMonth(),
		ExpirationYear:  in.GetCreditCard().GetExpirationYear(),
		CVV:             in.GetCreditCard().GetCvv(),
		Amount:          in.GetAmount(),
		Store:           in.GetStore(),
		Description:     in.GetDescription(),
	}

	transaction, err := t.ProcessTransactionUseCase.ProcessTransaction(transactionDto)
	if err != nil {
		return &empty.Empty{}, status.Error(code.Code_FAILED_PRECONDITION, err.Error())
	}

	if transaction.Status != "approved" {
		return &empty.Empty{}, status.Error(code.Code_FAILED_PRECONDITION, "transaction rejected by bank")
	}

	return &empty.Empty{}, nil
}
