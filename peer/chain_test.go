package peer_test

import (
	"blockchain/peer"
	"encoding/json"
	"fmt"
	"testing"
)

func TestVerifyBlocks(t *testing.T) {
	var block1 peer.Block
	var block2 peer.Block

	json1 := `{"hash": "f22bce9b4d603a5d270550edb3ecf75bdee3a0c094383ca57a06d8b000000000", "height": 0, "messages": ["Rob wuz here"], "minedBy": "Rob!", "nonce": "3823127703541898446058332", "type": "GET_BLOCK_REPLY"}`
	json2 := `{"hash": "7ce995e3383f951682f6ba56483f8297cb44bcc0f49133ee577a963d00000000", "height": 1, "messages": ["kindreds", "decedent", "liegeful", "Dishley", "vermin"], "minedBy": "Rob!", "nonce": "976471045", "type": "GET_BLOCK_REPLY"}`

	json.Unmarshal([]byte(json1), &block1)
	json.Unmarshal([]byte(json2), &block2)

	fmt.Println(block2)
	if !peer.VerifyBlock(block1, block2, 9) {
		t.Error("Should be verified - true")
	}

}
