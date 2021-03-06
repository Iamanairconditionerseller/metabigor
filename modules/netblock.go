package modules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/j3ssie/metabigor/core"
	"github.com/thoas/go-funk"
)

func getAsnNum(raw string) string {
	if strings.HasPrefix(strings.ToLower(raw), "as") {
		return raw[2:]
	}
	return raw
}

// IPInfo get CIDR from ASN
func IPInfo(options core.Options) []string {
	asn := getAsnNum(options.Net.Asn)
	url := fmt.Sprintf(`https://ipinfo.io/AS%v`, asn)
	var result []string
	core.InforF("Get data from: %v", url)
	content := core.RequestWithChrome(url, "ipTabContent")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return result
	}
	// searching for data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		s.Find("address").First()
		if !strings.Contains(s.Text(), "Netblock") {
			data := strings.Split(strings.TrimSpace(s.Text()), "  ")
			cidr := strings.TrimSpace(data[0])
			desc := strings.TrimSpace(data[len(data)-1])

			core.InforF(fmt.Sprintf("%s - %s", cidr, desc))
			result = append(result, fmt.Sprintf("%s", cidr))
		}

	})
	return result
}

// IPv4Info get CIDR from ASN via ipv4info.com
func IPv4Info(options core.Options) []string {
	asn := getAsnNum(options.Net.Asn)
	url := fmt.Sprintf(`http://ipv4info.com/?act=check&ip=AS%v`, asn)
	var result []string

	core.InforF("Get data from: %v", url)
	content := core.SendGET(url, options)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return result
	}

	// finding ID of block
	var ASNLink []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			if strings.HasPrefix(href, "/org/") {
				ASNLink = append(ASNLink, href)
			}
		}
	})

	// searching for data
	ASNLink = funk.Uniq(ASNLink).([]string)
	for _, link := range ASNLink {
		core.InforF("Get data from: %v", link)
		url := fmt.Sprintf(`http://ipv4info.com%v`, link)
		core.InforF("Get data from: %v", url)
		content := core.SendGET(url, options)
		// finding ID of block
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			return result
		}

		doc.Find("td").Each(func(i int, s *goquery.Selection) {
			style, _ := s.Attr("style")
			class, _ := s.Attr("class")
			if style == "padding: 0 0 0 0;" && class == "bold" {
				data := s.Text()
				result = append(result, data)
			}

		})
	}
	core.InforF("\n%v", strings.Join(result, "\n"))
	return result
}

// ASNBgpDotNet get ASN infor from bgp.net
func ASNBgpDotNet(options core.Options) []string {
	asn := getAsnNum(options.Net.Asn)
	url := fmt.Sprintf(`https://bgp.he.net/AS%v#_prefixes`, asn)
	core.InforF("Get data from: %v", url)
	var result []string
	content := core.RequestWithChrome(url, "prefixes")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return result
	}
	// searching for data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		data := strings.Split(strings.TrimSpace(s.Text()), "  ")
		cidr := strings.TrimSpace(data[0])
		if !strings.Contains(cidr, "Prefix") {
			desc := strings.TrimSpace(data[len(data)-1])
			core.InforF(fmt.Sprintf("%s - %s", cidr, desc))
			result = append(result, fmt.Sprintf("%s", cidr))
		}
	})
	return result
}

// ASNSpyse get ASN infor from spyse.com
func ASNSpyse(options core.Options) []string {
	asn := getAsnNum(options.Net.Asn)
	url := fmt.Sprintf(`https://spyse.com/target/as/%v#c-domain__anchor--3--%v`, asn, asn)
	var result []string
	core.InforF("Get data from: %v", url)
	content := core.RequestWithChrome(url, "asn-ipv4-ranges")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return result
	}
	// searching for data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		data := strings.Split(strings.TrimSpace(s.Text()), "  ")
		cidr := strings.TrimSpace(data[0])
		if !strings.Contains(cidr, "CIDR") {
			desc := strings.Split(data[len(data)-2], "\n")
			realDesc := desc[len(desc)-1]
			core.InforF(fmt.Sprintf("%s - %s", cidr, realDesc))
			result = append(result, fmt.Sprintf("%s", cidr))
		}
	})
	return result
}

/* Get IP range from Organization */

// OrgBgpDotNet get Org infor from bgp.net
func OrgBgpDotNet(options core.Options) []string {
	org := options.Net.Org
	url := fmt.Sprintf(`https://bgp.he.net/search?search%%5Bsearch%%5D=%v&commit=Search`, org)
	core.InforF("Get data from: %v", url)
	var result []string
	content := core.RequestWithChrome(url, "search")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return result
	}

	// searching for data
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		if !strings.Contains(s.Text(), "Result") && !strings.Contains(s.Text(), "Description") {
			data := strings.Split(strings.TrimSpace(s.Text()), "  ")[0]
			realdata := strings.Split(data, "\n")
			cidr := strings.TrimSpace(realdata[0])
			desc := strings.TrimSpace(realdata[len(realdata)-1])
			core.InforF(fmt.Sprintf("%s - %s", cidr, desc))
			result = append(result, fmt.Sprintf("%s", cidr))
		}
	})
	return result
}

// ASNLookup get Org CIDR from asnlookup
func ASNLookup(options core.Options) []string {
	org := options.Net.Org
	url := fmt.Sprintf(`http://asnlookup.com/api/lookup?org=%v`, org)
	core.InforF("Get data from: %v", url)
	data := core.SendGET(url, options)
	var result []string
	if data == "" {
		return result
	}
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return result
	}

	for _, item := range result {
		core.InforF(item)
	}
	return result
}
