/*
 * Implementation of a commit action.
 */

package my.company.myproject.actions;

import my.company.myproject.adapters.ClusterNode;
import my.company.myproject.state.NodeState;
import my.company.myproject.state.State;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;

public class CommitAction {
    private final State state;
    private final List<ClusterNode> cluster;
    private final Logger logger;

    public CommitAction(State state, List<ClusterNode> cluster) {
        this.state = state;
        this.cluster = cluster;
        this.logger = LoggerFactory.getLogger(getClass());
    }

    public void perform() {
        if (state.getCurrent() == NodeState.LEADER) {
            int newCommit = findNewCommit();
            if (newCommit > 0) {
                logger.info("Committing at {}", newCommit);
                state.getLog().commit(newCommit);
            }
        }
    }

    private int findNewCommit() {
        for (int newCommit = state.getLog().getCommitIndex() + 1; newCommit <= Integer.MAX_VALUE; newCommit++) {
            if (!matchIndexMatches(newCommit)) {
                return newCommit - 1;
            }
            if (state.getLog().get(newCommit) != null && state.getLog().get(newCommit).getTerm() == state.getTerm()) {
                continue;
            }
            return newCommit - 1;
        }
        return -1;
    }

    private boolean matchIndexMatches(int newCommit) {
        int matches = 0;
        for (ClusterNode node : cluster) {
            if (node.getMatchIndex() >= newCommit) {
                matches++;
            }
        }
        return matches > Math.floorDiv(cluster.size(), 2);
    }
}
