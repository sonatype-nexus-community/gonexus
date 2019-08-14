# gonexus [![DepShield Badge](https://depshield.sonatype.org/badges/sonatype-nexus-community/gonexus/depshield.svg)](https://depshield.github.io)

Provides a go client library for connecting to, and interacting with, [Sonatype](//www.sonatype.com) Nexus applications such as Nexus Repository Manager and Nexus IQ Server.

## Organization of this library

The library is broken into two packages. One for each application.

### nexusrm [![GoDoc](http://godoc.org/github.com/sonatype-nexus-community/gonexus/rm?status.png)](http://godoc.org/github.com/sonatype-nexus-community/gonexus/rm) [![nexusrm coverage](https://gocover.io/_badge/github.com/sonatype-nexus-community/gonexus/rm?0 "nexusrm coverage")](http://gocover.io/github.com/sonatype-nexus-community/gonexus/rm)

Create a connection to an instance of Nexus Repository Manager

```go
// import "github.com/sonatype-nexus-community/gonexus/rm"
rm, err := nexusrm.New("http://localhost:8081", "username", "password")
if err != nil {
    panic(err)
}
```

#### Supported RM Endpoints

| Endpoint                                                                                             |         Status         |
| ---------------------------------------------------------------------------------------------------- | :--------------------: |
| [Assets](https://help.sonatype.com/repomanager3/rest-and-integration-api/assets-api)                 |      :full_moon:       |
| [Components](https://help.sonatype.com/repomanager3/rest-and-integration-api/components-api)         | :waning_gibbous_moon:  |
| [Lifecycle](https://help.sonatype.com/repomanager3/rest-and-integration-api/lifecycle-api)           |       :new_moon:       |
| [Maintenance](https://help.sonatype.com/repomanager3/rest-and-integration-api/maintenance-api) _pro_ | :waning_crescent_moon: |
| [Nodes](https://help.sonatype.com/repomanager3/rest-and-integration-api/nodes-api) _pro_             |       :new_moon:       |
| [Read-Only](https://help.sonatype.com/repomanager3/rest-and-integration-api/read-only-api)           |      :full_moon:       |
| [Repositories](https://help.sonatype.com/repomanager3/rest-and-integration-api/repositories-api)     |      :full_moon:       |
| [Script](https://help.sonatype.com/repomanager3/rest-and-integration-api/script-api)                 |      :full_moon:       |
| [Search](https://help.sonatype.com/repomanager3/rest-and-integration-api/search-api)                 | :waning_gibbous_moon:  |
| [Staging](https://help.sonatype.com/repomanager3/staging) _pro_                                      | :waning_gibbous_moon:  |
| [Status](https://help.sonatype.com/repomanager3/rest-and-integration-api/status-api)                 |      :full_moon:       |
| [Support](https://help.sonatype.com/repomanager3/rest-and-integration-api/support-api)               |      :full_moon:       |
| [Tagging](https://help.sonatype.com/repomanager3/tagging) _pro_                                      | :waning_gibbous_moon:  |
| [Tasks](https://help.sonatype.com/repomanager3/rest-and-integration-api/tasks-api)                   |       :new_moon:       |

#### Supported Provisioning API

| API        |        Status         |
| ---------- | :-------------------: |
| Core       |      :new_moon:       |
| Security   |      :new_moon:       |
| Blob Store | :waning_gibbous_moon: |
| Repository | :waning_gibbous_moon: |

**Legend**: :full_moon: == Complete :new_moon: == untouched :waning_crescent_moon::last_quarter_moon::waning_gibbous_moon: == partial support

### nexusiq [![GoDoc](http://godoc.org/github.com/sonatype-nexus-community/gonexus/iq?status.png)](http://godoc.org/github.com/sonatype-nexus-community/gonexus/iq) [![nexusiq coverage](https://gocover.io/_badge/github.com/sonatype-nexus-community/gonexus/iq?0 "nexusiq coverage")](http://gocover.io/github.com/sonatype-nexus-community/gonexus/iq)

Create a connection to an instance of Nexus IQ Server

```go
// import "github.com/sonatype-nexus-community/gonexus/iq"
iq, err := nexusiq.New("http://localhost:8070", "username", "password")
if err != nil {
    panic(err)
}

```

#### Supported IQ Endpoints

| Endpoint                                                                                                             |        Status         |   Min IQ Version    |
| -------------------------------------------------------------------------------------------------------------------- | :-------------------: | :-----------------: |
| [Application](https://help.sonatype.com/iqserver/automating/rest-apis/application-rest-apis---v2)                    |      :full_moon:      |                     |
| [Authorization Configuration](https://help.sonatype.com/iqserver/automating/rest-apis)                               |      :full_moon:      |         r70         |
| [Component Details](https://help.sonatype.com/iqserver/automating/rest-apis/component-details-rest-api---v2)         |      :full_moon:      |                     |
| [Component Evaluation](https://help.sonatype.com/iqserver/automating/rest-apis/component-evaluation-rest-apis---v2)  |      :full_moon:      |                     |
| [Component Labels](https://help.sonatype.com/iqserver/automating/rest-apis/component-labels-rest-api---v2)           |      :full_moon:      |                     |
| [Component Remediation](https://help.sonatype.com/iqserver/automating/rest-apis/component-remediation-rest-api---v2) |      :full_moon:      |         r64         |
| [Component Search](https://help.sonatype.com/iqserver/automating/rest-apis/component-search-rest-apis---v2)          |      :full_moon:      |                     |
| [Component Versions](https://help.sonatype.com/iqserver/automating/rest-apis/component-versions-rest-api---v2)       |      :full_moon:      |                     |
| [Data Retention Policy](https://help.sonatype.com/iqserver/automating/rest-apis/data-retention-policy-rest-api---v2) |      :full_moon:      |                     |
| [Organization](https://help.sonatype.com/iqserver/automating/rest-apis/organization-rest-apis---v2)                  |      :full_moon:      |                     |
| [Promote Scan](https://help.sonatype.com/iqserver/automating/rest-apis/promote-scan-rest-api---v2)                   |      :new_moon:       |                     |
| [Role](https://help.sonatype.com/iqserver/automating/rest-apis/role-rest-api---v2)                                   |      :full_moon:      |         r70         |
| [Report-related](https://help.sonatype.com/iqserver/automating/rest-apis/report-related-rest-apis---v2)              |      :full_moon:      |                     |
| [Source Control](https://help.sonatype.com/integrations/nexus-iq-for-github)                                         |      :full_moon:      |                     |
| [Success Metrics Data](https://help.sonatype.com/iqserver/automating/rest-apis/success-metrics-data-rest-api---v2)   |      :new_moon:       |                     |
| [Users](https://help.sonatype.com/iqserver/automating/rest-apis/user-rest-api---v2)                                  |      :full_moon:      |         r70         |
| [Violation](https://help.sonatype.com/iqserver/automating/rest-apis/violation-rest-api---v2)                         |      :full_moon:      |                     |
| [Webhooks](https://help.sonatype.com/iqserver/automating/iq-server-webhooks)                                         |      :waning_crescent_moon:      |                     |

**Legend**: :full_moon: == Complete :new_moon: == untouched :waning_crescent_moon::last_quarter_moon::waning_gibbous_moon: == partial support

## The Fine Print

It is worth noting that this is **NOT SUPPORTED** by [Sonatype](//www.sonatype.com), and is a contribution of [@HokieGeek](https://github.com/HokieGeek)
plus us to the open source community (read: you!)

Remember:

- Use this contribution at the risk tolerance that you have
- Do **NOT** file Sonatype support tickets related to this
- **DO** file issues here on GitHub, so that the community can pitch in
