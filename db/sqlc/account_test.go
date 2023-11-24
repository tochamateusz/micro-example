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
	account     Account
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
	s.db.Raw("TRUNCATE accounts CASCADE")
}

// run before each test
func (s *AccountSuite) BeforeTest(suiteName, testName string) {
	log.Println("BeforeTest()", suiteName, testName)
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

	s.account = account
}

func (s *AccountSuite) TestGetAccount() {
	resultAccount, err := s.testQueries.GetAccountByID(context.Background(), s.account.ID)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resultAccount)

	require.Equal(s.T(), s.account.Currency, resultAccount.Currency)
	require.Equal(s.T(), s.account.Balance, resultAccount.Balance)
	require.Equal(s.T(), s.account.Owner, resultAccount.Owner)

	require.Equal(s.T(), s.account.ID, resultAccount.ID)
	require.Equal(s.T(), s.account.CreatedAt, resultAccount.CreatedAt)
}

func (s *AccountSuite) TestUpdateAccount() {
	arg := UpdateAccountParams{
		ID:      s.account.ID,
		Balance: 101,
	}

	resultAccount, err := s.testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(s.T(), err)

	require.NotEmpty(s.T(), resultAccount)

	require.Equal(s.T(), s.account.Currency, resultAccount.Currency)
	require.Equal(s.T(), resultAccount.Balance, arg.Balance)
	require.Equal(s.T(), s.account.Owner, resultAccount.Owner)

	require.Equal(s.T(), s.account.ID, resultAccount.ID)
	require.Equal(s.T(), s.account.CreatedAt, resultAccount.CreatedAt)
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
