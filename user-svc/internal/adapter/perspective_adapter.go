package adapter

import (
	"be-yourmoments/user-svc/internal/helper/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PerspectiveAdapter interface {
	IsToxicMessage(msg string) (bool, error)
}
type perspectiveAdapter struct {
	perspectiveAPIKey string
}

func NewPerspectiveAdapter() PerspectiveAdapter {
	perspectiveAPIKey := utils.GetEnv("PERSPECTIVE_API_KEY")
	return &perspectiveAdapter{
		perspectiveAPIKey: perspectiveAPIKey,
	}
}

type PerspectiveRequest struct {
	Comment struct {
		Text string `json:"text"`
	} `json:"comment"`
	Languages           []string               `json:"languages"`
	RequestedAttributes map[string]interface{} `json:"requestedAttributes"`
}

type PerspectiveResponse struct {
	AttributeScores map[string]struct {
		SummaryScore struct {
			Value float64 `json:"value"`
		} `json:"summaryScore"`
	} `json:"attributeScores"`
}

func (a *perspectiveAdapter) IsToxicMessage(msg string) (bool, error) {
	reqBody := PerspectiveRequest{
		Languages: []string{"en"},
		RequestedAttributes: map[string]interface{}{
			"TOXICITY":        map[string]interface{}{},
			"SEVERE_TOXICITY": map[string]interface{}{},
			"IDENTITY_ATTACK": map[string]interface{}{},
			"INSULT":          map[string]interface{}{},
			"PROFANITY":       map[string]interface{}{},
			"THREAT":          map[string]interface{}{},
		},
	}

	reqBody.Comment.Text = msg
	bodyBytes, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("https://commentanalyzer.googleapis.com/v1alpha1/comments:analyze?key=%s", a.perspectiveAPIKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	var result PerspectiveResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return false, err
	}

	toxicity := result.AttributeScores["TOXICITY"].SummaryScore.Value
	severeToxicity := result.AttributeScores["SEVERE_TOXICITY"].SummaryScore.Value
	identityAttack := result.AttributeScores["IDENTITY_ATTACK"].SummaryScore.Value
	insult := result.AttributeScores["INSULT"].SummaryScore.Value
	profanity := result.AttributeScores["PROFANITY"].SummaryScore.Value
	threat := result.AttributeScores["THREAT"].SummaryScore.Value

	if toxicity >= 0.8 || severeToxicity >= 0.8 || identityAttack >= 0.8 || insult >= 0.8 || profanity >= 0.8 || threat >= 0.8 {
		return true, nil
	}

	return false, nil
}
