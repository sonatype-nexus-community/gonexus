package nexusrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hokiegeek/gonexus"
)

// http://localhost:8081/service/rest/v1/components?continuationToken=foo&repository=bar
const restListComponentsByRepo = "%s/service/rest/v1/components?repository=%s"

const hashPart = 20

// RepositoryItem holds the data of a component in a repository
type RepositoryItem struct {
	ID         string `json:"id"`
	Repository string `json:"repository"`
	Format     string `json:"format"`
	Group      string `json:"group"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Assets     []struct {
		DownloadURL string `json:"downloadUrl"`
		Path        string `json:"path"`
		ID          string `json:"id"`
		Repository  string `json:"repository"`
		Format      string `json:"format"`
		Checksum    struct {
			Sha1 string `json:"sha1"`
			Md5  string `json:"md5"`
		} `json:"checksum"`
	} `json:"assets"`
	Tags []interface{} `json:"tags"`
}

// Hash is a hack which returns the most pertinent hash for a component
func (i *RepositoryItem) Hash() string {
	var hash string

	sumByExt := func(exts []string) string {
		ext := exts[0]
		for _, ass := range i.Assets {
			if strings.HasSuffix(ass.Path, ext) {
				return ass.Checksum.Sha1
			}
		}
		return ""
	}

	switch i.Format {
	case "maven2":
		hash = sumByExt([]string{"jar"})
	case "rubygems":
		hash = sumByExt([]string{"gem"})
	case "npm":
		hash = sumByExt([]string{"tar.gz"})
	case "pipy":
		hash = sumByExt([]string{"whl", "tar.gz"})
	default:
		hash = ""
	}
	if len(hash) < hashPart {
		return hash
	}
	return hash[0:hashPart]
}

type listComponentsResponse struct {
	Items             []RepositoryItem `json:"items"`
	ContinuationToken string           `json:"continuationToken"`
}

// RM holds basic and state info of the Repository Manager server we will connect to
type RM struct {
	rest *nexus.Server
}

// ListComponents returns a list of components in the indicated repository
func (rm *RM) ListComponents(repo string) (items []RepositoryItem, err error) {
	continuation := ""

	getComponents := func() (listResp listComponentsResponse, err error) {
		url := fmt.Sprintf(restListComponentsByRepo, rm.host, repo)

		if continuation != "" {
			url += "&continuationToken=" + continuation
		}

		fmt.Println("here")
		body, resp, err := rest.Get(url)
		if err != nil || resp.StatusCode == http.StatusOK {
			return
		}

		if err = json.Unmarshal(body, &listResp); err != nil {
			return
		}

		/*
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return
			}

			client := &http.Client{}
			req.SetBasicAuth(rm.username, rm.password)
			resp, err := client.Do(req)
			if err != nil {
				return
			}

			if resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)

				if err = json.Unmarshal(body, &listResp); err != nil {
					return
				}
			}
		*/

		return
	}

	for {
		resp, err := getComponents()
		if err != nil {
			return items, err
		}

		items = append(items, resp.Items...)

		// for _, item := range rmItems.Items {
		// 	iqc := rmItemToIQComponent(item)
		// 	if iqc.Hash != "" {
		// 		c = append(c, iqc)
		// 	}
		// }

		if resp.ContinuationToken == "" {
			break
		}
		continuation = resp.ContinuationToken
	}

	return
}

/*
func rmFwComponents(rmServer, iqServer serverInfo, repo string) (*results, error) {
	// fmt.Println(":: Getting components from proxy")
	rmComponents, err := getRepoComponents(rmServer, repo)
	if err != nil {
		return nil, err
	}

	iq, err := nexusiq.New(iqServer.host, iqServer.auth.username, iqServer.auth.password)
	if err != nil {
		return nil, err
	}

	// fmt.Println(":: Evaluating components")
	report, err := iq.EvaluateComponentsAsFirewall(rmComponents)
	if err != nil {
		return nil, err
	}

	if report.IsError {
		panic(fmt.Sprintf("%q", report.ErrorMessage))
	}

	r := new(results)

	for _, result := range report.Results {
		r.Components = append(r.Components, iqResultToComponent(result))
	}

	return r, nil
}
*/

// New creates a new Repository Manager instance
func New(host, username, password string) (rm *RM, err error) {
	rm = new(RM)
	rm.rest, err = nexus.NewServer(host, username, password)
	// rm.host = host
	// rm.username = username
	// rm.password = password
	return
}
