package speedtest_util

import (
	"context"
	"fmt"
	"log"

	"speedtest-tray/internal/config"
)

func (st *SpeedTester) initialize(ctx context.Context, updateCh chan<- Update) error {
	updateCh <- Update{Phase: INITIALIZING, Progress: config.ProgressInit}

	updateCh <- Update{Phase: GETTING_INFO, Progress: config.ProgressGetInfo}
	user, err := st.client.FetchUserInfoContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	st.user = user
	return nil
}

func (st *SpeedTester) selectBestServer(ctx context.Context, updateCh chan<- Update, res *Result) error {
	updateCh <- Update{Phase: FINDING_SERVERS, Progress: config.ProgressFindServers}
	servers, err := st.client.FetchServers()
	if err != nil {
		return fmt.Errorf("failed to fetch servers: %w", err)
	}
	st.servers = servers

	updateCh <- Update{Phase: SELECTING_SERVER, Progress: config.ProgressSelectServer}
	targets, err := st.servers.FindServer([]int{})
	if err != nil {
		return fmt.Errorf("failed to find target server: %w", err)
	}
	st.server = targets[0]
	serverInfo := fmt.Sprintf("%s (%s)", st.server.Name, st.server.Country)
	res.Server = serverInfo
	log.Printf("Selected server: %s\n", serverInfo)

	updateCh <- Update{Phase: SERVER_SELECTED, Progress: config.ProgressServerSelect, Server: serverInfo}
	return nil
}
