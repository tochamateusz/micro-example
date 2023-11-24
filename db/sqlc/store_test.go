package simple_bank

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StoreSuite struct {
	suite.Suite
	db          *gorm.DB
	testQueries *Queries
}

func TestStoreSuite(t *testing.T) {

	dns := fmt.Sprintf("host=0.0.0.0 user=root password=password_test dbname=simple_bank_test port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testingSuit := new(StoreSuite)
	testingSuit.db = db
	testingSuit.testQueries = New(testingSuit.db.ConnPool)

	suite.Run(t, testingSuit)

}

func (s *StoreSuite) SetupSuite() {
	log.Println("SetupSuite()")
}

func (s *StoreSuite) TearDownSuite() {
	log.Println("TearDownSuite()")

	s.db.Exec("TRUNCATE accounts CASCADE")
	s.db.Exec("TRUNCATE transfer CASCADE")
	s.db.Exec("TRUNCATE entries CASCADE")
}

// run before each test
func (s *StoreSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
}

// run after each test
func (s *StoreSuite) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func (s *StoreSuite) TestTransferTx() {
	db, err := s.db.DB()
	require.NoError(s.T(), err)
	store := ProvideStore(db)

	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}
	account1, err := s.testQueries.CreateAccount(context.Background(), arg)
	require.NoError(s.T(), err)

	arg = CreateAccountParams{
		Owner:    "tom2",
		Balance:  0,
		Currency: "USD",
	}
	account2, err := s.testQueries.CreateAccount(context.Background(), arg)
	require.NoError(s.T(), err)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	var amount int64 = 10
	var n = 5

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(s.T(), err)

		result := <-results
		require.NotEmpty(s.T(), result)

		require.NotEmpty(s.T(), result.Transfer)
		require.Equal(s.T(), account1.ID, result.Transfer.FromAccountID)
		require.Equal(s.T(), account2.ID, result.Transfer.ToAccountID)
		require.Equal(s.T(), amount, result.Transfer.Amount)

		require.NotZero(s.T(), result.Transfer.ID)
		require.NotZero(s.T(), result.Transfer.CreatedAt)

		_, err = store.GetTransferById(context.Background(), result.Transfer.ID)
		require.NoError(s.T(), err)

		require.NotEmpty(s.T(), result.FromEntry)
		require.Equal(s.T(), account1.ID, result.FromEntry.AccountID)
		require.Equal(s.T(), -amount, result.FromEntry.Amount)

		require.NotZero(s.T(), result.FromEntry.ID)
		require.NotZero(s.T(), result.FromEntry.CreatedAt)

		_, err = store.GetEntryByID(context.Background(), result.FromEntry.ID)
		require.NoError(s.T(), err)

		require.NotEmpty(s.T(), result.ToEntry)
		require.Equal(s.T(), account2.ID, result.ToEntry.AccountID)
		require.Equal(s.T(), amount, result.ToEntry.Amount)

		require.NotZero(s.T(), result.ToEntry.ID)
		require.NotZero(s.T(), result.ToEntry.CreatedAt)

		_, err = store.GetEntryByID(context.Background(), result.ToEntry.ID)
		require.NoError(s.T(), err)

		//TODO: Update

		fromAccount, err := store.GetAccountByID(context.Background(), result.Transfer.FromAccountID)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), fromAccount)

		require.Equal(s.T(), fromAccount.Balance, account1.Balance-int64((i+1)*10))

		toAccount, err := store.GetAccountByID(context.Background(), result.Transfer.ToAccountID)
		require.NoError(s.T(), err)
		require.NotEmpty(s.T(), toAccount)

		require.Equal(s.T(), toAccount.Balance, account2.Balance+int64((i+1)*10))

	}

}
