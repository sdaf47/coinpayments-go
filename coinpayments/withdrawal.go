package coinpayments

import (
	"net/http"

	"fmt"

	"github.com/dghubble/sling"
)

type WithdrawalService struct {
	sling        *sling.Sling
	ApiPublicKey string
	Params       WithdrawalBodyParams
}

type Withdrawal struct {
	Amount string `json:"amount"`
	Id     string `json:"id"`
	Status int    `json:"status"`
}

type WithdrawalResponse struct {
	Error  string      `json:"error"`
	Result *Withdrawal `json:"result"`
}

type WithdrawalParams struct {
	Amount      float64 `url:"amount"`
	Currency    string  `url:"currency"`
	Currency2   string  `url:"currency2"`
	Address     string  `url:"address"`
	Pbntag      string  `url:"pbntag"`
	DestTag     string  `url:"dest_tag"`
	IpnUrl      string  `url:"ipn_url"`
	AutoConfirm string  `url:"auto_confirm"` // 1
	IPNUrl      string  `url:"ipn_url"`
}

type WithdrawalBodyParams struct {
	APIParams
	WithdrawalParams
}

func newWithdrawalService(sling *sling.Sling, apiPublicKey string) *WithdrawalService {
	withdrawalService := &WithdrawalService{
		sling:        sling.Path("api.php"),
		ApiPublicKey: apiPublicKey,
	}
	withdrawalService.getParams()
	return withdrawalService
}

func (s *WithdrawalService) getHMAC() string {
	return getHMAC(getPayload(s.Params))
}

func (s *WithdrawalService) NewWithdrawal(withdrawalParams *WithdrawalParams) (WithdrawalResponse, *http.Response, error) {
	withdrawalResponse := new(WithdrawalResponse)
	s.Params.WithdrawalParams = *withdrawalParams
	fmt.Println(getPayload(s.Params))
	fmt.Println(getHMAC(getPayload(s.Params)))
	resp, err := s.sling.New().Set("HMAC", s.getHMAC()).Post(
		"api.php").BodyForm(s.Params).ReceiveSuccess(withdrawalResponse)
	return *withdrawalResponse, resp, err
}

func (s *WithdrawalService) getParams() {
	s.Params.Command = "create_withdrawal"
	s.Params.Key = s.ApiPublicKey
	s.Params.Version = "1"
}
