package my.company.raft.actions

import my.company.proto.AppendRequest
import my.company.raft.adapters.ClusterNode
import my.company.raft.state.NodeState
import my.company.raft.state.State
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.async
import kotlinx.coroutines.withTimeoutOrNull
import mu.KotlinLogging
import java.lang.RuntimeException

class HeartbeatAction(private val nodeState: State, private val clusterNodes: List<ClusterNode>) {

    private val logger = KotlinLogging.logger {}

    private val log = nodeState.log

    suspend fun sendHeartbeat() {
        if (nodeState.current != NodeState.LEADER) {
            return
        }

        val heartbeatResponses = clusterNodes.map { node ->
            val prevLogIndex = node.nextIndex - 1
            val prevLogTerm = if (prevLogIndex != -1) {
                log[prevLogIndex]?.term ?: throw RuntimeException("Unable to get previous log term")
            } else {
                -1
            }
            val entries = log.starting(prevLogIndex + 1)

            val appendRequest = AppendRequest.newBuilder()
                .setTerm(nodeState.term)
                .setLeaderId(nodeState.id)
                .setLeaderCommit(nodeState.log.commitIndex)
                .setPrevLogIndex(prevLogIndex)
                .setPrevLogTerm(prevLogTerm)
                .addAllEntries(entries)
                .build()

            GlobalScope.async {
                val response = node.appendEntries(appendRequest)
                if (response == null) {
                    null
                } else {
                    node to response
                }
            }
        }.map { asyncResponse ->
            withTimeoutOrNull(200) {
                asyncResponse.await()
            }
        }.filterNotNull()

        processHeartbeatResponses(heartbeatResponses)
    }

    private fun processHeartbeatResponses(responses: List<Pair<ClusterNode, AppendRequest.Response>>) {
        responses.forEach { (node, response) ->
            when {
                response.success -> {
                    node.nextIndex = log.lastIndex() + 1
                    node.matchIndex = node.nextIndex - 1
                }
                else -> {
                    logger.info { "Heartbeat response: ${response.success}-${response.term}" }
                    node.decreaseIndex()
                }
            }
        }
    }
}
