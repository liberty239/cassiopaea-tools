package forum

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"
)

func newSearchQuery(q string, users ...string) string {
	ret := fmt.Sprintf(
		"https://cassiopaea.org/forum/search/1/?q=%s&c[title_only]=1&o=date",
		url.QueryEscape(q),
	)

	if len(users) > 0 {
		for i := range users {
			users[i] = url.QueryEscape(users[i])
		}
		ret += "&c[users]=" + strings.Join(users, "%2C")
	}
	return ret
}

func allSessionsSearchQueryUrls() []string {
	return []string{
		newSearchQuery("Session", "Laura", "Chu", "Andromeda"),
		newSearchQuery("Sesssion", "Laura"),
	}
}

func isSearchPath(path string) bool {
	p := strings.Split(path, "/")
	if len(p) < 1 {
		return false
	}
	return p[2] == "search"
}

func timestampFromPath(path string) (ret string) {
	p := strings.Split(path, "/")
	if len(p) < 1 {
		return
	}

	if strings.HasPrefix(p[len(p)-1], "page") {
		return
	}

	ret = p[len(p)-2]
	ret = ret[:strings.Index(p[len(p)-2], ".")]
	return
}

func parseTimestamp(ts string) (ret time.Time, err error) {
	ret, err = time.Parse("session-2-January-2006", ts)
	if err != nil {
		ret, err = time.Parse("session-2-Jan-2006", ts)
		if err != nil {
			ret, err = time.Parse("sesssion-2-January-2006", ts) // case for 31 Oct 2001
		}
	}
	return
}

func absoluteUrl(rawurl string) string {
	url, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}

	if len(url.Host) == 0 {
		url.Scheme = "https"
		url.Host = "cassiopaea.org"
	}

	return url.String()
}

func packageName() string {
	type test struct{}
	return reflect.TypeOf(test{}).PkgPath()
}
