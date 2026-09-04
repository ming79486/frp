export interface ServerInfo {
  version: string
  config: ServerInfoConfig
  status: ServerInfoStatus
}

export interface ServerInfoConfig {
  bindPort: number
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  kcpBindPort: number
  quicBindPort: number
  subdomainHost: string
  maxPoolCount: number
  maxPortsPerClient: number
  heartbeatTimeout: number
  allowPortsStr: string
  tlsForce: boolean
  transportProtocol: string
  autoTransportEnabled: boolean
  autoTransportProtocols?: string[]
}

export interface ServerInfoStatus {
  totalTrafficIn: number
  totalTrafficOut: number
  curConns: number
  clientCounts: number
  proxyTypeCount: Record<string, number>
  autoNegotiationSuccess?: number
  autoNegotiationFailure?: number
  autoTransportSelections?: Record<string, number>
  autoTransportClientCounts?: Record<string, number>
  autoTransportSwitchCounts?: Record<string, number>
  autoTransportIllegalSelections?: Record<string, number>
}
