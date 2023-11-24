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

type AccountSuite struct {
	suite.Suite
	db          *gorm.DB
	testQueries *Queries
}

func TestAccountSuite(t *testing.T) {

	dns := fmt.Sprintf("host=0.0.0.0 user=root password=password_test dbname=simple_bank_test port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testingSuit := new(AccountSuite)
	testingSuit.db = db
	testingSuit.testQueries = New(testingSuit.db.ConnPool)

	suite.Run(t, testingSuit)

}

func (s *AccountSuite) SetupSuite() {
	log.Println("SetupSuite()")
}

func (s *AccountSuite) TearDownSuite() {
	log.Println("TearDownSuite()")
	s.db.Exec("TRUNCATE accounts CASCADE")
}

// run before each test
func (s *AccountSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
	s.db.Exec("TRUNCATE accounts CASCADE")
}

// run after each test
func (s *AccountSuite) AfterTest(suiteName, testName string) {
	log.Println("AfterTest()", suiteName, testName)
}

func (s *AccountSuite) TestCreateAccount() {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}
	account, err := s.testQueries.CreateAccount(context.Background(), arg)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), account)

	require.Equal(s.T(), arg.Owner, account.Owner)
	require.Equal(s.T(), arg.Balance, account.Balance)
	require.Equal(s.T(), arg.Currency, account.Currency)

	require.NotZero(s.T(), account.ID)
	require.NotZero(s.T(), account.CreatedAt)
}

func (s *AccountSuite) TestGetAccount() {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	account, err := s.testQueries.CreateAccount(context.Background(), arg)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), account)

	resultAccount, err := s.testQueries.GetAccountByID(context.Background(), account.ID)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resultAccount)

	require.Equal(s.T(), account.Currency, resultAccount.Currency)
	require.Equal(s.T(), account.Balance, resultAccount.Balance)
	require.Equal(s.T(), account.Owner, resultAccount.Owner)

	require.Equal(s.T(), account.ID, resultAccount.ID)
	require.Equal(s.T(), account.CreatedAt, resultAccount.CreatedAt)
}

func (s *AccountSuite) TestUpdateAccount() {

	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	account, err := s.testQueries.CreateAccount(context.Background(), arg)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), account)

	updateArgs := UpdateAccountParams{
		ID:      account.ID,
		Balance: 101,
	}

	resultAccount, err := s.testQueries.UpdateAccount(context.Background(), updateArgs)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resultAccount)

	require.Equal(s.T(), account.Currency, resultAccount.Currency)
	require.Equal(s.T(), resultAccount.Balance, updateArgs.Balance)
	require.Equal(s.T(), account.Owner, resultAccount.Owner)

	require.Equal(s.T(), account.ID, resultAccount.ID)
	require.Equal(s.T(), account.CreatedAt, resultAccount.CreatedAt)
}

func (s *AccountSuite) TestDeleteAccount() {
	arg := CreateAccountParams{
		Owner:    "tom",
		Balance:  100,
		Currency: "USD",
	}

	account, err := s.testQueries.CreateAccount(context.Background(), arg)

	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), account)

	err = s.testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(s.T(), err)

	resultAccount, err := s.testQueries.GetAccountByID(context.Background(), account.ID)
	require.Error(s.T(), err)
	require.Empty(s.T(), resultAccount)

}

func (s *AccountSuite) TestListAccount() {
	var accounts []Account

	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner:    fmt.Sprintf("tom_%d", i),
			Balance:  100 + int64(i),
			Currency: "USD",
		}
		account, err := s.testQueries.CreateAccount(context.Background(), arg)
		require.NoError(s.T(), err)
		accounts = append(accounts, account)
	}

	require.NotEmpty(s.T(), accounts)

	arg := GetAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	resultAccount, err := s.testQueries.GetAccounts(context.Background(), arg)
	require.NoError(s.T(), err)
	require.NotEmpty(s.T(), resultAccount)

	require.Equal(s.T(), len(resultAccount), 5)
	require.Equal(s.T(), accounts[0:5], resultAccount)

}
