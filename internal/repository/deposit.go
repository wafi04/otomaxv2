package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/wafi04/otomaxv2/internal/model"
)

type DepositRepository struct {
	db *sql.DB
}

func NewDepositRepository(db *sql.DB) *DepositRepository {
	return &DepositRepository{
		db: db,
	}
}

func (repo *DepositRepository) Create(c context.Context, req model.CreateDeposit) (bool, error) {
	query := `
		INSERT INTO deposits (
			username,
			method,
			amount,
			payment_referee,
			destination_number,
			status,
			created_at,
			updated_at
		) VALUES (
			$1,$2,$3,$4,$5,'PENDING',NOW(),NOW()
		)
	`
	_, err := repo.db.ExecContext(c, query, &req.Username, &req.Method, &req.PaymentReferee, &req.DestinationNumber)

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (repo *DepositRepository) GetAll(c context.Context, req model.FilterDeposit) ([]model.DepositData, int, error) {
	countQuery := `
		SELECT COUNT(*) 
		FROM categories
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR status  = $2)
	`

	var totalCount int
	err := repo.db.QueryRowContext(c, countQuery, req.Search, req.Status).Scan(&totalCount)
	if err != nil {
		log.Printf("GetAll Categories count error: %v", err)
		return nil, 0, err
	}

	query := `
		SELECT 
			id, 
			username,
			method,
			amount,
			payement_referee,
			destination_number,
			status,
			created_at, 
			updated_at
		FROM deposits
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`
	rows, err := repo.db.QueryContext(c, query, req.Search, req.Status, req.Limit, req.Offset)
	if err != nil {
		log.Printf("GetAll Deposits error: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var deposits []model.DepositData
	for rows.Next() {
		var dep model.DepositData
		err := rows.Scan(
			&dep.ID, &dep.Username, &dep.Method, &dep.Amount, &dep.PaymentReferee,
			&dep.DestinationNumber, &dep.Status, &dep.CreatedAt, &dep.UpdatedAt,
		)
		if err != nil {
			log.Printf("Scan Deposit error: %v", err)
			continue
		}
		deposits = append(deposits, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return deposits, totalCount, nil

}

func (repo *DepositRepository) GetByID(ctx context.Context, id int) (*model.DepositData, error) {
	query := `
		SELECT 
			id, 
			username,
			method,
			amount,
			payement_referee,
			destination_number,
			status,
			created_at, 
			updated_at
		FROM deposits
		WHERE id = $1
	`

	var dep model.DepositData
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&dep.ID, &dep.Username, &dep.Method, &dep.Amount, &dep.PaymentReferee,
		&dep.DestinationNumber, &dep.Status, &dep.CreatedAt, &dep.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("GetByID Category error: %v", err)
		return nil, err
	}
	return &dep, nil
}
