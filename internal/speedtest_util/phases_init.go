package speedtest_util

import (
	"context"
	"fmt"
	"log"
)

func (st *SpeedTester) initialize(ctx context.Context, updateCh chan<- Update) error {
	updateCh <- Update{Phase: INITIALIZING, Progress: 0.0}

	updateCh <- Update{Phase: GETTING_INFO, Progress: 0.05}
	user, err := st.client.FetchUserInfoContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	st.user = user
	return nil
}

func (st *SpeedTester) selectBestServer(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: FINDING_SERVERS, Progress: 0.10}
	servers, err := st.client.FetchServers()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}
	st.servers = servers

	updateCh <- Update{Phase: SELECTING_SERVER, Progress: 0.12}
	targets, err := st.servers.FindServer([]int{})
	if err != nil {
		return fmt.Errorf("failed to find target server: %w", err)
	}
	st.server = targets[0]
	serverInfo := fmt.Sprintf("%s (%s)", st.server.Name, st.server.Country)
	res.Server = serverInfo
	log.Printf("Selected server: %s\n", serverInfo)

	updateCh <- Update{Phase: SERVER_SELECTED, Progress: 0.15, Server: serverInfo}
	return nil
}
