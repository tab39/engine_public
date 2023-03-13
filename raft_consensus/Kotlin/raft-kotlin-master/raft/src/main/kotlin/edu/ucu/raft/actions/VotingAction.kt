// A class for managing the voting process
class VotingAction(private val state: State, private val cluster: List<ClusterNode>) {

    // Check the term of the vote response against the current state's term
    private fun checkTerm(response: VoteResponse) {
        if (response.term > state.term) {
            // If the response has a higher term, update the state to the new term
            // and change to follower mode to handle the new leader
            state.term = response.term
            state.changeToFollower()
        }
    }

    // Send vote requests to all cluster nodes and check if candidate wins the election
    suspend fun askVotes(): Boolean {
        // Check if the cluster is empty, return false if it is
        if (cluster.isEmpty()) {
            return false
        }
        
        // Calculate the number of votes required to win the election
        val requiredVotes = (cluster.size / 2) + 1
        
        // Create a VoteRequest object with the current state's data
        val request = VoteRequest.newBuilder()
            .setTerm(state.term)
            .setCandidateId(state.id)
            .setLastLogIndex(state.log.lastIndex())
            .setLastLogTerm(state.log.lastTerm() ?: -1)
            .build()
        
        // Send vote requests to all cluster nodes asynchronously and collect responses
        val responses = cluster.map { node ->
            GlobalScope.async {
                // Use a timeout to prevent waiting indefinitely for a response
                withTimeoutOrNull(VOTE_TIMEOUT) {
                    node.requestVote(request)
                }
            }
        }.mapNotNull { it.await() }

        // Check the term of each response and update the state if necessary
        responses.forEach { checkTerm(it) }
        
        // Count the number of granted votes in the responses
        val grantedVotes = responses.count { it.voteGranted }
        
        // Check if the candidate wins the election
        return grantedVotes >= requiredVotes
    }
    
    companion object {
        // Timeout for voting
        const val VOTE_TIMEOUT = 200L
    }
}
